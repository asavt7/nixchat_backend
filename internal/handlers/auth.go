package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const refreshTokenCookieName = "refresh-token"
const accessTokenCookieName = "access-token"
const currentUserID = "currentUserID"

type signInUserInput struct {
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type signInResponse struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}

// signIn godoc
// @Tags auth
// @Summary signIn
// @Description signIn and get access and refresh tokens
// @ID signIn
// @Accept  json
// @Produce  json
// @Param signInUserInput body signInUserInput true "body"
// @Success 200 {object} signInResponse
// @Failure 400 {object} Message
// @Failure 500 {object} Message
// @Router /sign-in [post]
func (h *APIHandler) signIn(c echo.Context) error {
	u := new(signInUserInput)
	if err := c.Bind(u); err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	if len(u.Password) == 0 || len(u.Username) == 0 {
		u.Username = c.FormValue("username")
		u.Password = c.FormValue("password")
	}

	err := h.validator.Struct(u)
	if err != nil {
		return responseMessage(http.StatusBadRequest, err.Error(), c)
	}

	user, err := h.service.AuthorizationService.CheckUserCredentials(u.Username, u.Password)
	if err != nil {
		return responseMessage(http.StatusUnauthorized, "Password or Username is incorrect", c)
	}

	accessToken, refreshToken, err := h.generateTokensAndSetCookies(user.ID, c)
	if err != nil {
		return responseMessage(http.StatusUnauthorized, "Token is incorrect", c)
	}

	return response(http.StatusOK, signInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, c)
}

func (h *APIHandler) generateTokensAndSetCookies(userID uuid.UUID, c echo.Context) (accessToken, refreshToken string, err error) {
	accessToken, refreshToken, accessExp, refreshExp, err := h.service.AuthorizationService.GenerateTokens(userID)
	if err != nil {
		return accessToken, refreshToken, err
	}

	h.setTokenCookie(accessTokenCookieName, accessToken, accessExp, c)
	h.setTokenCookie(refreshTokenCookieName, refreshToken, refreshExp, c)

	return accessToken, refreshToken, nil
}

func (h *APIHandler) setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func (h *APIHandler) setUserCookie(userID string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "userId"
	cookie.Value = userID
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}
