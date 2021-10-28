package handlers

import (
	"fmt"
	specs "github.com/asavt7/nixchat_backend/api"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

// APIHandler handler
type APIHandler struct {
	mainRouter *echo.Echo
	service    *services.Services
	validator  *validator.Validate
}

// NewAPIHandler constructs APIHandler
func NewAPIHandler(service *services.Services) *APIHandler {
	return &APIHandler{
		service:    service,
		validator:  validator.New(),
		mainRouter: echo.New(),
	}
}

func (h *APIHandler) InitRoutes(cfg *config.Config) http.Handler {
	h.mainRouter.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	h.mainRouter.Use(middleware.Recover())

	/*
		mainRouter.GET("/oauth/google/login", h.handleGoogleLogin)
		mainRouter.GET("/oauth/google/callback", h.handleGoogleCallback)
	*/
	h.mainRouter.GET("/swagger/*", echoSwagger.WrapHandler)
	specs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	/*
		mainRouter.POST("/sign-in", h.signIn)
	*/
	h.mainRouter.POST("/sign-up", h.signUp)

	h.mainRouter.GET("/health/readiness", h.readinessProbe)
	h.mainRouter.GET("/health/liveness", h.livenessProbe)

	/*
		API := mainRouter.Group(apiPath)

		API.Use(parseAccessToken(), h.tokenRefresherMiddleware)

		usersAPI := API.Group("/users/:userId")

		usersAPI.GET("/posts", h.getUserPosts)
		usersAPI.POST("/posts", h.createPost)

		usersAPI.GET("/posts/:postId", h.getUserPostByID)
		usersAPI.DELETE("/posts/:postId", h.deletePost)
		usersAPI.PUT("/posts/:postId", h.updatePost)

		usersAPI.GET("/posts/:postId/comments", h.getCommentsByPostID)
		usersAPI.POST("/posts/:postId/comments", h.createComment)

		usersAPI.DELETE("/posts/:postId/comments/:commentId", h.deleteComment)
		usersAPI.PUT("/posts/:postId/comments/:commentId", h.updateComment)
	*/

	return h.mainRouter
}
