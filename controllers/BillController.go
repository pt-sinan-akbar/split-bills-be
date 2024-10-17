package controllers

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/models"
)

// GetAllBills godoc
// @Summary Get all bills
// @Description Get all Bills from table
// @Tags bills
// @Produce  json
// @Success 200 {array} models.Bill
// @Failure 500 {object} gin.H
// @Router /bills [get]

func GetAll(c *gin.Context) {
	var bills []models.Bill
	result := initializers.DB.Preload("BillData").Preload("BillOwner").Find(&bills)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /bills/{id} [get]
func GetByID(c *gin.Context) {
	id := c.Param("id")
	var bill models.Bill
	result := initializers.DB.Preload("BillData").Preload("BillOwner").First(&bill, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bill)	
}