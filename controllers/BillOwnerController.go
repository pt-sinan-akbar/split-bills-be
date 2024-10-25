package controllers

import (
	"net/http"
	"strconv"

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
//	@Failure		404	{object}	helpers.ErrResponse
//	@Failure		500	{object}	helpers.ErrResponse
//	@Router			/billowners [get]
func (bc BillOwnerController) GetAll(c *gin.Context) {
    var billOwners []models.BillOwner
    result := bc.DB.Find(&billOwners)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, billOwners)
}

// GetBillOwnerByID godoc
//	@Summary		Get a Bill Owner by ID
//	@Description	Get Bill Owner by ID
//	@Tags			billowners
//	@Produce		json
//	@Param			id	path int true "data"
//	@Success		200	{array}		models.BillOwner
//	@Failure		404	{object}	helpers.ErrResponse
//	@Failure		500	{object}	helpers.ErrResponse
//	@Router			/billowners/{id} [get]
func (bc BillOwnerController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var billOwner models.BillOwner
	result := bc.DB.Where("id = ?", id).First(&billOwner)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, billOwner)
}

// CreateBillOwner godoc
// @Summary Create new bill owner
// @Description Create new bill owner
// @Tags billowners
// @Accept json
// @Produce json
// @Param data body models.BillOwner true "Bill Owner data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse
// @Failure		500	{object}	helpers.ErrResponse
// @Router /billowners [post]
func (bc BillOwnerController) CreateAsync(c *gin.Context){
	var billowner models.BillOwner

	if err := c.ShouldBindJSON(&billowner); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	result := bc.DB.Create(&billowner)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
	}

	c.JSON(http.StatusOK, billowner)
}

// DeleteBillOwner godoc
// @Summary Delete bill owner by ID
// @Description Delete bill owner by ID
// @Tags billowners
// @Produce json
// @Param id path int true "Bill Owner ID"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse
// @Failure		500	{object}	helpers.ErrResponse
// @Router /billowners/{id} [delete]
func (bc BillOwnerController) DeleteAsync(c *gin.Context){
	id := c.Param("id")
	var billowner models.BillOwner
	result := bc.DB.Where("id = ?", id).First(&billowner)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	result = bc.DB.Delete(&billowner)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, billowner)
}

// Edit Bill Owner godoc
// @Summary Edit Bill Owner
// @Description Edit Bill Owner
// @Tags billowners
// @Accept json
// @Produce json
// @Param id path integer true "Bill Owner ID"
// @Param data body models.BillOwner true "Edit Bill Owner Data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse
// @Failure		500	{object}	helpers.ErrResponse
// @Router /billowners/{id} [put]
func (bc BillOwnerController) EditAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var findBillOwner models.BillOwner
	if err := bc.DB.First(&findBillOwner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: "Bill Owner not found: " + err.Error()})
		return
	}

	var updateOwner models.BillOwner
	if err := c.ShouldBindJSON(&updateOwner); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid request Data: " + err.Error()})
		return
	}

	tx := bc.DB.Begin() 
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Failed to start transaction: " + tx.Error.Error()})
        return
    }

    if err := tx.Model(&findBillOwner).Updates(&updateOwner).Error; err != nil {
        tx.Rollback() 
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Could not update bill owner: " + err.Error()})
        return
    }

    if err := tx.Commit().Error; err != nil { 
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Transaction commit failed: " + err.Error()})
        return
    }

	c.JSON(http.StatusOK, findBillOwner)
}