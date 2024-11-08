package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
	"strconv"
)

type BillMemberController struct {
	BMM *manager.BillMemberManager
}

func NewBillMemberController(BMM *manager.BillMemberManager) BillMemberController {
	return BillMemberController{BMM}
}

// GetAllBillMembers godoc
//
//	@Description	Get all Bill Members
//	@Summary		Get all Bill Members
//	@Tags			billmembers
//	@Produce		json
//	@Success		200	{array}		models.BillMember
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billmembers [get]
func (bc BillMemberController) GetAll(c *gin.Context) {
	result, err := bc.BMM.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBillMemberByID godoc
//
//		@Summary		Get a bill member by ID
//		@Description	Get a bill member by ID
//		@Tags			billmembers
//		@Produce		json
//	 @Param			id	path int true "data"
//		@Success		200	{array}		models.BillMember
//		@Failure		404	{object}	helpers.ErrResponse "Page not found"
//		@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//		@Router			/billmembers/{id} [get]
func (bc BillMemberController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	result, err := bc.BMM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
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
	var obj models.BillMember
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BMM.CreateAsync(obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BMM.DeleteAsync(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill Member deleted successfully"})
}

// EditBillMember godoc
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
func (bc BillMemberController) EditAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var obj models.BillMember
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BMM.EditAsync(id, obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}
