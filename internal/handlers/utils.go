package handlers

import (
	"errors"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

// Message - common message
type Message struct {
	Message string `json:"message"`
}

func responseMessage(status int, message string, c echo.Context) error {
	return response(status, Message{Message: message}, c)
}

func response(status int, body interface{}, c echo.Context) error {
	ctype := c.Request().Header.Get(echo.HeaderContentType)
	acceptType := c.Request().Header.Get(echo.HeaderAccept)
	if len(acceptType) == 0 {
		acceptType = ctype
	}
	switch {
	case strings.Contains(acceptType, echo.MIMEApplicationJSON):
		return c.JSON(status, body)
	default:
		return c.JSON(status, body)
	}
}

const pageSizeParamName = "size"
const pageOffsetParamName = "offset"
const defaultPageSize = 10
const defaultPageOffset = 0
const maxPageSize = 1000

func parsePaginationArgs(c echo.Context) (model.PagedQuery, error) {
	sizeStr := c.QueryParam(pageSizeParamName)
	size := defaultPageSize
	if sizeStr != "" {
		parsedSize, err := strconv.Atoi(sizeStr)
		if err != nil {
			return model.PagedQuery{}, err
		}
		size = parsedSize
	}

	offsetStr := c.QueryParam(pageOffsetParamName)
	offset := defaultPageOffset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return model.PagedQuery{}, err
		}
		offset = parsedOffset
	}

	return model.PagedQuery{
		Size:   size,
		Offset: offset,
	}, nil
}

func validatePaginationArgs(validator *validator.Validate, query model.PagedQuery) error {
	if query.Size > maxPageSize {
		return errors.New("exceeded the maximum allowed size of page")
	}
	return validator.Struct(query)

}
