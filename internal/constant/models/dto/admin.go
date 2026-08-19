package dto

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminStoreFilterParams struct {
	PaginationParams
	Status             string `form:"status"`
	VerificationStatus string `form:"verification_status"`
	Category           string `form:"category"`
	Query              string `form:"query"`
}

type AdminUpdateStoreStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminStoreKYCFilterParams struct {
	PaginationParams
	Status string `form:"status"`
	Query  string `form:"query"`
}

type AdminOrderFilterParams struct {
	PaginationParams
	StoreID int64  `form:"store_id"`
	OrderID *int64 `form:"order_id"`
	Status  string `form:"status"`
	Query   string `form:"query"`
}

type AdminUserFilterParams struct {
	PaginationParams
	Role  string `form:"role"`
	Query string `form:"query"`
}

type AdminStoryFilterParams struct {
	PaginationParams
	IsActive *bool  `form:"is_active"`
	Type     string `form:"type"`
	Query    string `form:"query"`
}

type AdminDeliveryFilterParams struct {
	PaginationParams
	Status string `form:"status"`
	Query  string `form:"query"`
}

type AdminWithdrawalFilterParams struct {
	PaginationParams
	Status string `form:"status"`
}

type AdminTransactionFilterParams struct {
	PaginationParams
	Status string `form:"status"`
	Medium string `form:"medium"`
	Query  string `form:"query"`
}
