package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"time"
)

type BillManager struct {
	DB  *gorm.DB
	BIM *BillItemManager
	BDM *BillDataManager
}

func NewBillManager(DB *gorm.DB, BIM *BillItemManager, BDM *BillDataManager) BillManager {
	return BillManager{DB, BIM, BDM}
}

// CRUD

func (bm BillManager) GetAll() ([]models.Bill, error) {
	var obj []models.Bill
	result := bm.DB.Where("deleted_at IS NULL").Preload("BillData").Preload("BillOwner").Find(&obj)
	if result.Error != nil {
		return obj, fmt.Errorf("failed to get all bills: %v", result.Error)
	}
	return obj, nil
}

func (bm BillManager) GetByID(id string) (models.Bill, error) {
	var obj models.Bill
	result := bm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("BillData").Preload("BillOwner").First(&obj)
	if result.Error != nil {
		return obj, result.Error
	}
	return obj, nil
}

func (bm BillManager) CreateAsync(bill *models.Bill) error {
	result := bm.DB.Create(&bill)
	if result.Error != nil {
		return fmt.Errorf("failed to create bill: %v", result.Error)
	}
	return nil
}

func (bm BillManager) DeleteAsync(id string) error {
	obj, _ := bm.GetByID(id)
	obj.UpdatedAt = time.Now()
	obj.DeletedAt = &obj.UpdatedAt
	result := bm.DB.Save(&obj)
	if result.Error != nil {
		return fmt.Errorf("failed to delete bill: %v", result.Error)
	}
	return nil
}

func (bm BillManager) EditAsync(bill *models.Bill) error {
	obj, _ := bm.GetByID(bill.ID)

	tx := bm.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	updateData := map[string]interface{}{
		"Name":        bill.Name,
		"RawImage":    bill.RawImage,
		"ID":          bill.ID,
		"BillOwnerId": bill.BillOwnerId,
		"UpdatedAt":   time.Now(),
	}

	if err := tx.Model(&obj).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update bill: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Business Logic

func (bm BillManager) SaveImage(image *multipart.FileHeader, id string) error {
	imageBytes, err := image.Open()
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer imageBytes.Close()

	fileBytes, err := io.ReadAll(imageBytes)
	if err != nil {
		return fmt.Errorf("failed to read image: %v", err)
	}

	err2 := helpers.UploadFile(id, fileBytes)
	if err2 != nil {
		return err2
	}

	return nil
}

func (bm BillManager) UploadImageToPython(image *multipart.FileHeader, id string) error {
	//var bill models.Bill
	//bill.ID = id
	//bill.RawImage = id
	//
	//// send image to python
	//
	//// if response.Error != nil
	//// save bill to db
	//result := bm.DB.Create(&bill)
	//if result.Error != nil {
	//	return fmt.Errorf("failed to save bill: %v", result.Error)
	//}
	//// iterate over python response to create BillData and BillItem
	//
	return nil
}

func (bm BillManager) UploadBill(image *multipart.FileHeader) error {
	// generate id
	id, err := helpers.GenerateID()
	if err != nil {
		return err
	}
	fmt.Println("DONE GENERATING ID")

	// upload image to python
	if err := bm.UploadImageToPython(image, id); err != nil {
		return err
	}
	fmt.Println("DONE UPLOADING IMAGE TO PYTHON")

	// upload image to supabase
	if err := bm.SaveImage(image, id); err != nil {
		fmt.Println("ERROR UPLOADING IMAGE TO SUPABASE")
		fmt.Println("ID: ", id)
		return err
	}
	fmt.Println("DONE UPLOADING IMAGE TO SUPABASE")

	return nil
}

func (bm BillManager) DynamicUpdateItem(itemId int, price float64, quantity int64) error {
	//_, err := BillItemManager.DynamicUpdateItem(itemId, price, quantity)
	return nil
}
