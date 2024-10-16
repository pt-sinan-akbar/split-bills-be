package models

import "time"

type BillMember struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment"`
	// belongs to a Bill
	BillId string `gorm:"type:varchar(50)"`
	// many to many
	BillItem  []*BillItem `gorm:"many2many:bill_member_items;"`
	Bill      Bill        `gorm:"foreignKey:BillId;references:Id"`
	Name      string      `gorm:"type:varchar(50)"`
	PriceOwe  float64     `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time   `gorm:"type:datetime"`
	UpdatedAt time.Time   `gorm:"type:datetime"`
}
