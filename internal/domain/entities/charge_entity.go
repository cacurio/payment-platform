package entities

type ChargeEntity struct {
	ID            string  `dynamodbav:"id" json:"id"`
	Token         string  `dynamodbav:"token" json:"token"`
	Amount        float64 `dynamodbav:"amount" json:"amount"`
	Currency      string  `dynamodbav:"currency" json:"currency"`
	Status        string  `dynamodbav:"status" json:"status"`
	Created       string  `dynamodbav:"created" json:"created"`
	PaymentNumber string  `dynamodbav:"paymentNumber" json:"paymentNumber" `
	RefundId      string  `dynamodbav:"refundId" json:"refundId"`
}
