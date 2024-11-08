package manager

import (
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"time"
)

type BillOwnerManager struct {
	DB *gorm.DB
}

func NewBillOwnerManager(DB *gorm.DB) BillOwnerManager {
	return BillOwnerManager{DB}
}

func (bom BillOwnerManager) GetAll() ([]models.BillOwner, error) {
	var obj []models.BillOwner
	result := bom.DB.Where("deleted_at IS NULL").Find(&obj)
	return obj, result.Error
}

func (bom BillOwnerManager) GetById(id int) (models.BillOwner, error) {
	var obj models.BillOwner
	result := bom.DB.Where("id = ? AND deleted_at IS NULL", id).First(&obj)
	return obj, result.Error
}

func (bom BillOwnerManager) CreateAsync(obj models.BillOwner) error {
	result := bom.DB.Create(&obj)
	return result.Error
}

func (bom BillOwnerManager) DeleteAsync(id int) error {
	tx := bom.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bom.GetById(id)
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

func (bom BillOwnerManager) EditAsync(id int, obj models.BillOwner) error {
	tx := bom.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	oldObj, err := bom.GetById(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&oldObj).Updates(map[string]interface{}{
		"user_id":      obj.UserId,
		"name":         obj.Name,
		"contact":      obj.Contact,
		"bank_account": obj.BankAccount,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
