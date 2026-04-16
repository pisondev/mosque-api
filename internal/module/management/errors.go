package management

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrValidation      = errors.New("validation failed")
	ErrTagInUse        = errors.New("tag is still in use")
	ErrUnauthorizedCtx = errors.New("unauthorized tenant context")
	ErrPaymentRequired = errors.New("payment required before setup")
)
