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


type BillMemberController struct {
	DB *gorm.DB
}

func NewBillMemberController(DB *gorm.DB) BillMemberController {
	return BillMemberController{DB}
}

// GetAllBillMembers godoc
//	@Summary		Get all Bill Members
//	@Description	Get all Bill Members
//	@Tags			billmembers
//	@Produce		json
//	@Success		200	{array}		models.BillMember
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billmembers [get]
func (bc BillMemberController) GetAll(c *gin.Context) {
	var billMembers []models.BillMember
	result := bc.DB.Where("deleted_at IS NULL").Find(&billMembers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, billMembers)
}

// GetBillMemberByID godoc
//	@Summary		Get a bill member by ID
//	@Description	Get a bill member by ID
//	@Tags			billmembers
//	@Produce		json
//  @Param			id	path int true "data"
//	@Success		200	{array}		models.BillMember
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billmembers/{id} [get]
func (bc BillMemberController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var billMember models.BillMember
	result := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&billMember)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Error: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, billMember)
}

// CreateBillMember godoc
// @Summary Create a new bill member
// @Description Create a new bill member
// @Tags billmembers
// @Accept json
// @Produce json
// @Param data body models.BillMember true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billmembers [post]
func (bc BillMemberController) CreateAsync(c *gin.Context) {
	var billMember models.BillMember
	err := c.BindJSON(&billMember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	result := bc.DB.Create(&billMember)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, billMember)
}

// DeleteBillMember godoc
// @Summary Delete a bill member by ID
// @Description Delete a bill member by ID
// @Tags billmembers
// @Accept json
// @Produce json
// @Param id path int true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billmembers/{id} [delete]
func (bc BillMemberController) DeleteAsync(c *gin.Context) {
	id := c.Param("id")
	var billMember models.BillMember
	now := time.Now()
	result := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&billMember)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	result = bc.DB.Model(&billMember).Update("deleted_at", &now)
	result = bc.DB.Model(&billMember).Update("updated_at", &now)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, billMember)
}

//EditBillMember godoc
// @Summary Edit Bill Member
// @Description Edit Bill Member
// @Tags billmembers
// @Accept json
// @Produce json
// @Param id path int true "data"
// @Param data body models.BillMember true "data"
// @Success		200	{object}	models.BillOwner
// @Failure		404	{object}	helpers.ErrResponse "Page not found"
// @Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billmembers/{id} [put]
func (bc BillMemberController) EditAsync(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	now := time.Now()

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var findBillMember models.BillMember
	if err := bc.DB.Where("id = ? AND deleted_at IS NULL", id).First(&findBillMember).Error; err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var updateMember models.BillMember
	if err := c.ShouldBindJSON(&updateMember); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	tx := bc.DB.Begin()
	if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Failed to start transaction: " + tx.Error.Error()})
        return
    }

	updateData := map[string]interface{}{
		"BillId": updateMember.BillId,
		"Name": updateMember.Name,
		"PriceOwe": updateMember.PriceOwe,
		"UpdatedAt": now,
	}

	if err := tx.Model(&findBillMember).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil { 
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Transaction commit failed: " + err.Error()})
        return
    }

	c.JSON(http.StatusOK, findBillMember)
}