package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BillItemManager struct {
	DB *gorm.DB
}

func NewBillItemManager(DB *gorm.DB) BillItemManager {
	return BillItemManager{DB}
}

func (bim BillItemManager) DynamicUpdateItem(itemId int, price float64, quantity int64) (models.BillItem, error) {
	item, err := bim.GetByID(itemId)
	if err != nil {
		return item, fmt.Errorf("failed to get item: %v", err)
	}
	newSubtotal := price * float64(quantity)
	var taxPercent, newTax = 0.0, 0.0
	if item.Tax != 0 {
		taxPercent = item.Subtotal / item.Tax
		newTax = newSubtotal / taxPercent
	}
	var servicePercent, newService = 0.0, 0.0
	if item.Service != 0 {
		servicePercent = item.Subtotal / item.Service
		newService = newSubtotal / servicePercent
	}
	item.Subtotal = newSubtotal
	item.Tax = newTax
	item.Service = newService
	fmt.Println(newSubtotal, newTax, newService)

	//result := bim.DB.Save(&item)
	//if result.Error != nil {
	//	return fmt.Errorf("failed to update item: %v", result.Error)
	//}
	_, err = strconv.Atoi(item.BillId)
	if err != nil {
		return item, fmt.Errorf("internal server error")
	}
	return item, err
}

func (bim BillItemManager) GetByID(id int) (models.BillItem, error) {
	var obj models.BillItem
	result := bim.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").First(&obj)
	return obj, result.Error
}

func (bim BillItemManager) GetAll() ([]models.BillItem, error) {
	var obj []models.BillItem
	result := bim.DB.Where("deleted_at IS NULL").Preload("Bill").Find(&obj)
	return obj, result.Error
}

func (bim BillItemManager) CreateAsync(billItem models.BillItem) error {
	result := bim.DB.Create(&billItem)
	return result.Error
}

func (bim BillItemManager) DeleteAsync(id int) error {
	tx := bim.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bim.GetByID(id)
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

func (bim BillItemManager) EditAsync(id int, billItem models.BillItem) error {
	tx := bim.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	oldObj, err := bim.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&oldObj).Updates(map[string]interface{}{
		"bill_id":    billItem.BillId,
		"name":       billItem.Name,
		"qty":        billItem.Qty,
		"price":      billItem.Price,
		"subtotal":   billItem.Subtotal,
		"tax":        billItem.Tax,
		"service":    billItem.Service,
		"discount":   billItem.Discount,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
