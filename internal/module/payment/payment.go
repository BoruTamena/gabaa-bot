package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/lakipay"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type OrderCallbacks interface {
	OnPaymentSuccess(ctx context.Context, orderID int64) error
	OnPaymentFailed(ctx context.Context, orderID int64) error
}

type paymentModule struct {
	paymentStorage        storage.PaymentStorage
	paymentWebhookStorage storage.PaymentWebhookStorage
	escrowStorage         storage.EscrowStorage
	walletStorage         storage.WalletStorage
	withdrawalStorage     storage.WithdrawalStorage
	orderStorage          storage.OrderStorage
	lakipay               platform.LakiPay
	orderCallbacks        OrderCallbacks
}

func NewPaymentModule(
	paymentStorage storage.PaymentStorage,
	paymentWebhookStorage storage.PaymentWebhookStorage,
	escrowStorage storage.EscrowStorage,
	walletStorage storage.WalletStorage,
	withdrawalStorage storage.WithdrawalStorage,
	orderStorage storage.OrderStorage,
	lakipayClient platform.LakiPay,
	orderCallbacks OrderCallbacks,
) *paymentModule {
	return &paymentModule{
		paymentStorage:        paymentStorage,
		paymentWebhookStorage: paymentWebhookStorage,
		escrowStorage:         escrowStorage,
		walletStorage:         walletStorage,
		withdrawalStorage:     withdrawalStorage,
		orderStorage:          orderStorage,
		lakipay:               lakipayClient,
		orderCallbacks:        orderCallbacks,
	}
}

func (m *paymentModule) InitiateForOrder(ctx context.Context, order *db.Order, medium, phone string) (*dto.Payment, error) {
	phone = normalizePhone(phone)
	if phone == "" {
		return nil, fmt.Errorf("phone number is required for payment")
	}

	payment := &db.Payment{
		OrderID:     order.ID,
		Status:      constant.PaymentStatusInitiated,
		Method:      "lakipay",
		Amount:      order.TotalPrice,
		Currency:    "ETB",
		PhoneNumber: phone,
		Medium:      medium,
	}
	if err := m.paymentStorage.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	reference := fmt.Sprintf("ORDER-%d-%d", order.ID, payment.ID)
	payment.Reference = reference
	if err := m.paymentStorage.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	callbackURI := viper.GetString("lakipay.callback_url")
	if callbackURI == "" {
		appURL := strings.TrimRight(viper.GetString("app.url"), "/")
		callbackURI = appURL + "/api/v1/webhook/lakipay"
	}

	logger.Info("initiating lakipay payment",
		zap.Int64("order_id", order.ID),
		zap.Int64("payment_id", payment.ID),
		zap.String("reference", reference),
		zap.String("medium", medium),
		zap.String("callback_uri", callbackURI),
	)

	resp, err := m.lakipay.InitiateDirectPayment(ctx, lakipay.DirectPaymentRequest{
		Amount:      order.TotalPrice,
		Currency:    "ETB",
		PhoneNumber: phone,
		Medium:      medium,
		Description: fmt.Sprintf("Payment for order #%d", order.ID),
		Reference:   reference,
		CallbackURI: callbackURI,
	})
	if err != nil {
		logger.Error("lakipay direct payment failed",
			zap.Error(err),
			zap.Int64("order_id", order.ID),
			zap.Int64("payment_id", payment.ID),
			zap.String("reference", reference),
			zap.String("medium", medium),
			zap.String("phone", phone),
		)
		payment.Status = constant.PaymentStatusFailed
		payment.GatewayStatus = constant.GatewayPaymentStatusFailed
		_ = m.paymentStorage.UpdatePayment(ctx, payment)
		_ = m.orderCallbacks.OnPaymentFailed(ctx, order.ID)
		return nil, fmt.Errorf("payment initiation failed: %w", err)
	}

	gatewayResp, _ := json.Marshal(resp)
	txnID := resp.Data.TransactionID
	payment.Status = constant.PaymentStatusPending
	payment.TransactionID = &txnID
	payment.GatewayStatus = constant.ParseGatewayPaymentStatus(resp.Data.Status)
	payment.GatewayResponse = gatewayResp
	if err := m.paymentStorage.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return mapPaymentToDTO(payment), nil
}

