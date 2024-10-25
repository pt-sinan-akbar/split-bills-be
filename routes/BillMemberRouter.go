package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
)

type BillMemberRouterController struct {
	billmemberController controllers.BillMemberController
}

func NewBillMemberRouterController(bc controllers.BillMemberController) BillMemberRouterController {
	return BillMemberRouterController{bc}
}

func (brc *BillMemberRouterController) BillMemberRouter(rg *gin.RouterGroup){
rg.GET("/billmembers", brc.billmemberController.GetAll)
	rg.GET("/billmembers/:id", brc.billmemberController.GetByID)
	rg.POST("/billmembers", brc.billmemberController.CreateAsync)
	rg.PUT("/billmembers/:id", brc.billmemberController.EditAsync)
	rg.DELETE("/billmembers/:id", brc.billmemberController.DeleteAsync)
}