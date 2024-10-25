package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
)

type BillOwnerController struct {
	DB *gorm.DB
}

func NewBillOwnerController(DB *gorm.DB) BillOwnerController {
	return BillOwnerController{DB}
}

// GetAllBillOwner godoc
//	@Summary		Get all Bill Owner
//	@Description	Get all Bill Owner
//	@Tags			billowners
//	@Produce		json
//	@Success		200	{array}		models.BillOwner
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billowners [get]
func (bc BillOwnerController) GetAll(c *gin.Context) {
    var obj []models.BillOwner
    result := bc.DB.Where("deleted_at IS NULL").Find(&obj)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, obj)
}

// GetBillOwnerByID godoc
//	@Summary		Get a Bill Owner by ID
//	@Description	Get Bill Owner by ID
//	@Tags			billowners
//	@Produce		json
//	@Param			id	path int true "data"
//	@Success		200	{array}		models.BillOwner
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billowners/{id} [get]
func (bc BillOwnerController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var obj models.BillOwner
	result := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&obj)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Error: " + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, obj)
}

// CreateBillOwner godoc
// @Summary Create new bill owner
// @Description Create new bill owner
// @Tags billowners
// @Accept json
// @Produce json
// @Param data body models.BillOwner true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billowners [post]
func (bc BillOwnerController) CreateAsync(c *gin.Context){
	var obj models.BillOwner

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	// now := time.Now()

	result := bc.DB.Create(&obj)
	// result = bc.DB.Model(&billowner).Update("created_at", &now)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
	}

	c.JSON(http.StatusOK, obj)
}

// DeleteBillOwner godoc
// @Summary Delete bill owner by ID
// @Description Delete bill owner by ID
// @Tags billowners
// @Produce json
// @Param id path int true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billowners/{id} [delete]
func (bc BillOwnerController) DeleteAsync(c *gin.Context){
	id := c.Param("id")
	var obj models.BillOwner
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

// Edit Bill Owner godoc
// @Summary Edit Bill Owner
// @Description Edit Bill Owner
// @Tags billowners
// @Accept json
// @Produce json
// @Param id path integer true "data"
// @Param data body models.BillOwner true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billowners/{id} [put]
func (bc BillOwnerController) EditAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	now := time.Now()
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var obj models.BillOwner
	if err := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&obj).Error; err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: "Error: " + err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid request Data: " + err.Error()})
		return
	}

	tx := bc.DB.Begin() 
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Failed to start transaction: " + tx.Error.Error()})
        return
    }

	updateData := map[string]interface{}{
        "UserId":      obj.UserId,
        "Name":        obj.Name,
        "Contact":     obj.Contact,
        "BankAccount": obj.BankAccount,
        "UpdatedAt":   now,
    }

    if err := tx.Model(&obj).Updates(updateData).Error; err != nil {
        tx.Rollback() 
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Could not update bill owner: " + err.Error()})
        return
    }

    if err := tx.Commit().Error; err != nil { 
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Transaction commit failed: " + err.Error()})
        return
    }

	c.JSON(http.StatusOK, obj)
}