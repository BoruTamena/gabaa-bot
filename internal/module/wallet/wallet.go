package wallet

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
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
	return m.walletStorage.UpdateWalletBalance(ctx, storeID, amount)
}
