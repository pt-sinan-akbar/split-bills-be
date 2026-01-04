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
		internalGroup.POST("/upload", brc.billController.UploadImage)
		internalGroup.POST("/extract-bill-data", brc.billController.ExtractBillData)
		existingBillGroup := internalGroup.Group("/:id")
		{
			existingBillGroup.GET("", brc.billController.GetByID)
			existingBillGroup.PUT("", brc.billController.EditAsync)
			existingBillGroup.DELETE("", brc.billController.DeleteAsync)
			dynamicGroup := existingBillGroup.Group("/dynamic")
			{
				dynamicGroup.PUT("/data", brc.billController.DynamicUpdateData)
				dynamicGroup.POST("/item", brc.billController.DynamicCreateItem)
				dynamicGroup.PUT("/item/:item_id", brc.billController.DynamicUpdateItem)
				dynamicGroup.DELETE("/item/:item_id", brc.billController.DynamicDeleteItem)
				dynamicGroup.DELETE("/member/:member_id", brc.billController.DynamicDeleteMember)
				dynamicGroup.POST("/owner", brc.billController.UpsertOwner)
				dynamicGroup.POST("/member-items", brc.billController.UpsertMemberItems)
			}
			existingBillGroup.GET("/validate", brc.billController.ValidateBill)
			existingBillGroup.POST("/finalize", brc.billController.FinalizeBill)
			existingBillGroup.GET("/summary", brc.billController.GetBillSummary)
		}
	}
}