func (m *paymentModule) HandleWebhook(ctx context.Context, rawBody []byte) dto.WebhookResult {
	payloadMap, err := lakipay.PayloadToStringMap(rawBody)
	if err != nil {
		return dto.WebhookResult{StatusCode: 400, Message: "invalid payload"}
	}

	signature := payloadMap["signature"]
	verified, err := m.lakipay.VerifyWebhookSignature(payloadMap, signature)
	if err != nil {
		logger.Error("webhook signature verification error", zap.Error(err))
	}

	webhookEvent := &db.PaymentWebhook{
		TransactionID: payloadMap["transaction_id"],
		Event:         payloadMap["event"],
		Status:        payloadMap["status"],
		Payload:       rawBody,
		Signature:     signature,
		Verified:      verified,
		ReceivedAt:    time.Now(),
	}
	_ = m.paymentWebhookStorage.CreateWebhookEvent(ctx, webhookEvent)

	if !verified {
		return dto.WebhookResult{StatusCode: 401, Message: "invalid signature"}
	}

	event := strings.ToUpper(strings.TrimSpace(payloadMap["event"]))
	switch event {
	case constant.WebhookEventWithdrawal:
		return m.handleWithdrawalWebhook(ctx, payloadMap, webhookEvent)
	default:
		return m.handleDepositWebhook(ctx, payloadMap, webhookEvent)
	}
}

func (m *paymentModule) handleDepositWebhook(ctx context.Context, payloadMap map[string]string, webhookEvent *db.PaymentWebhook) dto.WebhookResult {
	transactionID := payloadMap["transaction_id"]
	reference := payloadMap["reference"]
	gatewayStatus := constant.ParseGatewayPaymentStatus(payloadMap["status"])

	var payment *db.Payment
	var err error
	if transactionID != "" {
		payment, err = m.paymentStorage.GetPaymentByTransactionID(ctx, transactionID)
	}
	if (err != nil || payment == nil) && reference != "" {
		payment, err = m.paymentStorage.GetPaymentByReference(ctx, reference)
	}
	if err != nil || payment == nil {
		logger.Error("payment not found for webhook", zap.String("transaction_id", transactionID), zap.String("reference", reference))
		return dto.WebhookResult{StatusCode: 404, Message: "payment not found"}
	}

	webhookEvent.PaymentID = &payment.ID

	if payment.Status.IsTerminal() {
		_ = m.paymentWebhookStorage.MarkWebhookProcessed(ctx, webhookEvent.ID)
		return dto.WebhookResult{StatusCode: 200, Message: "already processed"}
	}

	switch gatewayStatus {
	case constant.GatewayPaymentStatusSuccess:
		if err := m.processPaymentSuccess(ctx, payment, gatewayStatus); err != nil {
			logger.Error("failed to process payment success", zap.Error(err), zap.Int64("payment_id", payment.ID))
			return dto.WebhookResult{StatusCode: 500, Message: "processing failed"}
		}
	case constant.GatewayPaymentStatusFailed:
		if err := m.processPaymentFailed(ctx, payment, gatewayStatus); err != nil {
			logger.Error("failed to process payment failure", zap.Error(err), zap.Int64("payment_id", payment.ID))
			return dto.WebhookResult{StatusCode: 500, Message: "processing failed"}
		}
	default:
		payment.GatewayStatus = gatewayStatus
		_ = m.paymentStorage.UpdatePayment(ctx, payment)
	}

	_ = m.paymentWebhookStorage.MarkWebhookProcessed(ctx, webhookEvent.ID)
	return dto.WebhookResult{StatusCode: 200, Message: "ok"}
}

