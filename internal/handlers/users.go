package handlers

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// getUserInfo godoc
// @Tags users
// @Summary getUserInfo
// @Description getUserInfo
// @ID getUserInfo
// @Accept  json
// @Produce  json
// @Param userId path int true "userId"
// @Success 200 {object} model.User
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /api/v1/users/{userId} [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (h *APIHandler) getUserInfo(context echo.Context) error {
	userIDStr := context.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return responseMessage(http.StatusBadRequest, "missing or incorrect userId param, expected int", context)
	}

	user, err := h.service.UserService.GetByID(userID)
	if err != nil {
		return responseMessage(http.StatusInternalServerError, err.Error(), context)
	}
	return response(http.StatusOK, user, context)
}

// getUsers godoc
// @Tags users
// @Summary getUsers
// @Description getUsers
// @ID getUsers
// @Accept  json
// @Produce  json
// @Param size query int false "size" minimum(0) maximum(1000)
// @Param offset query int false "offset" minimum(0)
// @Success 200 {object} []model.User
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /api/v1/users [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (h *APIHandler) getUsers(c echo.Context) error {

	pagedQuery, err := parsePaginationArgs(c)
	if err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}
	if err := validatePaginationArgs(h.validator, pagedQuery); err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	users, err := h.service.UserService.GetAll(pagedQuery)
	if err != nil {
		return responseMessage(http.StatusInternalServerError, err.Error(), c)
	}
	return response(http.StatusOK, users, c)
}

// updateUser godoc
// @Tags users
// @Summary updateUser
// @Description updateUser
// @ID updateUser
// @Accept  json
// @Produce  json
// @Param userId path string true "userId"
// @Param post body model.UpdateUserInfo true "update input"
// @Success 200 {object} model.User
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /api/v1/users/{userId}/ [put]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (h *APIHandler) updateUser(context echo.Context) error {
	userIDStr := context.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return responseMessage(http.StatusBadRequest, "missing or incorrect userId param", context)
	}

	currentUser := context.Get(currentUserID).(uuid.UUID)
	if currentUser != userID {
		return responseMessage(http.StatusUnauthorized, "unauthorized : cannot change resource", context)
	}

	updateInput := new(model.UpdateUserInfo)
	if err := context.Bind(updateInput); err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), context)
	}
	if err := h.validator.Struct(updateInput); err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), context)
	}

	updatedUser, err := h.service.UserService.Update(userID, *updateInput)
	if err != nil {
		switch err.(type) {
		default:
			return responseMessage(http.StatusInternalServerError, err.Error(), context)
		}
	}
	return response(http.StatusOK, updatedUser, context)
}
