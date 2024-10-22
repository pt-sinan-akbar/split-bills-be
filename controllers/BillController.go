package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"net/http"
)

type BillController struct {
	DB *gorm.DB
}

func NewBillController(DB *gorm.DB) BillController {
	return BillController{DB}
}

// GetAllBills godoc
// @Summary Get all bills
// @Description Get all Bills from table
// @Tags bills
// @Produce  json
// @Success 200 {array} models.Bill
// @Failure 500 {object} helpers.ErrResponse
// @Router /bills [get]
func (bc BillController) GetAll(c *gin.Context) {
	var bills []models.Bill
	result := bc.DB.Preload("BillData").Preload("BillOwner").Find(&bills)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bills)
}

// GetBillByID godoc
// @Summary Get a bill by ID
// @Description Get bill by ID
// @Tags bills
// @Produce  json
// @Param id path string true "Bill ID"
// @Success 200 {object} models.Bill
// @Failure 404 {object} helpers.ErrResponse
// @Failure 500 {object} helpers.ErrResponse
// @Router /bills/{id} [get]
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
// @Summary Get a bill by ID
// @Description Get bill by ID
// @Tags bills
// @Produce  json
// @Param id path string true "Bill ID"
// @Success 200 {object} models.Bill
// @Failure 404 {object} helpers.ErrResponse
// @Failure 500 {object} helpers.ErrResponse
// @Router /bills/{id} [get]
func (bc BillController) CreateBill(c *gin.Context) {
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
// @Summary Delete a bill by ID
// @Description Delete bill by ID
// @Tags bills
// @Produce  json
// @Param id path string true "Bill ID"
// @Success 200 {object} helpers.ErrResponse
// @Failure 500 {object} helpers.ErrResponse
// @Router /bills/{id} [delete]
func (bc BillController) DeleteBill(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)
	result := bc.DB.Where("id = ?", id).Delete(&models.Bill{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Bill deleted successfully"})
}

// UpdateBill godoc
// @Summary Update a bill by ID
// @Description Update bill by ID
// @Tags bills
// @Produce  json
// @Param id path string true "Bill ID"
// @Param bill body models.Bill true "Bill Data"
// @Success 200 {object} models.Bill
// @Failure 400 {object} helpers.ErrResponse
// @Failure 500 {object} helpers.ErrResponse
// @Router /bills/{id} [put]
func (bc BillController) UpdateBill(c *gin.Context) {
    id := c.Param("id")
    var bill models.Bill

    if err := bc.DB.Preload("BillData").Where("bill_id = ?", id).First(&bill).Error; err != nil {
        c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: "Bill not found or doesn't exist"})
        return
    }

    if err := c.ShouldBindJSON(&bill); err != nil {
        c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
        return
    }

    tx := bc.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Model(&bill).Association("BillData").Replace(bill.BillData); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
        return
    }

    if err := tx.Save(&bill).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
        return
    }

    tx.Commit()
    c.JSON(http.StatusOK, bill)
}

