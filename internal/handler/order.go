package handler

import (
	"net/http"
	"orderTracking/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InputOrderItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Username    string `json:"username"`
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
		Track_id:    uuid.New().String(),
		Title:       order.Title,
		Description: order.Description,
		Price:       order.Price,
		Username:    order.Username,
		Date:        time.Now(),
		Id:          user_id.(int),
	}
	msg, err := h.services.OrderItem.CreateOrderItem(new_order)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, msg)
}

func (h *Handler) getAllItems(c *gin.Context) {

}
func (h *Handler) getItemById(c *gin.Context) {

}
func (h *Handler) updateItem(c *gin.Context) {

}
func (h *Handler) deleteItem(c *gin.Context) {
}