func (m *paymentModule) handleWithdrawalWebhook(ctx context.Context, payloadMap map[string]string, webhookEvent *db.PaymentWebhook) dto.WebhookResult {
	transactionID := payloadMap["transaction_id"]
	reference := payloadMap["reference"]
	gatewayStatus := constant.ParseGatewayPaymentStatus(payloadMap["status"])

	var withdrawal *db.Withdrawal
	var err error
	if transactionID != "" {
		withdrawal, err = m.withdrawalStorage.GetWithdrawalByTransactionID(ctx, transactionID)
	}
	if (err != nil || withdrawal == nil) && reference != "" {
		withdrawal, err = m.withdrawalStorage.GetWithdrawalByReference(ctx, reference)
	}
	if err != nil || withdrawal == nil {
		logger.Error("withdrawal not found for webhook", zap.String("transaction_id", transactionID), zap.String("reference", reference))
		return dto.WebhookResult{StatusCode: 404, Message: "withdrawal not found"}
	}

	webhookEvent.WithdrawalID = &withdrawal.ID

	if withdrawal.Status.IsTerminal() {
		_ = m.paymentWebhookStorage.MarkWebhookProcessed(ctx, webhookEvent.ID)
		return dto.WebhookResult{StatusCode: 200, Message: "already processed"}
	}

	switch gatewayStatus {
	case constant.GatewayPaymentStatusSuccess:
		if err := m.processWithdrawalSuccess(ctx, withdrawal, gatewayStatus); err != nil {
			logger.Error("failed to process withdrawal success", zap.Error(err), zap.Int64("withdrawal_id", withdrawal.ID))
			return dto.WebhookResult{StatusCode: 500, Message: "processing failed"}
		}
	case constant.GatewayPaymentStatusFailed:
		if err := m.processWithdrawalFailed(ctx, withdrawal, gatewayStatus, constant.WithdrawalStatusFailed); err != nil {
			logger.Error("failed to process withdrawal failure", zap.Error(err), zap.Int64("withdrawal_id", withdrawal.ID))
			return dto.WebhookResult{StatusCode: 500, Message: "processing failed"}
		}
	case constant.GatewayPaymentStatusCancelled:
		if err := m.processWithdrawalFailed(ctx, withdrawal, gatewayStatus, constant.WithdrawalStatusCancelled); err != nil {
			logger.Error("failed to process withdrawal cancellation", zap.Error(err), zap.Int64("withdrawal_id", withdrawal.ID))
			return dto.WebhookResult{StatusCode: 500, Message: "processing failed"}
		}
	default:
		withdrawal.GatewayStatus = gatewayStatus
		_ = m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal)
	}

	_ = m.paymentWebhookStorage.MarkWebhookProcessed(ctx, webhookEvent.ID)
	return dto.WebhookResult{StatusCode: 200, Message: "ok"}
}

func (m *paymentModule) processWithdrawalSuccess(ctx context.Context, withdrawal *db.Withdrawal, gatewayStatus constant.GatewayPaymentStatus) error {
	withdrawal.Status = constant.WithdrawalStatusSuccess
	withdrawal.GatewayStatus = gatewayStatus
	if err := m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal); err != nil {
		return err
	}
	return m.walletStorage.CompleteWithdrawal(ctx, withdrawal.StoreID, withdrawal.Amount)
}

func (m *paymentModule) processWithdrawalFailed(ctx context.Context, withdrawal *db.Withdrawal, gatewayStatus constant.GatewayPaymentStatus, status constant.WithdrawalStatus) error {
	withdrawal.Status = status
	withdrawal.GatewayStatus = gatewayStatus
	if err := m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal); err != nil {
		return err
	}
	return m.walletStorage.UnlockWithdrawal(ctx, withdrawal.StoreID, withdrawal.Amount)
}

func (m *paymentModule) processPaymentSuccess(ctx context.Context, payment *db.Payment, gatewayStatus constant.GatewayPaymentStatus) error {
	order, err := m.orderStorage.GetOrderByID(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	payment.Status = constant.PaymentStatusSuccess
	payment.GatewayStatus = gatewayStatus
	if err := m.paymentStorage.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	escrow := &db.Escrow{
		OrderID:  payment.OrderID,
		StoreID:  order.StoreID,
		Amount:   payment.Amount,
		Currency: payment.Currency,
		Status:   constant.EscrowStatusHeld,
	}
	if err := m.escrowStorage.CreateEscrow(ctx, escrow); err != nil {
		return err
	}

	if err := m.walletStorage.AddPendingBalance(ctx, order.StoreID, payment.Amount); err != nil {
		return err
	}

	return m.orderCallbacks.OnPaymentSuccess(ctx, payment.OrderID)
}

func (m *paymentModule) processPaymentFailed(ctx context.Context, payment *db.Payment, gatewayStatus constant.GatewayPaymentStatus) error {
	payment.Status = constant.PaymentStatusFailed
	payment.GatewayStatus = gatewayStatus
	if err := m.paymentStorage.UpdatePayment(ctx, payment); err != nil {
		return err
	}
	return m.orderCallbacks.OnPaymentFailed(ctx, payment.OrderID)
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	if phone == "" {
		return ""
	}
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}
	if strings.HasPrefix(phone, "0") && len(phone) == 10 {
		phone = "251" + phone[1:]
	}
	return phone
}

func mapPaymentToDTO(p *db.Payment) *dto.Payment {
	return &dto.Payment{
		ID:            p.ID,
		OrderID:       p.OrderID,
		Status:        p.Status,
		Method:        p.Method,
		Reference:     p.Reference,
		TransactionID: p.TransactionID,
		Amount:        p.Amount,
		Currency:      p.Currency,
		PhoneNumber:   p.PhoneNumber,
		Medium:        p.Medium,
		GatewayStatus: p.GatewayStatus,
		CreatedAt:     p.CreatedAt,
	}
}
