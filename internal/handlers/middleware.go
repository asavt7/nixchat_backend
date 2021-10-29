package handlers

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func (h *APIHandler) parseAccessToken() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &model.Claims{},
		SigningKey:              []byte(h.service.AuthorizationService.GetAccessSigningKey()),
		TokenLookup:             "header:Authorization,cookie:" + accessTokenCookieName,
		ErrorHandlerWithContext: jwtErrorChecker,
		SuccessHandler: func(c echo.Context) {
			tok := c.Get("user")
			accessToken := tok.(*jwt.Token)
			claims := accessToken.Claims.(*model.Claims)
			userID := claims.UserID
			c.Set(currentUserID, userID)
		},
	})
}

func (h *APIHandler) tokenAutoRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tok := c.Get("user")
		if tok == nil {
			return responseMessage(http.StatusUnauthorized, "invalid token", c)
		}
		accessToken := tok.(*jwt.Token)
		claims := accessToken.Claims.(*model.Claims)
		err := h.service.ValidateAccessToken(claims)
		if err != nil {
			return responseMessage(http.StatusUnauthorized, "invalid token", c)
		}
		if h.service.AuthorizationService.IsNeedToRefresh(claims) {
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				refreshClaims, err := h.service.ParseRefreshTokenToClaims(rc.Value)
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						return responseMessage(http.StatusUnauthorized, "invalid token signature", c)
					}
					return responseMessage(http.StatusUnauthorized, "invalid token", c)
				}
				if claims.UserID != refreshClaims.UserID {
					return responseMessage(http.StatusUnauthorized, "invalid token", c)
				}
				err = h.service.AuthorizationService.ValidateRefreshToken(refreshClaims)
				if err != nil {
					return responseMessage(http.StatusUnauthorized, "invalid token", c)
				}
				_, _, err = h.generateTokensAndSetCookies(claims.UserID, c)
				if err != nil {
					return responseMessage(http.StatusUnauthorized, "invalid token", c)
				}
			}
		}
		return next(c)
	}
}

func jwtErrorChecker(err error, c echo.Context) error {
	return responseMessage(http.StatusUnauthorized, err.Error(), c)
}
