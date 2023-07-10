package plainoracle

import "errors"

var (
	ErrInvalidAValue                = errors.New("invalid A value")
	ErrBalancesMustMatchMultipliers = errors.New("balances must match multipliers")
	ErrZero                         = errors.New("zero")
	ErrDDoesNotConverge             = errors.New("d does not converge")
	ErrTokenFromEqualsTokenTo       = errors.New("can't compare token to itself")
	ErrTokenIndexesOutOfRange       = errors.New("token index out of range")
	ErrAmountOutNotConverge         = errors.New("approximation did not converge")
	ErrTokenNotFound                = errors.New("token not found")
	ErrWithdrawMoreThanAvailable    = errors.New("cannot withdraw more than available")
	ErrD1LowerThanD0                = errors.New("d1 <= d0")
	ErrDenominatorZero              = errors.New("denominator should not be 0")
)
