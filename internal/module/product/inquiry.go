package product

import (
	"context"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

func (m *productModule) CreateProductInquiry(ctx context.Context, customerID, productID int64, req dto.CreateProductInquiryRequest) (*dto.ProductInquiry, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	product, err := m.productStorage.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product.Status != constant.ProductStatusPublished {
		return nil, fmt.Errorf("product is not available for inquiry")
	}
	if product.StoreID == nil {
		return nil, fmt.Errorf("product is not linked to a store")
	}

	inquiry := &db.ProductInquiry{
		ProductID:  product.ID,
		StoreID:    *product.StoreID,
		CustomerID: customerID,
		Status:     constant.InquiryStatusPending,
		Quantity:   req.Quantity,
		Note:       req.Note,
		Name:       req.Name,
		Phone:      req.Phone,
	}
	if err := m.inquiryStorage.CreateInquiry(ctx, inquiry); err != nil {
		return nil, err
	}

	created, err := m.inquiryStorage.GetInquiryByID(ctx, inquiry.ID)
	if err != nil {
		return nil, err
	}
	go m.notifyInquiryCreated(context.Background(), created)
	return m.mapInquiryToDTO(created), nil
}

func (m *productModule) ListMyInquiries(ctx context.Context, customerID int64, filter dto.ProductInquiryFilterParams) (*dto.PaginatedResponse, error) {
	inquiries, total, err := m.inquiryStorage.ListInquiriesByCustomer(ctx, customerID, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProductInquiry, len(inquiries))
	for i := range inquiries {
		items[i] = *m.mapInquiryToDTO(&inquiries[i])
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}

func (m *productModule) ListStoreInquiries(ctx context.Context, storeID int64, filter dto.ProductInquiryFilterParams) (*dto.PaginatedResponse, error) {
	inquiries, total, err := m.inquiryStorage.ListInquiriesByStore(ctx, storeID, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProductInquiry, len(inquiries))
	for i := range inquiries {
		items[i] = *m.mapInquiryToDTO(&inquiries[i])
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}

func (m *productModule) GetStoreInquiry(ctx context.Context, storeID, inquiryID int64) (*dto.ProductInquiry, error) {
	inquiry, err := m.inquiryStorage.GetInquiryByID(ctx, inquiryID)
	if err != nil {
		return nil, err
	}
	if inquiry.StoreID != storeID {
		return nil, fmt.Errorf("inquiry does not belong to your store")
	}
	return m.mapInquiryToDTO(inquiry), nil
}

func (m *productModule) ApproveProductInquiry(ctx context.Context, storeID, reviewerID, inquiryID int64, req dto.ReviewProductInquiryRequest) (*dto.ProductInquiry, error) {
	return m.reviewProductInquiry(ctx, storeID, reviewerID, inquiryID, constant.InquiryStatusApproved, req)
}

func (m *productModule) RejectProductInquiry(ctx context.Context, storeID, reviewerID, inquiryID int64, req dto.ReviewProductInquiryRequest) (*dto.ProductInquiry, error) {
	return m.reviewProductInquiry(ctx, storeID, reviewerID, inquiryID, constant.InquiryStatusRejected, req)
}

func (m *productModule) reviewProductInquiry(ctx context.Context, storeID, reviewerID, inquiryID int64, status string, req dto.ReviewProductInquiryRequest) (*dto.ProductInquiry, error) {
	inquiry, err := m.inquiryStorage.GetInquiryByID(ctx, inquiryID)
	if err != nil {
		return nil, err
	}
	if inquiry.StoreID != storeID {
		return nil, fmt.Errorf("inquiry does not belong to your store")
	}
	if inquiry.Status != constant.InquiryStatusPending {
		return nil, fmt.Errorf("only pending inquiries can be reviewed")
	}

	now := time.Now()
	inquiry.Status = status
	inquiry.ReviewedBy = &reviewerID
	inquiry.ReviewedAt = &now
	inquiry.ReviewNote = req.ReviewNote

	if err := m.inquiryStorage.UpdateInquiry(ctx, inquiry); err != nil {
		return nil, err
	}

	updated, err := m.inquiryStorage.GetInquiryByID(ctx, inquiry.ID)
	if err != nil {
		return nil, err
	}
	go m.notifyInquiryReviewed(context.Background(), updated)
	return m.mapInquiryToDTO(updated), nil
}

func (m *productModule) mapInquiryToDTO(inquiry *db.ProductInquiry) *dto.ProductInquiry {
	resp := &dto.ProductInquiry{
		ID:         inquiry.ID,
		ProductID:  inquiry.ProductID,
		StoreID:    inquiry.StoreID,
		CustomerID: inquiry.CustomerID,
		Status:     inquiry.Status,
		Quantity:   inquiry.Quantity,
		Note:       inquiry.Note,
		Name:       inquiry.Name,
		Phone:      inquiry.Phone,
		ReviewedBy: inquiry.ReviewedBy,
		ReviewedAt: inquiry.ReviewedAt,
		ReviewNote: inquiry.ReviewNote,
		CreatedAt:  inquiry.CreatedAt,
		UpdatedAt:  inquiry.UpdatedAt,
	}
	if inquiry.Product.ID != 0 {
		resp.Product = m.mapToDTO(&inquiry.Product)
	}
	return resp
}

func logInquiryAsyncFailure(action string, inquiryID int64, err error) {
	logger.Warn("product inquiry async action failed",
		zap.String("action", action),
		zap.Int64("inquiry_id", inquiryID),
		zap.Error(err),
	)
}

func (m *productModule) notifyInquiryCreated(ctx context.Context, inquiry *db.ProductInquiry) {
	store, err := m.storeStorage.GetStoreByID(ctx, inquiry.StoreID)
	if err != nil {
		logInquiryAsyncFailure("load_store_for_inquiry_create", inquiry.ID, err)
		return
	}

	merchant, err := m.userStorage.GetUserByID(ctx, store.SellerID)
	if err != nil {
		logInquiryAsyncFailure("load_merchant_for_inquiry_create", inquiry.ID, err)
		return
	}
	if merchant.TelegramUserID == nil {
		return
	}

	if err := m.tele.SendProductInquiryNotification(*merchant.TelegramUserID, *m.mapInquiryToDTO(inquiry)); err != nil {
		logInquiryAsyncFailure("send_inquiry_create_notification", inquiry.ID, err)
	}
}

func (m *productModule) notifyInquiryReviewed(ctx context.Context, inquiry *db.ProductInquiry) {
	customer, err := m.userStorage.GetUserByID(ctx, inquiry.CustomerID)
	if err != nil {
		logInquiryAsyncFailure("load_customer_for_inquiry_review", inquiry.ID, err)
		return
	}
	if customer.TelegramUserID == nil {
		return
	}

	if err := m.tele.SendProductInquiryReviewNotification(*customer.TelegramUserID, *m.mapInquiryToDTO(inquiry)); err != nil {
		logInquiryAsyncFailure("send_inquiry_review_notification", inquiry.ID, err)
	}
}
