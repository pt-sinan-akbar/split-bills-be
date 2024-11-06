package manager

import (
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"time"
)

type BillDataManager struct {
	DB *gorm.DB
}

func NewBillDataManager(DB *gorm.DB) BillDataManager {
	return BillDataManager{DB}
}

func (bdm BillDataManager) GetAll() ([]models.BillData, error) {
	var obj []models.BillData
	result := bdm.DB.Where("deleted_at IS NULL").Preload("Bill").Find(&obj)
	return obj, result.Error
}

func (bdm BillDataManager) GetByID(id int) (models.BillData, error) {
	var obj models.BillData
	result := bdm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").First(&obj)
	return obj, result.Error
}

func (bdm BillDataManager) CreateAsync(billData *models.BillData) error {
	result := bdm.DB.Create(&billData)
	return result.Error
}

func (bdm BillDataManager) DeleteAsync(id int) error {
	obj, _ := bdm.GetByID(id)
	obj.UpdatedAt = time.Now()
	obj.DeletedAt = &obj.UpdatedAt
	result := bdm.DB.Save(&obj)
	return result.Error
}

func (bdm BillDataManager) EditAsync(id int, obj *models.BillData) error {
	old, _ := bdm.GetByID(id)

	tx := bdm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	updateData := map[string]interface{}{
		"BillId":    obj.BillId,
		"StoreName": obj.StoreName,
		"SubTotal":  obj.SubTotal,
		"Discount":  obj.Discount,
		"Tax":       obj.Tax,
		"Service":   obj.Service,
		"Total":     obj.Total,
		"Misc":      obj.Misc,
		"UpdatedAt": time.Now(),
	}

	if err := tx.Model(&old).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
