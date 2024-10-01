package domain

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidCharge    = errors.New("invalid charge")
	ErrInvalidRefund    = errors.New("invalid refund")
	ErrNoFoundCharge    = errors.New("no found charge")
	ErrNoCreatedCharge  = errors.New("no created charge")
	ErrNoProcessPayment = errors.New("no process payment")
)
