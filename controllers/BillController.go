package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
)

type BillController struct {
	DB *gorm.DB
}

func NewBillController(DB *gorm.DB) BillController {
	return BillController{DB}
}

// GetAllBills godoc
//
//	@Summary		Get all bills
//	@Description	Get all Bills from table
//	@Tags			bills
//	@Produce		json
//	@Success		200	{array}		models.Bill
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills [get]
func (bc BillController) GetAll(c *gin.Context) {
	var bills []models.Bill
	// result := bc.DB.Preload("BillData").Preload("BillOwner").Find(&bills)
	result := bc.DB.Where("deleted_at IS NULL").Find(&bills)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bills)
}

// GetBillByID godoc
//
//	@Summary		Get a bill by ID
//	@Description	Get bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id	path		string	true	"data"
//	@Success		200	{object}	models.Bill
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [get]
func (bc BillController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill
	result := bc.DB.Preload("BillData").Preload("BillOwner").Where("id = ?", id).First(&bill)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bill)
}

// CreateBill godoc
//
//	@Summary		Create a new bill
//	@Description	Create a new bill
//	@Tags			bills
//	@Accept			json
//	@Produce		json
//	@Param			data	body	models.Bill	true	"data"
//	@Success		201	{object}	models.Bill
//	@Failure		400	{object}	helpers.ErrResponse
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills [post]
func (bc BillController) CreateAsync(c *gin.Context) {
	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	result := bc.DB.Create(&bill)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, bill)
}

// DeleteBill godoc
//
//	@Summary		Delete a bill by ID
//	@Description	Delete bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id	path		string	true	"data"
//	@Success		200	{object}	helpers.ErrResponse
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [delete]
func (bc BillController) DeleteAsync(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill
	result := bc.DB.Where("id = ?", id).First(&bill)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	now := time.Now()
	result = bc.DB.Model(&bill).Update("deleted_at", &now)
	result = bc.DB.Model(&bill).Update("updated_at", &now)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}
	
	// c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Bill deleted successfully"})
	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted record"}) 
}

// UpdateBill godoc
//
//	@Summary		Update a bill by ID
//	@Description	Update bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id		path		string		true	"data"
//	@Param			data	body		models.Bill	true	"data"
//	@Success		200		{object}	models.Bill
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [put]
func (bc BillController) EditAsync(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	var findbill models.Bill
	if err := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&findbill).Error; err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var updateBill models.Bill
	if err := c.ShouldBindJSON(&updateBill); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	tx := bc.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: tx.Error.Error()})
		return
	}

	updateData := map[string]interface{}{
		"Name": updateBill.Name,
		"RawImage": updateBill.RawImage,
		"ID": updateBill.ID,
		"BillOwnerId": updateBill.BillOwnerId,
		"UpdatedAt": now,
	}

	if err := tx.Model(&findbill).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, findbill)
}
