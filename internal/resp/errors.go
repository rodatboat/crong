package resp

import "errors"

var (
	ErrBadRequest          = errors.New("Bad Request")
	ErrNotFound            = errors.New("Not Found")
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrValidation          = errors.New("Validation Error")

	ErrInvalidCron = errors.New("Invalid Cron Expression")
)
