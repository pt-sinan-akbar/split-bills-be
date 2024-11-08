package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
	"strconv"
)

type BillOwnerController struct {
	BOM *manager.BillOwnerManager
}

func NewBillOwnerController(BOM *manager.BillOwnerManager) BillOwnerController {
	return BillOwnerController{BOM}
}

// GetAllBillOwner godoc
//
//	@Summary		Get all Bill Owner
//	@Description	Get all Bill Owner
//	@Tags			billowners
//	@Produce		json
//	@Success		200	{array}		models.BillOwner
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/billowners [get]
func (bc BillOwnerController) GetAll(c *gin.Context) {
	obj, err := bc.BOM.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: "Error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, obj)
}

// GetBillOwnerByID godoc
//
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	obj, err := bc.BOM.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
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
func (bc BillOwnerController) CreateAsync(c *gin.Context) {
	var obj models.BillOwner
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BOM.CreateAsync(obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
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
func (bc BillOwnerController) DeleteAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BOM.DeleteAsync(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
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
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var obj models.BillOwner
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BOM.EditAsync(id, obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}
