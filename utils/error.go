package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// NewError is the uniform interface for HTTP based API request errors
func NewError(ctx *gin.Context, status int, err error) {
	HTTPReqErrTotal.With(prometheus.Labels{
		"method": ctx.Request.Method,
		"path":   ctx.Request.URL.Path,
		"status": strconv.Itoa(status),
	}).Inc()

	e := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, e)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
