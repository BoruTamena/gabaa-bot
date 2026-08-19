package store

import (
	"context"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

func (m *storeModule) AdminListStores(ctx context.Context, filter dto.AdminStoreFilterParams) (*dto.PaginatedResponse, error) {
	stores, total, err := m.storeStorage.ListStoresAdmin(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]dto.Store, len(stores))
	for i := range stores {
		items[i] = *m.mapToDTO(&stores[i])
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}

func (m *storeModule) AdminUpdateStoreStatus(ctx context.Context, storeID int64, status string) (*dto.Store, error) {
	allowed := map[string]bool{
		constant.StoreStatusPending:  true,
		constant.StoreStatusLaunched: true,
		constant.StoreStatusBanned:   true,
	}
	if !allowed[status] {
		return nil, fmt.Errorf("invalid store status")
	}
	if err := m.storeStorage.UpdateStoreStatus(ctx, storeID, status); err != nil {
		return nil, err
	}
	return m.GetStore(ctx, storeID)
}

func (m *storeModule) ListStoreVerificationsAdmin(ctx context.Context, filter dto.AdminStoreKYCFilterParams) (*dto.PaginatedResponse, error) {
	records, total, err := m.storeKYCStorage.ListStoreKYCAdmin(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]dto.StoreKYCResponse, 0, len(records))
	for i := range records {
		store := records[i].Store
		if store.ID == 0 {
			storePtr, storeErr := m.storeStorage.GetStoreByID(ctx, records[i].StoreID)
			if storeErr != nil {
				continue
			}
			store = *storePtr
		}
		items = append(items, *m.buildKYCResponse(&store, &records[i]))
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}

func (m *storeModule) AdminUpsertStoreKYC(ctx context.Context, storeID int64, req dto.SubmitStoreKYCRequest) (*dto.StoreKYCResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("store not found")
	}

	now := time.Now()
	kyc := &db.StoreKYC{
		StoreID:                    storeID,
		TINNumber:                  req.TINNumber,
		BusinessRegistrationNumber: req.BusinessRegistrationNumber,
		TINCertificateURL:          req.TINCertificateURL,
		BusinessLicenseURL:         req.BusinessLicenseURL,
		SubmittedAt:                now,
	}
	if err := m.storeKYCStorage.UpsertStoreKYC(ctx, kyc); err != nil {
		return nil, err
	}
	if err := m.storeStorage.UpdateStoreVerificationStatus(ctx, storeID, constant.StoreVerificationPendingReview); err != nil {
		return nil, err
	}
	store.VerificationStatus = constant.StoreVerificationPendingReview
	return m.buildKYCResponse(store, kyc), nil
}
