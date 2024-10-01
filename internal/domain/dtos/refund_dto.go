package dtos

type RefundDTO struct {
	PaymentId string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}
