package dto

type PaginationParams struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *PaginationParams) GetLimit() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

func (p *PaginationParams) GetOffset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.GetLimit()
}

func (p *PaginationParams) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

type PaginatedResponse struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	Data        interface{} `json:"data"`
}

func NewPaginatedResponse(data interface{}, total int64, params PaginationParams) *PaginatedResponse {
	page := params.GetPage()
	pageSize := params.GetLimit()
	return &PaginatedResponse{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		HasNext:     int64(page*pageSize) < total,
		HasPrevious: page > 1,
		Data:        data,
	}
}
