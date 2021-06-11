package statuserror

import "errors"

const (
	StatusCodeBadParams  = "invalid params"
	StatusCodeServerErr  = "server err"
	StatusCodeNullUser   = "user err"
	StatusCodeInvalidJWT = "invalid-jwt"
	StatusNotFilled = "all sections must be filled"
)

type IStatusError interface {
	Error() string
	StatusCode() string
	HttpCode() int
}

type StatusError struct {
	httpCode   int
	statusCode string
	err        error
}

func New(httpCode int, statusCode string, err error) *StatusError {
	return &StatusError{
		httpCode:   httpCode,
		statusCode: statusCode,
		err:        err,
	}
}

func (statusError *StatusError) Error() string {
	return statusError.err.Error()
}

func (statusError *StatusError) StatusCode() string {
	return statusError.statusCode
}

func (statusError *StatusError) HttpCode() int {
	return statusError.httpCode
}

var NotAuthorized = New(401, StatusCodeInvalidJWT, errors.New("not authorized"))
var NullUser = New(500, StatusCodeNullUser, errors.New("user is null"))

