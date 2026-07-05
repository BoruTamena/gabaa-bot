package store

import (
	"context"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type storeModule struct {
	storeStorage    storage.StoreStorage
	storeKYCStorage storage.StoreKYCStorage
	userStorage     storage.UserStorage
	tele            platform.Telegram
}

func NewStoreModule(
	sStorage storage.StoreStorage,
	kycStorage storage.StoreKYCStorage,
	uStorage storage.UserStorage,
	tele platform.Telegram,
) module.StoreModule {
	return &storeModule{
		storeStorage:    sStorage,
		storeKYCStorage: kycStorage,
		userStorage:     uStorage,
		tele:            tele,
	}
}

func (m *storeModule) CreateStore(ctx context.Context, userID int64, req dto.CreateStoreRequest) (*dto.Store, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.TelegramChatID != 0 {
		user, err := m.userStorage.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}

		if user.TelegramUserID == nil {
			return nil, fmt.Errorf("user has no telegram id linked")
		}

		isAdmin, err := m.tele.IsChatAdmin(req.TelegramChatID, *user.TelegramUserID)
		if err != nil || !isAdmin {
			logger.Warn("user attempted to create store for chat without admin rights",
				zap.Int64("user_id", userID), zap.Int64("chat_id", req.TelegramChatID))
			return nil, fmt.Errorf("you must be an admin of the chat to create a store")
		}
	}

	store := &db.Store{
		SellerID:           userID,
		TelegramChatID:     req.TelegramChatID,
		VerificationStatus: constant.StoreVerificationUnverified,
		Name:               req.Name,
		Category:           req.Category,
		Description:        req.Description,
		LogoImage:          req.LogoImage,
		CoverImage:         req.CoverImage,
		Phone:              req.Phone,
		Email:              req.Email,
		Location:           req.Location,
	}

	if err := m.storeStorage.CreateStore(ctx, store); err != nil {
		logger.Error("failed to create store", zap.Error(err), zap.Int64("seller_id", userID))
		return nil, err
	}

	logger.Info("store created successfully", zap.Int64("store_id", store.ID), zap.Int64("seller_id", userID))

	return m.mapToDTO(store), nil
}

func (m *storeModule) GetAdminDashboard(ctx context.Context, userID int64, chatID int64) (string, *dto.Store, error) {
	store, err := m.storeStorage.GetStoreByChatID(ctx, chatID)

	user, userErr := m.userStorage.GetUserByID(ctx, userID)
	tgUserID := int64(0)
	if userErr == nil && user.TelegramUserID != nil {
		tgUserID = *user.TelegramUserID
	}

	if err == nil {
		if tgUserID != 0 {
			isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUserID)
			if isAdmin {
				logger.Info("dashboard: manage (group admin)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
				return "manage", m.mapToDTO(store), nil
			}
		}
		logger.Info("dashboard: storefront", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
		return "storefront", m.mapToDTO(store), nil
	}

	stores, _ := m.storeStorage.GetStoresBySellerID(ctx, userID)
	if len(stores) > 0 {
		if chatID > 0 {
			logger.Info("dashboard: manage (private chat fallback)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
			return "manage", m.mapToDTO(&stores[0]), nil
		}
	}

	if tgUserID != 0 {
		isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUserID)
		if isAdmin {
			logger.Info("dashboard: setup", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
			return "setup", nil, nil
		}
	}

	logger.Info("dashboard: storefront (no store)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
	return "storefront", nil, nil
}

func (m *storeModule) GetStore(ctx context.Context, id int64) (*dto.Store, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.mapToDTO(store), nil
}

func (m *storeModule) GetStoreStatus(ctx context.Context, id int64) (string, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return "", err
	}
	return store.Status, nil
}

func (m *storeModule) UpdateStore(ctx context.Context, id int64, req dto.UpdateStoreRequest) (*dto.Store, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Category != "" {
		store.Category = req.Category
	}
	if req.Description != "" {
		store.Description = req.Description
	}
	if req.LogoImage != "" {
		store.LogoImage = req.LogoImage
	}
	if req.CoverImage != "" {
		store.CoverImage = req.CoverImage
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}
	if req.Email != "" {
		store.Email = req.Email
	}
	if req.Location != "" {
		store.Location = req.Location
	}

	if err := m.storeStorage.UpdateStore(ctx, store); err != nil {
		logger.Error("failed to update store", zap.Error(err), zap.Int64("store_id", id))
		return nil, err
	}

	logger.Info("store updated successfully", zap.Int64("store_id", id))

	return m.mapToDTO(store), nil
}

func (m *storeModule) SubmitStoreKYC(ctx context.Context, storeID, sellerID int64, req dto.SubmitStoreKYCRequest) (*dto.StoreKYCResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("store not found")
	}
	if store.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized to submit KYC for this store")
	}

	switch store.VerificationStatus {
	case constant.StoreVerificationPendingReview:
		return nil, fmt.Errorf("KYC submission is already under review")
	case constant.StoreVerificationVerified:
		return nil, fmt.Errorf("store is already verified")
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

func (m *storeModule) GetStoreKYC(ctx context.Context, storeID, sellerID int64) (*dto.StoreKYCResponse, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("store not found")
	}
	if store.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized to view KYC for this store")
	}

	kyc, err := m.storeKYCStorage.GetStoreKYCByStoreID(ctx, storeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &dto.StoreKYCResponse{
				StoreID:            storeID,
				StoreName:          store.Name,
				VerificationStatus: store.VerificationStatus,
			}, nil
		}
		return nil, err
	}

	return m.buildKYCResponse(store, kyc), nil
}

