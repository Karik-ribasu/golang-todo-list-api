package errors

import (
	"database/sql"
)

type HttpError struct {
	StatusCode int
	Message    interface{}
}

// func (s *HttpError) Error() string {
// 	return fmt.Sprint("%d%s", s.StatusCode, s.Message)
// }

func SQLErrorCheck(err error) (httpError *HttpError) {
	switch err {
	case sql.ErrNoRows:
		return &HttpError{StatusCode: 404, Message: `{"error": "Not Found"}`}
	default:
		return &HttpError{StatusCode: 503, Message: `{error: "Internal Server Error"}`}
	}
}
