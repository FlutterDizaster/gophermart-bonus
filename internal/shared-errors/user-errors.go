package sharederrors

import "errors"

var (
	ErrUserAlreadyExist     = errors.New("error user already exist")
	ErrWrongLoginOrPassword = errors.New("error wrong login or password")
)