func (m *storeModule) ListStoreVerifications(ctx context.Context, status string) ([]dto.StoreKYCResponse, error) {
	if status == "" {
		status = constant.StoreVerificationPendingReview
	}

	records, err := m.storeKYCStorage.ListStoreKYCByVerificationStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.StoreKYCResponse, 0, len(records))
	for i := range records {
		store := records[i].Store
		if store.ID == 0 {
			storePtr, storeErr := m.storeStorage.GetStoreByID(ctx, records[i].StoreID)
			if storeErr != nil {
				continue
			}
			store = *storePtr
		}
		responses = append(responses, *m.buildKYCResponse(&store, &records[i]))
	}
	return responses, nil
}

func (m *storeModule) ApproveStoreKYC(ctx context.Context, storeID int64) (*dto.StoreKYCResponse, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("store not found")
	}
	if store.VerificationStatus != constant.StoreVerificationPendingReview {
		return nil, fmt.Errorf("store KYC is not pending review")
	}

	now := time.Now()
	if err := m.storeKYCStorage.UpdateStoreKYCReview(ctx, storeID, "", now); err != nil {
		return nil, err
	}
	if err := m.storeStorage.UpdateStoreVerificationStatus(ctx, storeID, constant.StoreVerificationVerified); err != nil {
		return nil, err
	}

	store.VerificationStatus = constant.StoreVerificationVerified
	kyc, _ := m.storeKYCStorage.GetStoreKYCByStoreID(ctx, storeID)
	return m.buildKYCResponse(store, kyc), nil
}

func (m *storeModule) RejectStoreKYC(ctx context.Context, storeID int64, req dto.RejectStoreKYCRequest) (*dto.StoreKYCResponse, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("store not found")
	}
	if store.VerificationStatus != constant.StoreVerificationPendingReview {
		return nil, fmt.Errorf("store KYC is not pending review")
	}

	now := time.Now()
	if err := m.storeKYCStorage.UpdateStoreKYCReview(ctx, storeID, req.ReviewNote, now); err != nil {
		return nil, err
	}
	if err := m.storeStorage.UpdateStoreVerificationStatus(ctx, storeID, constant.StoreVerificationRejected); err != nil {
		return nil, err
	}

	store.VerificationStatus = constant.StoreVerificationRejected
	kyc, _ := m.storeKYCStorage.GetStoreKYCByStoreID(ctx, storeID)
	return m.buildKYCResponse(store, kyc), nil
}

func (m *storeModule) IsStoreVerified(ctx context.Context, storeID int64) (bool, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return false, err
	}
	return store.VerificationStatus == constant.StoreVerificationVerified, nil
}

func (m *storeModule) buildKYCResponse(store *db.Store, kyc *db.StoreKYC) *dto.StoreKYCResponse {
	resp := &dto.StoreKYCResponse{
		StoreID:            store.ID,
		StoreName:          store.Name,
		VerificationStatus: store.VerificationStatus,
	}
	if kyc == nil {
		return resp
	}

	submittedAt := kyc.SubmittedAt
	resp.TINNumber = kyc.TINNumber
	resp.BusinessRegistrationNumber = kyc.BusinessRegistrationNumber
	resp.TINCertificateURL = kyc.TINCertificateURL
	resp.BusinessLicenseURL = kyc.BusinessLicenseURL
	resp.ReviewNote = kyc.ReviewNote
	resp.SubmittedAt = &submittedAt
	resp.ReviewedAt = kyc.ReviewedAt
	return resp
}

func (m *storeModule) mapToDTO(store *db.Store) *dto.Store {
	return &dto.Store{
		ID:                 store.ID,
		SellerID:           store.SellerID,
		TelegramChatID:     store.TelegramChatID,
		TelegramChatTitle:  store.TelegramChatTitle,
		Status:             store.Status,
		VerificationStatus: store.VerificationStatus,
		Name:               store.Name,
		Category:           store.Category,
		Description:        store.Description,
		LogoImage:          store.LogoImage,
		CoverImage:         store.CoverImage,
		Phone:              store.Phone,
		Email:              store.Email,
		Location:           store.Location,
	}
}
