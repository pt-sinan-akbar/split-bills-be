package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
)

type BillItemRouterController struct {
	billitemController controllers.BillItemController
}

func NewBillItemRouterController(bc controllers.BillItemController) BillItemRouterController {
	return BillItemRouterController{bc}
}

func (brc *BillItemRouterController) BillItemRouter(rg *gin.RouterGroup) {
	rg.GET("/billitems", brc.billitemController.GetAll)
	rg.GET("/billitems/:id", brc.billitemController.GetByID)
	rg.POST("/billitems", brc.billitemController.CreateAsync)
	// rg.PUT("/billitems/:id", brc.billitemController.EditAsync)
	rg.DELETE("/billitems/:id", brc.billitemController.DeleteAsync)
}