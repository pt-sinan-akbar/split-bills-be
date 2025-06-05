package manager

import (
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"time"
)

type BillMemberManager struct {
	DB *gorm.DB
}

func NewBillMemberManager(DB *gorm.DB) BillMemberManager {
	return BillMemberManager{DB}
}

func (bmm BillMemberManager) GetAll() ([]models.BillMember, error) {
	var obj []models.BillMember
	result := bmm.DB.Where("deleted_at IS NULL").Preload("Bill").Find(&obj)
	return obj, result.Error
}

func (bmm BillMemberManager) GetByID(id int) (models.BillMember, error) {
	var obj models.BillMember
	result := bmm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").Preload("BillItem").First(&obj)
	return obj, result.Error
}

func (bmm BillMemberManager) CreateAsync(obj models.BillMember) error {
	result := bmm.DB.Create(&obj)
	return result.Error
}

func (bmm BillMemberManager) DeleteAsync(id int) error {
	tx := bmm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bmm.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	obj.UpdatedAt = time.Now()
	obj.DeletedAt = &obj.UpdatedAt

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (bmm BillMemberManager) EditAsync(id int, obj models.BillMember) error {
	tx := bmm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	oldObj, err := bmm.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&oldObj).Updates(map[string]interface{}{
		"bill_id":    obj.BillId,
		"name":       obj.Name,
		"price_owe":  obj.PriceOwe,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
