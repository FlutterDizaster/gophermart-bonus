package sharederrors

import "errors"

var (
	ErrOrderAlreadyLoaded      = errors.New("error order already loaded")
	ErrOrderLoadedByAnotherUsr = errors.New("error order loaded by another user")
)
