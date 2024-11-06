package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"strconv"
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
	if err != nil {
		return item, fmt.Errorf("failed to update data: %v", err)
	}
	return item, err
}

func (bim BillItemManager) GetByID(id int) (models.BillItem, error) {
	var obj models.BillItem
	result := bim.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").First(&obj)
	return obj, result.Error
}
