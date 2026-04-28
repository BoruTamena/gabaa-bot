package wallet

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type walletModule struct {
	walletStorage storage.WalletStorage
}

func NewWalletModule(wStorage storage.WalletStorage) module.WalletModule {
	return &walletModule{
		walletStorage: wStorage,
	}
}

func (m *walletModule) GetBalance(ctx context.Context, storeID int64) (float64, error) {
	wallet, err := m.walletStorage.GetWalletByStoreID(ctx, storeID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (m *walletModule) CreditWallet(ctx context.Context, storeID int64, amount float64) error {
	err := m.walletStorage.UpdateWalletBalance(ctx, storeID, amount)
	if err != nil {
		logger.Error("failed to credit wallet", zap.Error(err), zap.Int64("store_id", storeID), zap.Float64("amount", amount))
	} else {
		logger.Info("wallet credited successfully", zap.Int64("store_id", storeID), zap.Float64("amount", amount))
	}
	return err
}
