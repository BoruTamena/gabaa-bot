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
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/lakipay"
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
		return nil, fmt.Errorf("store not found: %w", err)
	}
	if store.VerificationStatus != constant.StoreVerificationVerified {
		return nil, fmt.Errorf("store is not verified for withdrawals")
	}

	phone := normalizePhone(req.PhoneNumber)
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	wallet, err := m.walletStorage.GetOrCreateWallet(ctx, storeID)
	if err != nil {
		return nil, err
	}
	if wallet.AvailableBalance < req.Amount {
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
		return nil, err
	}

	reference := fmt.Sprintf("WITHDRAW-%d-%d", storeID, withdrawal.ID)
	withdrawal.Reference = reference
	if err := m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal); err != nil {
		return nil, err
	}

	if err := m.walletStorage.LockForWithdrawal(ctx, storeID, req.Amount); err != nil {
		withdrawal.Status = constant.WithdrawalStatusFailed
		_ = m.withdrawalStorage.UpdateWithdrawal(ctx, withdrawal)
		return nil, fmt.Errorf("failed to lock funds: %w", err)
	}

	callbackURL := lakipay.ResolveCallbackURL()
	if callbackURL == "" {
		return nil, fmt.Errorf("lakipay callback URL is not configured (set LAKIPAY_CALLBACK_URL)")
	}

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
		return nil, err
	}

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
