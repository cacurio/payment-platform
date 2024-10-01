package dtos

type ChargeDTO struct {
	Token    string  `json:"token"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
