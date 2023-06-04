package config

import "errors"

var (
	ErrVariableNotFound = errors.New("variable not found")
	ErrInvalidValue     = errors.New("invalid value")
	ErrUnknowType       = errors.New("unknow type")
	ErrInitFactory      = errors.New("init factory")
)
