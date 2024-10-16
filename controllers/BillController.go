package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
	"net/http"
)

type BillController struct {
	DB *gorm.DB
}

type errResponse struct {
	Message string `json:"message"`
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
// @Failure 500 {object} errResponse
// @Router /bills [get]
func (bc BillController) GetAll(c *gin.Context) {
	var bills []models.Bill
	result := bc.DB.Preload("BillData").Preload("BillOwner").Find(&bills)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, errResponse{Message: result.Error.Error()})
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
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /bills/{id} [get]
func (bc BillController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill
	result := bc.DB.Preload("BillData").Preload("BillOwner").Where("id = ?", id).First(&bill)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, errResponse{Message: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bill)
}

// func (bc BillController) CreateBill(c *gin.Context) {
// 	var bill models.Bill
// 	if err := c.ShouldBindJSON(&bill); err != nil {
// 		c.JSON(http.StatusBadRequest, errResponse{Message: err.Error()})
// 		return
// 	}

// 	result := bc.DB.Create(&bill)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, errResponse{Message: result.Error.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, bill)
// }

// func (bc BillController) DeleteBill(c *gin.Context) {
// 	id := c.Param("id")
// 	result := bc.DB.Delete(&models.Bill{}, id)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, errResponse{Message: result.Error.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, errResponse{Message: "Bill deleted successfully"})
// }

// func (bc BillController) UpdateBill(c *gin.Context) {
// 	id := c.Param("id")
// 	var bill models.Bill
// 	if err := c.ShouldBindJSON(&bill); err != nil {
// 		c.JSON(http.StatusBadRequest, errResponse{Message: err.Error()})
// 		return
// 	}

// 	result := bc.DB.Model(&bill).Where("id = ?", id).Updates(&bill)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, errResponse{Message: result.Error.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, bill)
// }
