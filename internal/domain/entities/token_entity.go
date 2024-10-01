package entities

import (
	"fmt"
	"strings"
)

type TokenEntity struct {
	Token          string `dynamodbav:"tokenId" json:"token"`
	CardNumber     string `dynamodbav:"cardNumber" json:"cardNumber"`
	CardHolderName string `dynamodbav:"cardHolderName" json:"cardHolderName"`
	ExpirationDate string `dynamodbav:"expirationDate" json:"expirationDate"`
	CVV            string `dynamodbav:"cvv" json:"cvv"`
}

func (t *TokenEntity) GetExpirationYear() string {
	expirationDateSplit := strings.Split(t.ExpirationDate, "/")
	fmt.Println(expirationDateSplit)
	expiryYear := expirationDateSplit[1]
	return expiryYear
}
func (t *TokenEntity) GetExpirationMonth() string {
	expirationDateSplit := strings.Split(t.ExpirationDate, "/")
	expirationMonth := expirationDateSplit[0]
	return expirationMonth
}
