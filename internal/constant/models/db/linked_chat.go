package db

type LinkedChat struct {
	BaseModel
	StoreID   int64  `gorm:"column:store_id;not null" json:"store_id"`
	ChatID    int64  `gorm:"column:chat_id;uniqueIndex;not null" json:"chat_id"`
	ChatType  string `gorm:"column:chat_type" json:"chat_type"`
	Title     string `gorm:"column:title" json:"title"`
	Status    string `gorm:"column:status;default:'active'" json:"status"`
	Store     Store  `gorm:"foreignKey:StoreID;references:ID" json:"store"`
}
