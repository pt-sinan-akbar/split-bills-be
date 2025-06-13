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
	internalGroup := rg.Group("/bills")
	{
		internalGroup.GET("/", brc.billController.GetAll)
		internalGroup.POST("/", brc.billController.CreateAsync)
		internalGroup.GET("/:id", brc.billController.GetByID)
		internalGroup.PUT("/:id", brc.billController.EditAsync)
		internalGroup.DELETE("/:id", brc.billController.DeleteAsync)
		internalGroup.POST("/upload", brc.billController.UploadImage)
		dynamicGroup := internalGroup.Group("/dynamic")
		{
			dynamicGroup.PUT("/:id/data", brc.billController.DynamicUpdateData)
			dynamicGroup.POST("/:id/item/", brc.billController.DynamicCreateItem)
			dynamicGroup.PUT("/:id/item/:item_id", brc.billController.DynamicUpdateItem)
			dynamicGroup.DELETE("/:id/item/:item_id", brc.billController.DynamicDeleteItem)
			dynamicGroup.DELETE("/:id/member/:member_id", brc.billController.DynamicDeleteMember)
		}
	}
}
