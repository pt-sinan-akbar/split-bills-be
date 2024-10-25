package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
)

type BillOwnerRouterController struct {
	billownerController controllers.BillOwnerController
}

func NewBillOwnerRouterController(bc controllers.BillOwnerController) BillOwnerRouterController {
	return BillOwnerRouterController{bc}
}

func (brc *BillOwnerRouterController) BillOwnerRouter(rg *gin.RouterGroup){
	rg.GET("/billowners", brc.billownerController.GetAll)
	rg.GET("/billowners/:id", brc.billownerController.GetByID)
	rg.POST("/billowners", brc.billownerController.CreateAsync)
	rg.PUT("/billowners/:id", brc.billownerController.EditAsync)
	rg.DELETE("/billowners/:id", brc.billownerController.DeleteAsync)
}