package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type responseError struct {
	Status  int    `json:"-"`
	Message string `json:"error_message"`
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func dbErrorHandler(err error) responseError {
	var finalErr = responseError{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	}

	if err == sql.ErrNoRows || strings.Contains(
		sql.ErrNoRows.Error(), err.Error()) {
		finalErr = responseError{
			Status:  http.StatusNotFound,
			Message: "not found",
		}
	} else if pqErr, ok := err.(*pgconn.PgError); ok {
		switch pqErr.Code {
		case "23505": // duplicate
			finalErr = responseError{
				Status:  http.StatusForbidden,
				Message: pqErr.Detail,
			}
		case "23503": // fk constraint
			finalErr = responseError{
				Status:  http.StatusForbidden,
				Message: pqErr.Detail,
			}
		}
	}
	return finalErr
}
