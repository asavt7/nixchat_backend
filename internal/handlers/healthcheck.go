package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// readinessProbe
// @Summary readinessProbe
// @Description indicates that apps is ready to serve traffic
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/readiness [get]
func (h *APIHandler) readinessProbe(c echo.Context) error {
	// todo https://github.com/asavt7/nixchat_backend/issues/11
	return response(http.StatusOK, nil, c)
}

// livenessProbe
// @Summary livenessProbe
// @Description indicates that apps is alive or dead
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/liveness [get]
func (h *APIHandler) livenessProbe(c echo.Context) error {
	// todo https://github.com/asavt7/nixchat_backend/issues/11
	return response(http.StatusOK, nil, c)
}
