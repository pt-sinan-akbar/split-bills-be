package models

import (
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID string `gorm:"type:varchar(50);primary_key"`
	// has many BillData
	BillData []BillData `gorm:"foreignkey:BillId;references:Id"`
	// belongs to a BillOwner
	BillOwnerId int64          `gorm:"type:int"`
	BillOwner   BillOwner      `gorm:"foreignKey:BillOwnerId;references:Id"`
	Name        string         `gorm:"type:varchar(50)"`
	RawImage    string         `gorm:"type:varchar(50)"`
	CreatedAt   time.Time      `gorm:"type:datetime"`
	UpdatedAt   time.Time      `gorm:"type:datetime"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime"`
}
