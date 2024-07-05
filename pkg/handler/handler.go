package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp())
		auth.POST("/sign-in", h.signIn())
	}
	api := router.Group("/api")
	{
		lists := api.Group("/lists")
		{
			list.POST("", h.createList)
			list.GET("/", h.getAllLists)
			list.GET("/:id", h.getListById)
			list.PUT("/:id", h.updateList)
			list.DELETE("/:id", h.deleteList)
			items := lists.Group("/items")
			{
				items.POST("/", h.createItem())
				items.GET("/", h.getAllItems())
				items.GET("/:item_id", h.getItemById())
				items.PUT("/:item_id", h.updateItem())
				items.DELETE("/:item_id", h.deleteItem())
			}
		}
	}
	return router
}
