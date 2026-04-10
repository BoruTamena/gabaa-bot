package db

type Category struct {
	BaseModel
	StoreID int64  `gorm:"column:store_id" json:"store_id"`
	Name    string `gorm:"column:name;not null" json:"name"`
}
