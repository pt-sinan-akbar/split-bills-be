package models

import "time"

type BillItem struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment"`
	// belongs to a Bill
	BillId string `gorm:"type:varchar(50)"`
	Bill   Bill   `gorm:"foreignKey:BillId;references:Id"`
	// many to many
	BillMember []*BillMember `gorm:"many2many:bill_member_items;"`
	Name       string        `gorm:"type:varchar(50)"`
	Qty        int64         `gorm:"type:int"`
	Price      float64       `gorm:"type:decimal(10,2)"`
	Tax        float64       `gorm:"type:decimal(10,2)"`
	Service    float64       `gorm:"type:decimal(10,2)"`
	Discount   float64       `gorm:"type:decimal(10,2)"`
	CreatedAt  time.Time     `gorm:"type:datetime"`
	UpdatedAt  time.Time     `gorm:"type:datetime"`
}
