package dtos

type TokenDTO struct {
	Token          string `json:"token"`
	CardNumber     string `json:"cardNumber"`
	CardHolderName string `json:"cardHolderName"`
	ExpirationDate string `json:"expirationDate"`
	CVV            string `json:"cvv"`
}
