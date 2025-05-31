package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/dto"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
)

type BillController struct {
	BM *manager.BillManager
}

func NewBillController(billManager *manager.BillManager) BillController {
	return BillController{billManager}
}

// CRUD punya maul

// GetAllBills godoc
//
//		@Summary		Get all bills
//		@Description	Get all Bills from table
//		@Tags			bills
//		@Produce		json
//		@Success		200	{array}		models.Bill
//		@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	 	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//		@Router			/bills [get]
func (bc BillController) GetAll(c *gin.Context) {
	obj, err := bc.BM.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// GetBillByID godoc
//
//	@Summary		Get a bill by ID
//	@Description	Get bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id	path		string	true	"data"
//	@Success		200	{object}	models.Bill
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [get]
func (bc BillController) GetByID(c *gin.Context) {
	id := c.Param("id")
	obj, err := bc.BM.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, obj)
}

// CreateBill godoc
//
//	@Summary		Create a new bill
//	@Description	Create a new bill
//	@Tags			bills
//	@Accept			json
//	@Produce		json
//	@Param			data	body	models.Bill	true	"data"
//	@Success		201	{object}	models.Bill
//	@Failure		400	{object}	helpers.ErrResponse
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills [post]
func (bc BillController) CreateAsync(c *gin.Context) {
	var obj models.Bill
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BM.CreateAsync(obj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obj)
}

// DeleteBill godoc
//
//	@Summary		Delete a bill by ID
//	@Description	Delete bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id	path		string	true	"data"
//	@Success		200	{object}	helpers.ErrResponse
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [delete]
func (bc BillController) DeleteAsync(c *gin.Context) {
	id := c.Param("id")

	if _, err := bc.BM.GetByID(id); err != nil {
		c.JSON(http.StatusNotFound, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BM.DeleteAsync(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted record"})
}

// UpdateBill godoc
//
//	@Summary		Update a bill by ID
//	@Description	Update bill by ID
//	@Tags			bills
//	@Produce		json
//	@Param			id		path		string		true	"data"
//	@Param			data	body		models.Bill	true	"data"
//	@Success		200		{object}	models.Bill
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/{id} [put]
func (bc BillController) EditAsync(c *gin.Context) {
	id := c.Param("id")

	var updatedObj models.Bill
	if err := c.ShouldBindJSON(&updatedObj); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BM.EditAsync(id, &updatedObj); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedObj)
}

// Business Logic

// UploadImage godoc
//
//	@Summary		Upload image
//	@Description	Upload image
//	@Tags			bills
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			image	formData	file	true	"image"
//	@Success		200	{object}	helpers.ErrResponse
//	@Failure		400	{object}	helpers.ErrResponse
//	@Failure		404	{object}	helpers.ErrResponse "Page not found"
//	@Failure		500	{object}	helpers.ErrResponse "Internal Server Error: Server failed to process the request"
//	@Router			/bills/upload [post]
func (bc BillController) UploadImage(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}

	if err := bc.BM.UploadBill(image); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully uploaded image"})
}

func (bc BillController) DynamicUpdate(c *gin.Context) {
	var req dto.DynamicUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	if req.BillId == "" {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "plz give me bill id onii-chan :3"})
		return
	}
	// deny request with all values filled
	if req.Tax != 0 && req.Service != 0 && req.Item.Id != 0 && req.Item.Quantity != 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "wtf too many values"})
		return
	}
	// deny request with negative values
	if req.Tax < 0 || req.Service < 0 || req.Item.Price < 0 || req.Item.Quantity < 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "no negative values, ok?"})
		return
	}
	// update item
	if req.Item.Id != 0 && req.Item.Quantity != 0 {
		err := bc.BM.DynamicUpdateItem(req.BillId, req.Item.Id, req.Item.Price, req.Item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Item successfully updated"})
		return
	}
	// update tax/service, btw these can be zero so no validation
	err := bc.BM.DynamicUpdateData(req.BillId, req.Tax, req.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Tax/Service successfully updated"})
	return
}
