package errorhandler

import (
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ResponseError struct {
	Status  Error  `json:"-"`
	Message string `json:"error_message"`
}

func DbErrorHandler(err error) ResponseError {
	var finalErr = ResponseError{
		Status:  InternealServer,
		Message: err.Error(),
	}

	if err == sql.ErrNoRows || strings.Contains(
		sql.ErrNoRows.Error(), err.Error()) {
		finalErr = ResponseError{
			Status:  NotFound,
			Message: "not found",
		}
	} else if pqErr, ok := err.(*pgconn.PgError); ok {
		switch pqErr.Code {
		case "23505": // duplicate
			finalErr = ResponseError{
				Status:  AlreadyExist,
				Message: pqErr.Detail,
			}
		case "23503": // fk constraint
			finalErr = ResponseError{
				Status:  Forbidden,
				Message: pqErr.Detail,
			}
		}
	}
	return finalErr
}
