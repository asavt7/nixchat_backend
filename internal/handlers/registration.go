package handlers

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type signUpUserInput struct {
	Username string `json:"username" binding:"required" validate:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// signUp godoc
// @Tags auth
// @Summary signUp
// @Description register new user
// @ID signUp
// @Accept  json
// @Produce  json
// @Param signUpUserInput body signUpUserInput true "a body"
// @Success 200 {object} model.User
// @Failure 400 {object} Message
// @Failure 500 {object} Message
// @Router /sign-up [post]
func (h *APIHandler) signUp(c echo.Context) error {
	u := new(signUpUserInput)
	if err := c.Bind(u); err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	err := h.validator.Struct(u)
	if err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	newUser := model.User{
		Username: u.Username,
		Email:    u.Email,
	}

	createdUser, err := h.service.UserService.CreateUser(newUser, u.Password)
	if err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	return response(http.StatusCreated, createdUser, c)
}
