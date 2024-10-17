package models

import "time"

type BillData struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment"`
	// belongs to a Bill
	BillId    string    `gorm:"type:varchar(50);not null"`
	Bill      Bill      `gorm:"foreignKey:BillId"`
	StoreName string    `gorm:"type:varchar(50)"`
	SubTotal  float64   `gorm:"type:decimal(10,2)"`
	Discount  float64   `gorm:"type:decimal(10,2)"`
	Tax       float64   `gorm:"type:decimal(10,2)"`
	Service   float64   `gorm:"type:decimal(10,2)"`
	Total     float64   `gorm:"type:decimal(10,2)"`
	Misc      string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}
