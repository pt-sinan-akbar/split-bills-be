package models

import "time"

type BillOwner struct {
	ID int64 `gorm:"type:int;primary_key;auto_increment" json:"id,omitempty"`
	// has many Bill
	Bill        []Bill    `gorm:"foreignkey:BillOwnerId" json:"-"`
	UserId      int64     `gorm:"type:int" json:"user_id,omitempty"`
	Name        string    `gorm:"type:varchar(50)" json:"name,omitempty"`
	Contact     string    `gorm:"type:varchar(50)" json:"contact,omitempty"`
	BankAccount string    `gorm:"type:varchar(50)" json:"bank_account,omitempty"`
	CreatedAt   time.Time `gorm:"type:timestamp" json:"-"`
	UpdatedAt   time.Time `gorm:"type:timestamp" json:"-"`
	DeletedAt   *time.Time `gorm:"type:timestamp" json:"-"`
}
