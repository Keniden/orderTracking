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
