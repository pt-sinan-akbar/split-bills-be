package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
)

type BillItemController struct {
	DB *gorm.DB
}

func NewBillItemController(DB *gorm.DB) BillItemController {
	return BillItemController{DB}
}

// GetAllBillItem godoc
//	@Summary		Get all Bill Item
//	@Description	Get all Bill Item
//	@Tags			billitems
//	@Produce		json
//  @Success		200	{array}		models.BillItem
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billitems [get]
func (bc BillItemController) GetAll(c *gin.Context){
	var obj []models.BillItem
	result := bc.DB.Where("deleted_at IS NULL").Preload("Bill").Find(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// GetBillItemById godoc
//	@Summary		Get a Bill Item by ID
//	@Description	Get Bill Item by ID
//	@Tags			billitems
//	@Produce		json
//	@Param			id	path int true "data"
//	@Success		200	{array}		models.BillItem
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billitems/{id} [get]
func (bc BillItemController) GetByID(c *gin.Context){
	id := c.Param("id")

	var obj models.BillItem
	result := bc.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("Bill").First(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// CreateBillItem godoc
// @Summary Create new Bill Item
// @Description Create new Bill Item
// @Tags billitems
// @Accept json
// @Produce json
// @Param data body models.BillItem true "data"
// @Success		200	{object}	models.BillItem
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router			/billitems [post]
func (bc BillItemController) CreateAsync(c *gin.Context){
	var obj models.BillItem

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	result := bc.DB.Create(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// DeleteBillItem godoc
// @Summary Delete a Bill Item
// @Description Delete a Bill Item
// @Tags billitems
// @Produce json
// @Param id path int true "data"
// @Success	200	{object}	helpers.ErrResponse
// @Failure	404	{object}	helpers.ErrResponse "Page not found"
// @Failure	500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billitems/{id} [delete]
func (bc BillItemController) DeleteAsync(c *gin.Context){
	id := c.Param("id")
	var obj models.BillItem

	now := time.Now()
	result := bc.DB.Where("id = ?", id).First(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	result = bc.DB.Model(&obj).Update("deleted_at", &now)
	result = bc.DB.Model(&obj).Update("updated_at", &now)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted record"}) 
}

