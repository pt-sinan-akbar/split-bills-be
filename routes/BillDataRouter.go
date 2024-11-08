package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
)

type BillDataRouterController struct {
	billdataController controllers.BillDataController
}

func NewBillDataRouterController(bc controllers.BillDataController) BillDataRouterController {
	return BillDataRouterController{bc}
}

func (brc *BillDataRouterController) BillDataRouter(rg *gin.RouterGroup) {
	rg.GET("/billdatas", brc.billdataController.GetAll)
	rg.GET("/billdatas/:id", brc.billdataController.GetByID)
	rg.POST("/billdatas", brc.billdataController.CreateAsync)
	rg.PUT("/billdatas/:id", brc.billdataController.EditAsync)
	rg.DELETE("/billdatas/:id", brc.billdataController.DeleteAsync)
}
