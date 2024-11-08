package manager

import (
	"fmt"
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

func (bdm BillDataManager) CreateAsync(billData models.BillData) error {
	result := bdm.DB.Create(&billData)
	return result.Error
}

func (bdm BillDataManager) DeleteAsync(id int) error {
	tx := bdm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bdm.GetByID(id)
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

func (bdm BillDataManager) EditAsync(id int, obj *models.BillData) error {
	tx := bdm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	oldObj, err := bdm.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&oldObj).Updates(map[string]interface{}{
		"bill_id":    obj.BillId,
		"store_name": obj.StoreName,
		"sub_total":  obj.SubTotal,
		"discount":   obj.Discount,
		"tax":        obj.Tax,
		"service":    obj.Service,
		"total":      obj.Total,
		"misc":       obj.Misc,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (bdm BillDataManager) DynamicUpdateRecalculateData(billData *models.BillData, subtotal []float64) error {
	var newSubtotal float64
	for _, v := range subtotal {
		newSubtotal += v
	}
	oldSubtotal := billData.SubTotal
	oldTax := billData.Tax
	oldService := billData.Service
	oldTaxPercent := oldTax / oldSubtotal
	oldServicePercent := oldService / oldSubtotal
	newTax := newSubtotal * oldTaxPercent
	newService := newSubtotal * oldServicePercent
	newTotal := newSubtotal + newTax + newService
	billData.SubTotal = newSubtotal
	billData.Tax = newTax
	billData.Service = newService
	billData.Total = newTotal
	if err := bdm.EditAsync(int(billData.ID), billData); err != nil {
		return err
	}
	fmt.Println(newSubtotal, newTax, newService, newTotal)
	return nil
}
