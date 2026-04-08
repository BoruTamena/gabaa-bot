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

type PaginatedResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
