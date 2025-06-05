package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
	"strconv"
)

type BillItemController struct {
	BIM *manager.BillItemManager
}

func NewBillItemController(billItemManager *manager.BillItemManager) BillItemController {
	return BillItemController{billItemManager}
}

// GetAllBillItem godoc
//
//		@Summary		Get all Bill Item
//		@Description	Get all Bill Item
//		@Tags			billitems
//		@Produce		json
//	 @Success		200	{array}		models.BillItem
//		@Failure		404	{object}	helpers.ErrResponse "Page not found"
//		@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//		@Router			/billitems [get]
func (bc BillItemController) GetAll(c *gin.Context) {
	obj, err := bc.BIM.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// GetBillItemById godoc
//
//	@Summary		Get a Bill Item by ID
//	@Description	Get Bill Item by ID
//	@Tags			billitems
//	@Produce		json
//	@Param			id	path int true "data"
//	@Success		200	{object}		models.BillItem
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billitems/{id} [get]
func (bc BillItemController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid ID"})
		return
	}

	obj, err := bc.BIM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
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
func (bc BillItemController) CreateAsync(c *gin.Context) {
	var obj models.BillItem
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	billItem, err := bc.BIM.CreateAsync(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, billItem)
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
func (bc BillItemController) DeleteAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid ID"})
		return
	}

	if err := bc.BIM.DeleteAsync(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted record"})
}

// EditBillItem godoc
// @Summary Edit a Bill Item
// @Description Edit a Bill Item
// @Tags billitems
// @Accept json
// @Produce json
// @Param id path int true "data"
// @Param data body models.BillItem true "data"
// @Success	200	{object}	models.BillItem
// @Failure	404	{object}	helpers.ErrResponse "Page not found"
// @Failure	500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billitems/{id} [put]
func (bc BillItemController) EditAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid ID"})
		return
	}

	var updatedObj models.BillItem
	if err := c.ShouldBindJSON(&updatedObj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BIM.EditAsync(id, updatedObj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedObj)
}
