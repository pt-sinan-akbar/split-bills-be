package models

import "time"

type BillData struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment" json:"id,omitempty"`
	// belongs to a Bill
	BillId    string     `gorm:"type:varchar(50);not null" json:"bill_id,omitempty"`
	StoreName string     `gorm:"type:varchar(50)" json:"store_name,omitempty"`
	SubTotal  float64    `gorm:"type:decimal(10,2)" json:"sub_total,omitempty"`
	Discount  float64    `gorm:"type:decimal(10,2)" json:"discount,omitempty"`
	Tax       float64    `gorm:"type:decimal(10,2)" json:"tax,omitempty"`
	Service   float64    `gorm:"type:decimal(10,2)" json:"service,omitempty"`
	Total     float64    `gorm:"type:decimal(10,2)" json:"total,omitempty"`
	Misc      string     `gorm:"type:varchar(50)" json:"misc,omitempty"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"-"`

	Bill Bill `gorm:"foreignKey:BillId " json:"bill,omitempty"`
}
