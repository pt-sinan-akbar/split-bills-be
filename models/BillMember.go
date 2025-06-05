package models

import "time"

type BillMember struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment" json:"id,omitempty"`
	// belongs to a Bill
	BillId    string     `gorm:"type:varchar(50);not null" json:"bill_id,omitempty"`
	Name      string     `gorm:"type:varchar(50)" json:"name,omitempty"`
	PriceOwe  *float64   `gorm:"type:decimal(10,2)" json:"price_owe"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"-"`

	Bill *Bill `gorm:"foreignKey:BillId" json:"bill,omitempty"`
	// many to many
	BillItem []BillItem `gorm:"many2many:bill_member_items;" json:"bill_item"`
}
