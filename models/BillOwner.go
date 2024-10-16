package models

import "time"

type BillOwner struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment"`
	// has many Bill
	Bill        []Bill    `gorm:"foreignkey:BillOwnerId;references:Id"`
	UserId      int64     `gorm:"type:int"`
	Name        string    `gorm:"type:varchar(50)"`
	Contact     string    `gorm:"type:varchar(50)"`
	BankAccount string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time `gorm:"type:datetime"`
	UpdatedAt   time.Time `gorm:"type:datetime"`
}
