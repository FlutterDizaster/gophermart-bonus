package sharederrors

import "errors"

var (
	ErrNotEnoughFunds     = errors.New("error not enougs funds")
	ErrWithdrawNotAllowed = errors.New("error withdraw not allowed")
	ErrWrongOrderID       = errors.New("error wrong order id")
)
