package ports

type BankGateway interface {
	ProcessPayment(cardNumber, expiryMonth, expiryYear, currency string, amount float64) (string, error)
	RefundPayment(paymentNumber, currency string, amount float64) (string, error)
}
