package manager

import (
	"fmt"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type BillManager struct {
	DB   *gorm.DB
	BIM  *BillItemManager
	BDM  *BillDataManager
	BMIM *BillMemberItemManager
	BMM  *BillMemberManager
	BOM  *BillOwnerManager
}

func NewBillManager(DB *gorm.DB, BIM *BillItemManager, BDM *BillDataManager, BMIM *BillMemberItemManager, BMM *BillMemberManager, BOM *BillOwnerManager) BillManager {
	return BillManager{DB, BIM, BDM, BMIM, BMM, BOM}
}

// CRUD

func (bm BillManager) GetAll() ([]models.Bill, error) {
	var obj []models.Bill
	result := bm.DB.Where("deleted_at IS NULL").Preload("BillData").Preload("BillOwner").Preload("BillItem", "deleted_at IS NULL").Preload("BillMember", "deleted_at IS NULL").Preload("BillMemberItem", "deleted_at IS NULL").Find(&obj)
	return obj, result.Error
}

func (bm BillManager) GetByID(id string) (models.Bill, error) {
	var obj models.Bill
	result := bm.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("BillData").Preload("BillOwner").Preload("BillItem", "deleted_at IS NULL").Preload("BillMember", "deleted_at IS NULL").Preload("BillMemberItem", "deleted_at IS NULL").First(&obj)
	return obj, result.Error
}

func (bm BillManager) CreateAsync(bill models.Bill) (models.Bill, error) {
	generatedId, err := helpers.GenerateID()
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to generate id: %v", err)
	}
	bill.ID = generatedId
	// auto create bill data too if not given
	if bill.BillData == nil {
		bill.BillData = &models.BillData{
			BillId:    generatedId,
			StoreName: "",
			SubTotal:  0.0,
			Discount:  0.0,
			Tax:       0.0,
			Service:   0.0,
			Total:     0.0,
		}
	}
	result := bm.DB.Create(&bill)
	return bill, result.Error
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

// TODO: these dynamic shits isn't atomic, one failed operation will cause corrupted data
// TODO: also updated data can be not re-fetched if using pointer, the engineer is skill issue

func (bm BillManager) DynamicUpdateItem(billId string, itemId int, name string, price float64, quantity int) (models.Bill, error) {
	bill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get bill: %v", err)
	}
	isOnlyNameChanged, err := bm.BIM.DynamicUpdateItem(itemId, name, price, quantity)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to update item: %v", err)
	}
	if !isOnlyNameChanged {
		itemsSubtotal := 0.0
		for _, item := range bill.BillItem {
			if item.ID == int64(itemId) {
				itemsSubtotal += price * float64(quantity)
			} else {
				itemsSubtotal += item.Subtotal
			}
		}
		err = bm.BDM.DynamicUpdateRecalculateData(bill.BillData, itemsSubtotal)
		if err != nil {
			return models.Bill{}, fmt.Errorf("failed to update data: %v", err)
		}
	}
	updatedBill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get updated bill: %v", err)
	}
	return updatedBill, nil
}

func (bm BillManager) DynamicUpdateData(id string, tax float64, service float64) (models.Bill, error) {
	bill, err := bm.GetByID(id)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get bill: %v", err)
	}
	taxPercent, servicePercent, err := bm.BDM.DynamicUpdateData(bill.BillData, tax, service)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to update data: %v", err)
	}
	for _, item := range bill.BillItem {
		err = bm.BIM.DynamicUpdateRecalculateItem(item, taxPercent, servicePercent)
		if err != nil {
			return models.Bill{}, fmt.Errorf("failed to update item: %v", err)
		}
	}
	updatedBill, err := bm.GetByID(id)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get updated bill: %v", err)
	}
	return updatedBill, nil
}

func (bm BillManager) DynamicCreateItem(billId string, name string, price float64, quantity int) (models.Bill, error) {
	bill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get bill: %v", err)
	}
	taxPercent, servicePercent := bm.BDM.GetBillRates(bill.BillData)
	item, err := bm.BIM.DynamicCreate(billId, name, price, quantity, taxPercent, servicePercent)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to create item: %v", err)
	}
	itemsSubtotal := 0.0
	for _, item := range bill.BillItem {
		itemsSubtotal += item.Subtotal
	}
	itemsSubtotal += item.Subtotal
	err = bm.BDM.DynamicUpdateRecalculateData(bill.BillData, itemsSubtotal)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to update data: %v", err)
	}
	updatedBill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get updated bill: %v", err)
	}
	return updatedBill, nil
}

