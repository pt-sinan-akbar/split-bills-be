package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/models"
	"net/http"
	"strconv"
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
	bill, err := bc.BM.CreateAsync(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bill)
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

func (bc BillController) DynamicUpdateData(c *gin.Context) {
	var req struct {
		Tax     float64 `json:"tax"`
		Service float64 `json:"service"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	billId := c.Param("id")
	// deny request with negative values
	if req.Tax < 0 || req.Service < 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "no negative values, ok?"})
		return
	}
	// update tax/service, btw these can be zero so no validation
	updatedBill, err := bc.BM.DynamicUpdateData(billId, req.Tax, req.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedBill)
	return
}

func (bc BillController) DynamicUpdateItem(c *gin.Context) {
	var req struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	billId := c.Param("id")
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "wtf is this item_id"})
		return
	}
	// deny request with negative values
	if req.Price < 0 || req.Quantity < 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "no negative values, ok?"})
		return
	}
	// update item
	if itemId != 0 && req.Quantity != 0 {
		updatedBill, err := bc.BM.DynamicUpdateItem(billId, itemId, req.Name, req.Price, req.Quantity)
		if err != nil {
			if err.Error() == "failed to update item: no changes detected" {
				c.JSON(http.StatusAccepted, gin.H{"message": "no changes detected"})
				return
			}
			c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedBill)
		return
	}
	// if itemId or quantity is zero, return error
	c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "where is the item_id or quantity la u stupid or what?"})
	return
}

func (bc BillController) DynamicCreateItem(c *gin.Context) {
	var req struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	billId := c.Param("id")
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "where is the name la u stupid or what?"})
		return
	}
	// update item
	if req.Price > 0 && req.Quantity > 0 {
		updatedBill, err := bc.BM.DynamicCreateItem(billId, req.Name, req.Price, req.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, updatedBill)
		return
	}
	// if itemId or quantity is zero, return error
	c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "bad price or quantity, pls recheck your request"})
	return
}

func (bc BillController) DynamicDeleteItem(c *gin.Context) {
	billId := c.Param("id")
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil || itemId <= 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "wtf is this item_id"})
		return
	}
	updatedBill, err := bc.BM.DynamicDeleteItem(billId, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedBill)
}

func (bc BillController) DynamicDeleteMember(c *gin.Context) {
	billId := c.Param("id")
	memberId, err := strconv.Atoi(c.Param("member_id"))
	if err != nil || memberId <= 0 {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "wtf is this member_id"})
		return
	}
	if err := bc.BM.DynamicDeleteMember(billId, memberId); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, helpers.ErrResponse{Message: "Successfully deleted member from bill"})
}

func (bc BillController) UpsertOwner(c *gin.Context) {
	billId := c.Param("id")
	var req models.BillOwner
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: err.Error()})
		return
	}
	if req.Name == "" || req.Contact == "" || req.BankAccount == "" {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Message: "one or more fields are empty"})
		return
	}
	owner, err := bc.BM.UpsertOwner(billId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, owner)
}
