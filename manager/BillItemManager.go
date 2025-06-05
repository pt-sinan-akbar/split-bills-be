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

func (bim BillItemManager) DynamicUpdateItem(itemId int, name string, price float64, quantity int) (bool, error) {
	item, err := bim.GetByID(itemId)
	if err != nil {
		return false, fmt.Errorf("failed to get item: %v", err)
	}
	if item.Price == price && item.Qty == int64(quantity) {
		if item.Name != name {
			// only name has changed, update it
			item.Name = name
			err = bim.EditAsync(itemId, item)
			if err != nil {
				return false, fmt.Errorf("failed to update item name: %v", err)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("no changes detected")
		}
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
	item.Name = name
	item.Qty = int64(quantity)
	item.Price = price
	item.Subtotal = newSubtotal
	item.Tax = newTax
	item.Service = newService

	err = bim.EditAsync(itemId, item)
	if err != nil {
		return false, fmt.Errorf("failed to update item: %v", err)
	}
	return false, nil
}

func (bim BillItemManager) GetByID(id int) (models.BillItem, error) {
	var obj models.BillItem
	result := bim.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").Preload("BillMember").First(&obj)
	return obj, result.Error
}

func (bim BillItemManager) GetAll() ([]models.BillItem, error) {
	var obj []models.BillItem
	result := bim.DB.Where("deleted_at IS NULL").Preload("Bill").Find(&obj)
	return obj, result.Error
}

func (bim BillItemManager) CreateAsync(billItem models.BillItem) (models.BillItem, error) {
	result := bim.DB.Create(&billItem)
	return billItem, result.Error
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

	if err := tx.Model(&obj).Updates(map[string]interface{}{
		"name":       updateObj.Name,
		"qty":        updateObj.Qty,
		"price":      updateObj.Price,
		"subtotal":   updateObj.Subtotal,
		"tax":        updateObj.Tax,
		"service":    updateObj.Service,
		"discount":   updateObj.Discount,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (bim BillItemManager) DynamicUpdateRecalculateItem(item models.BillItem, taxPercent float64, servicePercent float64) error {
	item.Tax = item.Subtotal * taxPercent
	item.Service = item.Subtotal * servicePercent

	err := bim.EditAsync(int(item.ID), item)
	if err != nil {
		return fmt.Errorf("id %v, error: %v", item.ID, err)
	}
	return nil
}

func (bim BillItemManager) DynamicCreate(id string, name string, price float64, quantity int, taxPercent float64, servicePercent float64) (models.BillItem, error) {
	tax := price * float64(quantity) * taxPercent
	service := price * float64(quantity) * servicePercent
	subtotal := price * float64(quantity)
	item := models.BillItem{
		BillId:   id,
		Name:     name,
		Qty:      int64(quantity),
		Price:    price,
		Subtotal: subtotal,
		Tax:      tax,
		Service:  service,
		Discount: 0.0,
	}
	result, err := bim.CreateAsync(item)
	if err != nil {
		return models.BillItem{}, fmt.Errorf("failed to create item: %v", err)
	}
	return result, nil
}