func (bm BillManager) DynamicDeleteItem(billId string, itemId int) (models.Bill, error) {
	bill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get bill: %v", err)
	}
	err = bm.BIM.DeleteAsync(itemId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to delete item: %v", err)
	}
	err = bm.BMIM.DeleteByItemId(itemId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to delete member items: %v", err)
	}
	itemsSubtotal := 0.0
	for _, item := range bill.BillItem {
		if item.ID != int64(itemId) {
			itemsSubtotal += item.Subtotal
		}
	}
	err = bm.BDM.DynamicUpdateRecalculateData(bill.BillData, itemsSubtotal)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to update data: %v", err)
	}
	updatedBill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get updated bill: %v", err)
	}
	return updatedBill, nil
}

func (bm BillManager) DynamicDeleteMember(billId string, memberId int) error {
	// ini bisa diubah, tujuannya buat ngecek aja ada billnya atau ngga
	_, err := bm.GetByID(billId)
	if err != nil {
		return fmt.Errorf("failed to get bill: %v", err)
	}
	err = bm.BMM.DeleteAsync(memberId)
	if err != nil {
		return fmt.Errorf("failed to delete member: %v", err)
	}
	err = bm.BMIM.DeleteByMemberId(memberId)
	if err != nil {
		return fmt.Errorf("failed to delete member items: %v", err)
	}
	return nil
}

func (bm BillManager) UpsertOwner(billId string, req models.BillOwner) (models.BillOwner, error) {
	owner := models.BillOwner{}
	bill, err := bm.GetByID(billId)
	if err != nil {
		return models.BillOwner{}, fmt.Errorf("failed to get bill: %v", err)
	}
	if req.ID == 0 {
		req.ID = 0
		owner, err = bm.BOM.CreateAsync(req)
		if err != nil {
			return owner, fmt.Errorf("failed to create owner: %v", err)
		}
		bill.BillOwnerId = &owner.ID
		if err := bm.EditAsync(billId, &bill); err != nil {
			return owner, fmt.Errorf("failed to update bill with new owner: %v", err)
		}
	} else {
		owner, err = bm.BOM.EditAsync(int(req.ID), req)
		if err != nil {
			return owner, fmt.Errorf("failed to edit owner: %v", err)
		}
	}
	return owner, nil
}

func (bm BillManager) UpsertMemberItems(billId string, memberItems []models.BillMemberItem) (models.Bill, error) {
	err := bm.BMIM.UpsertMany(memberItems)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to upsert member items: %v", err)
	}
	bill, err := bm.GetByID(billId)
	if err != nil {
		return models.Bill{}, fmt.Errorf("failed to get bill: %v", err)
	}
	return bill, nil
}

