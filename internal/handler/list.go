package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	id, _ := c.Get(userCTX)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
