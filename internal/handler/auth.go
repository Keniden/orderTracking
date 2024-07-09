package handler //обработчики HTTP-запросов

import (
	"net/http"
	"orderTracking/internal/models"

	"github.com/gin-gonic/gin"
)

// signUp godoc
//
//	@Summary		Sign Up
//	@Tags			Account
//	@Description	Создание нового пользователя в системе.
//	@ID				sign-up
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.User	true	"Информация о пользователе"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signIn godoc
//
//	@Summary		Sign In
//	@Tags			Account
//	@Description	Авторизация пользователя.
//	@ID				sign-in
//	@Accept			json
//	@Produce		json
//	@Param			body	body		signInInput	true	"Информация для входа"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": token,
	})
}
