package models

import (
	"time"
)

type Bill struct {
	ID          string     `gorm:"type:varchar(50);primary_key" json:"id,omitempty"`
	BillOwnerId *int64     `gorm:"type:int" json:"bill_owner_id,omitempty"`
	Name        string     `gorm:"type:varchar(50)" json:"name,omitempty"`
	RawImage    *string    `gorm:"type:varchar(50)" json:"raw_image"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at,omitempty"`
	UpdatedAt   time.Time  `gorm:"type:timestamp" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"type:timestamp" json:"-"`

	// has one BillData
	BillData *BillData `gorm:"foreignkey:BillId" json:"bill_data,omitempty"`
	// belongs to a BillOwner
	BillOwner *BillOwner `gorm:"foreignKey:BillOwnerId" json:"bill_owner"`
	// has many BillItem
	BillItem []BillItem `gorm:"foreignkey:BillId" json:"bill_item"`
	// has many BillMember
	BillMember []BillMember `gorm:"foreignkey:BillId" json:"bill_member"`
	// has many BillMemberItems
	BillMemberItem []BillMemberItem `gorm:"foreignkey:BillId" json:"bill_member_item"`
}