func (bm BillManager) ValidateBill(billId string) (error, error) {
	bill, err := bm.GetByID(billId)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill: %v", err)
	}
	// can be extended to check other things like if bill is empty, or if items are empty, etc.
	err = bm.validateMemberItemQty(bill)
	if err != nil {
		return err, nil
	}
	err = bm.validateData(bill)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (bm BillManager) validateMemberItemQty(bill models.Bill) error {
	// below is vibe coded, but it's correct tho
	// Map to track allocated quantities per item
	allocatedQty := make(map[int64]int64)
	// Count members using each item (for null quantity cases)
	memberCount := make(map[int]int)
	// Track items with explicitly set quantities
	itemsWithExplicitQty := make(map[int]bool)

	// Count members and track explicit allocations
	for _, memberItem := range bill.BillMemberItem {
		if memberItem.DeletedAt == nil {
			if memberItem.Quantity != nil {
				allocatedQty[int64(memberItem.BillItemId)] += *memberItem.Quantity
				itemsWithExplicitQty[memberItem.BillItemId] = true
			} else {
				// Count this member for equal splitting
				memberCount[memberItem.BillItemId]++
			}
		}
	}

	// Check if allocated quantities match item quantities
	var errors []string
	for _, item := range bill.BillItem {
		itemID := int(item.ID)

		// Skip check if item has a mixture of null and non-null quantities
		if itemsWithExplicitQty[itemID] && memberCount[itemID] > 0 {
			errors = append(errors,
				fmt.Sprintf("%s has mixed allocation (both per person and equal split)", item.Name))
			continue
		}

		totalAllocated := allocatedQty[item.ID]

		// If using equal splitting (all null quantities)
		if !itemsWithExplicitQty[itemID] && memberCount[itemID] > 0 {
			// This is valid - members will split equally
			continue
		}

		// If using explicit quantities, ensure they match
		if totalAllocated != item.Qty {
			errors = append(errors,
				fmt.Sprintf("%s (allocated: %d, total: %d)",
					item.Name, totalAllocated, item.Qty))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("item allocation issues: %s", strings.Join(errors, ", "))
	}
	return nil
}

func (bm BillManager) validateData(bill models.Bill) error {
	errorPrefix := "data issues: "
	taxExist := bill.BillData.Tax != 0
	serviceExist := bill.BillData.Service != 0
	totalTaxItems, totalServiceItems := 0.0, 0.0
	for _, item := range bill.BillItem {
		if taxExist && item.Tax != 0 {
			totalTaxItems += item.Tax
		}
		if serviceExist && item.Service != 0 {
			totalServiceItems += item.Service
		}
	}
	if taxExist && totalTaxItems != bill.BillData.Tax {
		return fmt.Errorf("%s mismatch between total tax from items (%f) and bill data tax (%f)", errorPrefix, totalTaxItems, bill.BillData.Tax)
	}
	if serviceExist && totalServiceItems != bill.BillData.Service {
		return fmt.Errorf("%s mismatch between total service from items (%f) and bill data service (%f)", errorPrefix, totalServiceItems, bill.BillData.Service)
	}
	if bill.BillData.SubTotal != bill.BillData.Total-bill.BillData.Tax-bill.BillData.Service-bill.BillData.Discount {
		return fmt.Errorf("%s bill data sub total (%f) does not match calculated total (%f)",
			errorPrefix, bill.BillData.SubTotal, bill.BillData.Total-bill.BillData.Tax-bill.BillData.Service-bill.BillData.Discount)
	}
	return nil
}

func (bm BillManager) FinalizeBill(billId string) error {
	bill, err := bm.GetByID(billId)
	if err != nil {
		return fmt.Errorf("failed to get bill: %v", err)
	}
	// Create a map to track member's owed amounts
	memberTotal := make(map[int]float64)

	// Process each item to calculate member costs
	for _, item := range bill.BillItem {
		if item.DeletedAt != nil {
			continue // Skip deleted items
		}
		// Find all non-deleted member items for this item
		var memberItems []models.BillMemberItem
		for _, mi := range bill.BillMemberItem {
			if mi.BillItemId == int(item.ID) && mi.DeletedAt == nil {
				memberItems = append(memberItems, mi)
			}
		}

		if len(memberItems) == 0 {
			return fmt.Errorf("item %s having no member assigned, err: %v", item.Name, err)
		}

		// Calculate total price of this item including tax and service
		itemTotalPrice := item.Subtotal + item.Tax + item.Service - item.Discount

		// split per person
		memberQty := make(map[int]int64)
		// split equal
		memberEqualAssigned := make(map[int]bool)

		for _, mi := range memberItems {
			if mi.Quantity != nil {
				memberQty[mi.BillMemberId] += *mi.Quantity
			} else {
				memberEqualAssigned[mi.BillMemberId] = true
			}
		}

		// if both explicit and null quantities exist, reject
		if len(memberQty) > 0 && len(memberEqualAssigned) > 0 {
			return fmt.Errorf("item %s has mixed allocation (both explicit and equal split), err: %v", item.Name, err)
		}

		// Calculate price per quantity unit
		pricePerUnit := itemTotalPrice / float64(item.Qty)

		// Distribute owed amounts to each member
		// Split per person
		for memberId, qty := range memberQty {
			if qty <= 0 {
				return fmt.Errorf("item %s has member %d with non-positive quantity %d, err: %v", item.Name, memberId, qty, err)
			}
			// Calculate owed amount for this member
			owed := pricePerUnit * float64(qty)
			memberTotal[memberId] += owed
		}
		// Split equal
		if len(memberEqualAssigned) > 0 {
			equalOwed := itemTotalPrice / float64(len(memberEqualAssigned))
			for memberId := range memberEqualAssigned {
				memberTotal[memberId] += equalOwed
			}
		}
	}

	// Update each member's PriceOwe field
	for _, member := range bill.BillMember {
		if member.DeletedAt != nil {
			continue
		}

		memberID := int(member.ID)
		owed := memberTotal[memberID]
		member.PriceOwe = &owed
		if _, err := bm.BMM.EditAsync(memberID, member); err != nil {
			return fmt.Errorf("failed to update member %d: %v", memberID, err)
		}
	}

	return nil
}
