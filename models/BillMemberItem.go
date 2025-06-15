package models

import "time"

type BillMemberItem struct {
	BillId       string      `gorm:"type:varchar(50);" json:"bill_id,omitempty"`
	Bill         *Bill       `gorm:"foreignKey:BillId" json:"bill,omitempty"`
	BillItemId   int         `gorm:"type:int;primaryKey" json:"bill_item_id,omitempty"`
	BillItem     *BillItem   `gorm:"foreignKey:BillItemId;references:ID" json:"bill_item,omitempty"`
	BillMemberId int         `gorm:"type:int;primaryKey;" json:"bill_member_id,omitempty"`
	BillMember   *BillMember `gorm:"foreignKey:BillMemberId;references:ID" json:"bill_member,omitempty"`
	Quantity     *int64      `gorm:"type:int" json:"qty"`
	DeletedAt    *time.Time  `gorm:"type:timestamp" json:"deleted_at"`
}
