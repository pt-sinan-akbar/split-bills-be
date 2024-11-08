package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/dto"
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
	return obj, result.Error
}

func (bm BillManager) GetByID(id string) (models.Bill, error) {
	var obj models.Bill
	result := bm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("BillData").Preload("BillOwner").Preload("BillItem").First(&obj)
	return obj, result.Error
}

func (bm BillManager) CreateAsync(bill models.Bill) error {
	generatedId, err := helpers.GenerateID()
	if err != nil {
		return fmt.Errorf("failed to generate id: %v", err)
	}
	bill.ID = generatedId
	result := bm.DB.Create(&bill)
	return result.Error
}

func (bm BillManager) DeleteAsync(id string) error {
	tx := bm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bm.GetByID(id)
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

func (bm BillManager) EditAsync(id string, bill *models.Bill) error {
	tx := bm.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	obj, err := bm.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	updateData := map[string]interface{}{
		"bill_owner_id": bill.BillOwnerId,
		"name":          bill.Name,
		"raw_image":     bill.RawImage,
		"updated_at":    time.Now(),
	}

	if err := tx.Model(&obj).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return tx.Error
	}

	return tx.Commit().Error
}

// Business Logic

func (bm BillManager) SaveImage(image *multipart.FileHeader, id string) error {
	imageBytes, err := image.Open()
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer func(imageBytes multipart.File) {
		err := imageBytes.Close()
		if err != nil {
			fmt.Println("failed to close image: ", err)
		}
	}(imageBytes)

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

func (bm BillManager) DynamicUpdateItem(req dto.DynamicUpdateRequest) error {
	err := bm.BIM.DynamicUpdateItem(req.Item.Id, req.Item.Price, req.Item.Quantity)
	if err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}
	bill, err := bm.GetByID(req.BillId)
	if err != nil {
		return fmt.Errorf("failed to get bill: %v", err)
	}
	var itemsSubtotal []float64
	for _, item := range bill.BillItem {
		itemsSubtotal = append(itemsSubtotal, item.Subtotal)
	}
	err = bm.BDM.DynamicUpdateRecalculateData(bill.BillData, itemsSubtotal)
	if err != nil {
		return fmt.Errorf("failed to recalculate data: %v", err)
	}
	return nil
}
