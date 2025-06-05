package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"log"
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

func (bdm BillDataManager) DynamicUpdateRecalculateData(billData *models.BillData, newSubtotal float64) error {
	oldSubtotal := billData.SubTotal
	oldTotal := billData.Total
	newTotal := newSubtotal
	oldTax := billData.Tax
	oldService := billData.Service

	if oldTax != 0.0 {
		oldTaxPercent := oldTax / oldSubtotal
		newTax := newSubtotal * oldTaxPercent
		billData.Tax = newTax
		newTotal += newTax
	}
	if oldService != 0.0 {
		oldServicePercent := oldService / oldSubtotal
		newService := newSubtotal * oldServicePercent
		billData.Service = newService
		newTotal += newService
	}
	if oldSubtotal != newSubtotal {
		billData.SubTotal = newSubtotal
	}
	if oldTotal != newTotal {
		billData.Total = newTotal
	} else {
		log.Println("wtf no changes")
	}
	if err := bdm.EditAsync(int(billData.ID), billData); err != nil {
		return err
	}
	return nil
}

func (bdm BillDataManager) DynamicUpdateData(data *models.BillData, newTax float64, newService float64) (float64, float64, error) {
	if data.Tax == newTax && data.Service == newService {
		return 0.0, 0.0, fmt.Errorf("no changes")
	}
	taxPercent, servicePercent := bdm.GetBillRates(data)
	if newTax != 0.0 {
		taxPercent = newTax / data.SubTotal
	} else {
		newTax = 0.0
		taxPercent = 0.0
	}
	if newService != 0.0 {
		servicePercent = newService / data.SubTotal
	} else {
		newService = 0.0
		servicePercent = 0.0
	}
	data.Tax = newTax
	data.Service = newService
	data.Total = data.SubTotal + newTax + newService
	if err := bdm.EditAsync(int(data.ID), data); err != nil {
		return taxPercent, servicePercent, err
	}
	return taxPercent, servicePercent, nil
}

func (bdm BillDataManager) GetBillRates(data *models.BillData) (taxPercent float64, servicePercent float64) {
	if data.Tax != 0.0 {
		taxPercent = data.Tax / data.SubTotal
	} else {
		taxPercent = 0.0
	}
	if data.Service != 0.0 {
		servicePercent = data.Service / data.SubTotal
	} else {
		servicePercent = 0.0
	}
	return taxPercent, servicePercent
}
