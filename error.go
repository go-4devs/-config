package config

import "errors"

var (
	ErrValueNotFound = errors.New("value not found")
	ErrInvalidValue  = errors.New("invalid value")
	ErrUnknowType    = errors.New("unknow type")
	ErrInitFactory   = errors.New("init factory")
	ErrStopWatch     = errors.New("stop watch")
)
