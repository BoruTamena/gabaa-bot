package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/lakipay"
	"go.uber.org/zap"
)

type walletModule struct {
	walletStorage     storage.WalletStorage
	withdrawalStorage storage.WithdrawalStorage
	storeStorage      storage.StoreStorage
	lakipay           platform.LakiPay
}

func NewWalletModule(
	wStorage storage.WalletStorage,
	withdrawalStorage storage.WithdrawalStorage,
	storeStorage storage.StoreStorage,
	lakipayClient platform.LakiPay,
) module.WalletModule {
	return &walletModule{
		walletStorage:     wStorage,
		withdrawalStorage: withdrawalStorage,
		storeStorage:      storeStorage,
		lakipay:           lakipayClient,
	}
}

func (m *walletModule) GetWalletSummary(ctx context.Context, storeID int64) (*dto.Wallet, error) {
	wallet, err := m.walletStorage.GetOrCreateWallet(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return mapWalletToDTO(wallet), nil
}

func (m *walletModule) RequestWithdrawal(ctx context.Context, storeID int64, req dto.WithdrawalRequest) (*dto.Withdrawal, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		logger.Error("withdrawal: store not found", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, fmt.Errorf("store not found: %w", err)
	}
	// if store.VerificationStatus != constant.StoreVerificationVerified {
	// 	logger.Error("withdrawal: store not verified",
	// 		zap.Int64("store_id", storeID),
	// 		zap.String("verification_status", store.VerificationStatus),
	// 	)
	// 	return nil, fmt.Errorf("store is not verified for withdrawals")
	// }
	logger.Info("the withdrawal request", zap.Int64("store_id", store.ID), zap.Float64("amount", req.Amount), zap.String("phone_number", req.PhoneNumber), zap.String("medium", req.Medium))

	phone := normalizePhone(req.PhoneNumber)
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	wallet, err := m.walletStorage.GetOrCreateWallet(ctx, storeID)
	if err != nil {
		logger.Error("withdrawal: failed to load wallet", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}
	if wallet.AvailableBalance < req.Amount {
		logger.Error("withdrawal: insufficient balance",
			zap.Int64("store_id", storeID),
			zap.Float64("requested", req.Amount),
			zap.Float64("available", wallet.AvailableBalance),
		)
		return nil, fmt.Errorf("insufficient available balance")
	}

	withdrawal := &db.Withdrawal{
		StoreID:     storeID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		PhoneNumber: phone,
		Medium:      req.Medium,
		Status:      constant.WithdrawalStatusInitiated,
	}
	if err := m.withdrawalStorage.CreateWithdrawal(ctx, withdrawal); err != nil {
		logger.Error("withdrawal: failed to create record", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}

	reference := fmt.Sprintf("WITHDRAW-%d-%d", storeID, withdrawal.ID)
	withdrawal.Reference = reference
	if err := m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal); err != nil {
		logger.Error("withdrawal: failed to set reference", zap.Error(err), zap.Int64("withdrawal_id", withdrawal.ID))
		return nil, err
	}

	if err := m.walletStorage.LockForWithdrawal(ctx, storeID, req.Amount); err != nil {
		withdrawal.Status = constant.WithdrawalStatusFailed
		_ = m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal)
		logger.Error("withdrawal: failed to lock funds", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, fmt.Errorf("failed to lock funds: %w", err)
	}

	callbackURL := lakipay.ResolveCallbackURL()
	if callbackURL == "" {
		_ = m.walletStorage.UnlockWithdrawal(ctx, storeID, req.Amount)
		return nil, fmt.Errorf("lakipay callback URL is not configured (set LAKIPAY_CALLBACK_URL)")
	}

	logger.Info("initiating lakipay withdrawal",
		zap.Int64("store_id", storeID),
		zap.Int64("withdrawal_id", withdrawal.ID),
		zap.String("reference", reference),
		zap.String("medium", req.Medium),
		zap.String("callback_url", callbackURL),
	)

	resp, err := m.lakipay.InitiateWithdrawal(ctx, lakipay.WithdrawalRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		PhoneNumber: phone,
		Medium:      req.Medium,
		Reference:   reference,
		CallbackURL: callbackURL,
	})
	if err != nil {
		_ = m.walletStorage.UnlockWithdrawal(ctx, storeID, req.Amount)
		withdrawal.Status = constant.WithdrawalStatusFailed
		withdrawal.GatewayStatus = constant.GatewayPaymentStatusFailed
		_ = m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal)
		logger.Error("withdrawal: lakipay initiation failed",
			zap.Error(err),
			zap.Int64("store_id", storeID),
			zap.Int64("withdrawal_id", withdrawal.ID),
		)
		return nil, fmt.Errorf("withdrawal initiation failed: %w", err)
	}

	gatewayResp, _ := json.Marshal(resp)
	withdrawal.Status = constant.WithdrawalStatusPending
	if txnID := resp.TransactionID(); txnID != "" {
		withdrawal.TransactionID = &txnID
	}
	withdrawal.GatewayStatus = constant.ParseGatewayPaymentStatus(resp.GatewayStatus())
	withdrawal.GatewayResponse = gatewayResp
	if err := m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal); err != nil {
		logger.Error("withdrawal: failed to update after lakipay", zap.Error(err), zap.Int64("withdrawal_id", withdrawal.ID))
		return nil, err
	}

	logger.Info("withdrawal initiated successfully",
		zap.Int64("store_id", storeID),
		zap.Int64("withdrawal_id", withdrawal.ID),
		zap.String("status", string(withdrawal.Status)),
	)

	return mapWithdrawalToDTO(withdrawal), nil
}

func (m *walletModule) ListWithdrawals(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	limit := params.GetLimit()
	offset := params.GetOffset()

	withdrawals, total, err := m.withdrawalStorage.ListWithdrawalsByStoreID(ctx, storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	dtoList := make([]dto.Withdrawal, len(withdrawals))
	for i, w := range withdrawals {
		dtoList[i] = *mapWithdrawalToDTO(&w)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoList,
	}, nil
}

func (m *walletModule) GetMyStoreWithdrawal(ctx context.Context, storeID, withdrawalID int64) (*dto.Withdrawal, error) {
	withdrawal, err := m.withdrawalStorage.GetWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		return nil, err
	}
	if withdrawal.StoreID != storeID {
		return nil, fmt.Errorf("withdrawal does not belong to your store")
	}
	return mapWithdrawalToDTO(withdrawal), nil
}

func mapWalletToDTO(wallet *db.Wallet) *dto.Wallet {
	return &dto.Wallet{
		ID:               wallet.ID,
		StoreID:          wallet.StoreID,
		Currency:         wallet.Currency,
		PendingBalance:   wallet.PendingBalance,
		AvailableBalance: wallet.AvailableBalance,
		LockedBalance:    wallet.LockedBalance,
		TotalEarned:      wallet.TotalEarned,
		TotalWithdrawn:   wallet.TotalWithdrawn,
	}
}

func mapWithdrawalToDTO(w *db.Withdrawal) *dto.Withdrawal {
	return &dto.Withdrawal{
		ID:            w.ID,
		StoreID:       w.StoreID,
		Amount:        w.Amount,
		Currency:      w.Currency,
		PhoneNumber:   w.PhoneNumber,
		Medium:        w.Medium,
		Reference:     w.Reference,
		TransactionID: w.TransactionID,
		Status:        w.Status,
		GatewayStatus: w.GatewayStatus,
		CreatedAt:     w.CreatedAt,
	}
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
