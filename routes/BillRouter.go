package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
)

type BillRouterController struct {
	billController controllers.BillController
}

func NewBillRouterController(bc controllers.BillController) BillRouterController {
	return BillRouterController{bc}
}

func (brc *BillRouterController) BillRouter(rg *gin.RouterGroup) {
	rg.GET("/bills", brc.billController.GetAll)
	rg.GET("/bills/:id", brc.billController.GetByID)
	rg.POST("/bills", brc.billController.CreateBill)
	rg.PUT("/bills/:id", brc.billController.UpdateBill)
	rg.DELETE("/bills/:id", brc.billController.DeleteBill)
}
