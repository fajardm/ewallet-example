package errorcode

import "github.com/pkg/errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested data is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your data already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
)
