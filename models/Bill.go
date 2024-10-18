package models

import (
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID string `gorm:"type:varchar(50);primary_key"`
	// has many BillData
	BillData []BillData `gorm:"foreignkey:BillId"`
	// belongs to a BillOwner
	BillOwnerId int64          `gorm:"type:int;not null"`
	BillOwner   BillOwner      `gorm:"foreignKey:BillOwnerId"`
	Name        string         `gorm:"type:varchar(50)"`
	RawImage    string         `gorm:"type:varchar(50)"`
	CreatedAt   time.Time      `gorm:"type:timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp" swaggerignore:"true"`
}
