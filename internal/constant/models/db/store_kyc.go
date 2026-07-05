package db

import "time"

type StoreKYC struct {
	BaseModel
	StoreID                    int64      `gorm:"column:store_id;not null;uniqueIndex" json:"store_id"`
	TINNumber                  string     `gorm:"column:tin_number;not null" json:"tin_number"`
	BusinessRegistrationNumber string     `gorm:"column:business_registration_number;not null" json:"business_registration_number"`
	TINCertificateURL          string     `gorm:"column:tin_certificate_url;not null" json:"tin_certificate_url"`
	BusinessLicenseURL         string     `gorm:"column:business_license_url;not null" json:"business_license_url"`
	ReviewNote                 string     `gorm:"column:review_note" json:"review_note"`
	SubmittedAt                time.Time  `gorm:"column:submitted_at;not null" json:"submitted_at"`
	ReviewedAt                 *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	Store                      Store      `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
}

func (StoreKYC) TableName() string {
	return "store_kyc"
}
