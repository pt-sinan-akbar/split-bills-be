package models

import "time"

type BillItem struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment" json:"id,omitempty"`
	// belongs to a Bill
	BillId    string     `gorm:"type:varchar(50);not null" json:"bill_id,omitempty"`
	Name      string     `gorm:"type:varchar(50)" json:"name,omitempty"`
	Qty       int64      `gorm:"type:int" json:"qty,omitempty"`
	Price     float64    `gorm:"type:decimal(10,2)" json:"price,omitempty"`
	Subtotal  float64    `gorm:"type:decimal(10,2)" json:"subtotal,omitempty"`
	Tax       float64    `gorm:"type:decimal(10,2)" json:"tax,omitempty"`
	Service   float64    `gorm:"type:decimal(10,2)" json:"service,omitempty"`
	Discount  float64    `gorm:"type:decimal(10,2)" json:"discount,omitempty"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"-"`

	Bill Bill `gorm:"foreignKey:BillId" json:"bill,omitempty"`
	// many to many
	BillMember []*BillMember `gorm:"many2many:bill_member_items;" json:"bill_member,omitempty"`
}
