package handlers

import (
	"net/http"

	errorhandler "github.com/ak-karimzai/bank-api/internel/error_handler"
	"github.com/gin-gonic/gin"
)

var errMap map[errorhandler.Error]int = map[errorhandler.Error]int{
	errorhandler.InternealServer: http.StatusInternalServerError,
	errorhandler.Forbidden:       http.StatusForbidden,
	errorhandler.NotFound:        http.StatusNotFound,
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func toHttpError(err errorhandler.ResponseError) int {
	httpErr, ok := errMap[err.Status]
	if !ok {
		return http.StatusBadRequest
	}
	return httpErr
}
