package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"time"
)

type BillItemManager struct {
	DB *gorm.DB
}

func NewBillItemManager(DB *gorm.DB) BillItemManager {
	return BillItemManager{DB}
}

func (bim BillItemManager) DynamicUpdateItem(itemId int, price float64, quantity int) error {
	item, err := bim.GetByID(itemId)
	if err != nil {
		return fmt.Errorf("failed to get item: %v", err)
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
	item.Qty = int64(quantity)
	item.Price = price
	item.Subtotal = newSubtotal
	item.Tax = newTax
	item.Service = newService
	fmt.Println(int64(quantity), price, newSubtotal, newTax, newService)

	err = bim.EditAsync(itemId, item)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}
	return nil
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

func (bim BillItemManager) EditAsync(id int, updateObj models.BillItem) error {
	tx := bim.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bim.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	obj.Name = updateObj.Name
	obj.Qty = updateObj.Qty
	obj.Price = updateObj.Price
	obj.Subtotal = updateObj.Subtotal
	obj.Tax = updateObj.Tax
	obj.Service = updateObj.Service
	obj.Discount = updateObj.Discount
	obj.UpdatedAt = time.Now()

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
