package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BillMemberItemManager struct {
	DB *gorm.DB
}

func NewBillMemberItemManager(DB *gorm.DB) BillMemberItemManager {
	return BillMemberItemManager{DB}
}

func (bmim BillMemberItemManager) internalGet(billId *string, memberId *int, itemId *int) ([]models.BillMemberItem, error) {
	queryParts := []string{"deleted_at IS NULL"}
	if billId != nil {
		queryParts = append(queryParts, "bill_id = '"+*billId+"'")
	}
	if memberId != nil {
		queryParts = append(queryParts, "bill_member_id = "+fmt.Sprint(*memberId))
	}
	if itemId != nil {
		queryParts = append(queryParts, "bill_item_id = "+fmt.Sprint(*itemId))
	}

	query := strings.Join(queryParts, " AND ")
	var obj []models.BillMemberItem
	result := bmim.DB.Where(query).Preload("BillMember").Find(&obj)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get bill member items: %v", result.Error)
	}
	return obj, nil
}

func (bmim BillMemberItemManager) GetAll() ([]models.BillMemberItem, error) {
	return bmim.internalGet(nil, nil, nil)
}

func (bmim BillMemberItemManager) GetByIds(billId string, memberId int, itemId int) (models.BillMemberItem, error) {
	result, err := bmim.internalGet(&billId, &memberId, &itemId)
	if err != nil {
		return models.BillMemberItem{}, fmt.Errorf("failed to get bill member item: %v", err)
	}
	return result[0], nil
}

func (bmim BillMemberItemManager) GetByItemId(itemId int) ([]models.BillMemberItem, error) {
	return bmim.internalGet(nil, nil, &itemId)
}

func (bmim BillMemberItemManager) GetByMemberId(memberId int) ([]models.BillMemberItem, error) {
	return bmim.internalGet(nil, &memberId, nil)
}

func (bmim BillMemberItemManager) DeleteModel(obj models.BillMemberItem) error {
	tx := bmim.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	now := time.Now()
	obj.DeletedAt = &now

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (bmim BillMemberItemManager) DeleteByItemId(itemId int) error {
	memberItems, err := bmim.GetByItemId(itemId)
	if err != nil {
		return fmt.Errorf("failed to get member items: %v", err)
	}
	for _, memberItem := range memberItems {
		err = bmim.DeleteModel(memberItem)
		if err != nil {
			return fmt.Errorf("failed to delete member item: %v", err)
		}
	}
	return nil
}

func (bmim BillMemberItemManager) DeleteByMemberId(memberId int) error {
	memberItems, err := bmim.GetByMemberId(memberId)
	if err != nil {
		return fmt.Errorf("failed to get member items: %v", err)
	}
	for _, memberItem := range memberItems {
		err = bmim.DeleteModel(memberItem)
		if err != nil {
			return fmt.Errorf("failed to delete member item: %v", err)
		}
	}
	return nil
}
