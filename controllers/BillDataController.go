package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
	"strconv"
)

type BillDataController struct {
	BDM *manager.BillDataManager
}

func NewBillDataController(dataManager *manager.BillDataManager) BillDataController {
	return BillDataController{dataManager}
}

// GetAllBillData godoc
// @Summary Get all Bill Data
// @Description Get all Bill Data
// @Tags billdatas
// @Produce json
// @Success 200 {array} models.BillData
// @Failure 404 {object} helpers.ErrResponse "Page not found"
// @Failure 500 {object} helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billdatas [get]
func (bc BillDataController) GetAll(c *gin.Context) {
	obj, err := bc.BDM.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// GetBillDataByID godoc
// @Summary Get a Bill Data by ID
// @Description Get Bill Data by ID
// @Tags billdatas
// @Produce json
// @Param id path int true "data"
// @Success 200 {array} models.BillData
// @Failure 404 {object} helpers.ErrResponse "Page not found"
// @Failure 500 {object} helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billdatas/{id} [get]
func (bc BillDataController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "Invalid ID"})
		return
	}
	obj, err := bc.BDM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, obj)
}

// CreateBillData godoc
// @Summary Create new Bill Data
// @Description Create new Bill Data
// @Tags billdatas
// @Accept json
// @Produce json
// @Param data body models.BillData true "data"
// @Success 200 {object} models.BillData
// @Failure 404 {object} helpers.ErrResponse "Page not found"
// @Failure 500 {object} helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billdatas [post]
func (bc BillDataController) CreateAsync(c *gin.Context) {
	var obj models.BillData

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	err := bc.BDM.CreateAsync(&obj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// DeleteBillData godoc
// @Summary Delete a Bill Data by ID
// @Description Delete Bill Data by ID
// @Tags billdatas
// @Produce json
// @Param id path int true "data"
// @Success 200 {object} helpers.ErrResponse
// @Failure 404 {object} helpers.ErrResponse "Page not found"
// @Failure 500 {object} helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billdatas/{id} [delete]
func (bc BillDataController) DeleteAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	_, err = bc.BDM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err = bc.BDM.DeleteAsync(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted record"})
}

// Edit Bill Data godoc
// @Summary Edit Bill Data
// @Description Edit Bill Data
// @Tags billdatas
// @Accept json
// @Produce json
// @Param id path int true "data"
// @Param data body models.BillData true "data"
// @Success 200 {object} models.BillData
// @Failure 404 {object} helpers.ErrResponse "Page not found"
// @Failure 500 {object} helpers.ErrResponse "Internal Server Error: Server failed to process the request"
// @Router /billdatas/{id} [put]
func (bc BillDataController) EditAsync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	var obj models.BillData
	obj, err = bc.BDM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BDM.EditAsync(id, &obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}
