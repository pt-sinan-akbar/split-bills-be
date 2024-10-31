package controllers

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/manager"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
)

type BillController struct {
	DB     *gorm.DB
	Config initializers.Config
}

func NewBillController(DB *gorm.DB, Config initializers.Config) BillController {
	return BillController{DB, Config}
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
	var obj []models.Bill
	result := bc.DB.Where("deleted_at IS NULL").Preload("BillData").Preload("BillOwner").Find(&obj)
	// result := bc.DB.Where("deleted_at IS NULL").Find(&bills)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
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
	var obj models.Bill
	result := bc.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("BillData").Preload("BillOwner").First(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
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
	var obj models.Bill
	var imageName string

	billManager := manager.NewBillManager(bc.DB, bc.Config)

	image, err := c.FormFile("image")

	if image != nil {
		imageName, err = billManager.SaveImage(image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
			return
		}
	}

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if obj.ID == "" {
		id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
			return
		}
		obj.ID = id
	}

	if image != nil {
		obj.RawImage = imageName
	}

	result := bc.DB.Create(&obj)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, obj)
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
	var obj models.Bill
	result := bc.DB.Where("id = ?", id).First(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	now := time.Now()
	result = bc.DB.Model(&obj).Update("deleted_at", &now)
	result = bc.DB.Model(&obj).Update("updated_at", &now)

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

	var obj models.Bill
	if err := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&obj).Error; err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	tx := bc.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: tx.Error.Error()})
		return
	}

	updateData := map[string]interface{}{
		"Name":        obj.Name,
		"RawImage":    obj.RawImage,
		"ID":          obj.ID,
		"BillOwnerId": obj.BillOwnerId,
		"UpdatedAt":   now,
	}

	if err := tx.Model(&obj).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}
