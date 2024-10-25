package models

import (
	"time"
)

type Bill struct {
	ID string `gorm:"type:varchar(50);primary_key" json:"id,omitempty"`
	// has many BillData
	BillData []BillData `gorm:"foreignkey:BillId" json:"-"`
	// belongs to a BillOwner
	BillOwnerId *int64     `gorm:"type:int" json:"bill_owner_id,omitempty"`
	BillOwner   *BillOwner `gorm:"foreignKey:BillOwnerId" json:"-"`
	Name        string     `gorm:"type:varchar(50)" json:"name,omitempty"`
	RawImage    string     `gorm:"type:varchar(50)" json:"raw_image,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"-"`
	UpdatedAt   time.Time  `gorm:"type:timestamp" json:"-"`
	DeletedAt   *time.Time `gorm:"type:timestamp" json:"-"`	
}
