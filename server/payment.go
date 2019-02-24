package server

import (
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func paymentHandler(price int, email string) (bool, string) {

	var sh_key = os.Getenv("SECRET_KEY")
	stripe.Key = sh_key
	fmt.Println(sh_key)
	params := &stripe.ChargeParams{
		Amount:       stripe.Int64(int64(price)),
		Currency:     stripe.String(string(stripe.CurrencyJPY)),
		ReceiptEmail: stripe.String(email),
	}
	token := "tok_mastercard"
	params.SetSource(token)

	ch, err := charge.New(params)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", ch.ID)
	return ch.Paid, ch.FailureMessage
}
