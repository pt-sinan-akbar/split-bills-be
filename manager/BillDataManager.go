package manager

import (
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
)

type BillDataManager struct {
	DB *gorm.DB
}

func NewBillDataManager(DB *gorm.DB) BillDataManager {
	return BillDataManager{DB}
}

func (bdm BillDataManager) GetByID(id int) (models.BillData, error) {
	var obj models.BillData
	result := bdm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").First(&obj)
	return obj, result.Error
}
