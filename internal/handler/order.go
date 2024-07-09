package handler

import (
	"net/http"
	"orderTracking/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InputOrderItem struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	Username      string `json:"username"`
	From_location string `json:"from_location"`
	To_location   string `json:"to_location"`
	Status        string `json:"status"`
}

// type Answer struct {
// 	Message string `json:"message"`
// }

// createOrder godoc
//
//	@Summary		Create Order
//	@Security		ApiKeyAuth
//	@Tags			Order
//	@Description	Create a new order
//	@ID				create-order
//	@Accept			json
//	@Produce		json
//	@Param			body	body		InputOrderItem	true	"Order info"
//	@Success		200		{object}	models.OrderItem
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/lists/add_order [post]
func (h *Handler) createOrder(c *gin.Context) {
	var order InputOrderItem
	if err := c.BindJSON(&order); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user_id, _ := c.Get("userId")
	new_order := models.OrderItem{
		Track_id:      uuid.New().String(),
		Title:         order.Title,
		Description:   order.Description,
		Price:         order.Price,
		Date:          time.Now(),
		Id:            user_id.(int),
		From_location: order.From_location,
		To_location:   order.To_location,
	}
	msg, err := h.services.OrderItem.CreateOrderItem(new_order)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, msg)
}

// getStatus godoc
//
//	@Summary		Get Order Status
//	@Security		ApiKeyAuth
//	@Tags			Order
//	@Description	Get the status of an order by track ID
//	@ID				get-order-status
//	@Accept			json
//	@Produce		json
//	@Param			track_id	query	string	true	"Track ID"
//	@Success		200		{object}	models.OrderItem
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/lists/status [get]
func (h *Handler) getStatus(c *gin.Context) {
	// track_id, _ := c.Get("userId")
	id := c.Query("track_id")
	msg, err := h.services.OrderItem.GetOrderDetailsByTrackID(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, msg)
}
