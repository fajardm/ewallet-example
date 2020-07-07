package errorcode

import (
	"github.com/pkg/errors"
	"net/http"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested data is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your data already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
	// ErrUnauthorized will throw if actor not authorized to access usecase
	ErrUnauthorized = errors.New("unauthorized")
)

var statusCode = map[error]int{
	ErrInternalServerError: http.StatusInternalServerError,
	ErrNotFound:            http.StatusNotFound,
	ErrConflict:            http.StatusConflict,
	ErrBadParamInput:       http.StatusBadRequest,
	ErrUnauthorized:        http.StatusUnauthorized,
}

func StatusCode(err error) int {
	if c, ok := statusCode[err]; ok {
		return c
	}
	return http.StatusInternalServerError
}
