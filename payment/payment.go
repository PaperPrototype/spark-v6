package payment

import (
	"log"

	"github.com/stripe/stripe-go/accountlink"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
)

// Then import the package

var key string = "sk_test_51K6JrUHOarh3IyKxXzJgsICsUp5UtWPqzX1eNCpTdvFEbm0utEthnYiD0h9jyAdHgLxzzNQhSFBmBa1byIQ78ASR001TfvNRMK"

func Setup() {

}

func createConnectedAccount() (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Type: stripe.String(string(stripe.AccountTypeExpress)),
	}
	result, err := account.New(params)
	if err != nil {
		log.Println("payment ERROR creating connected account:", err)
		return nil, err
	}

	return result, nil
}

func CreateAccountLink() error {
	params := &stripe.AccountLinkParams{
		Account:    stripe.String("acct_1032D82eZvKYlo2C"),
		RefreshURL: stripe.String("https://example.com/reauth"),
		ReturnURL:  stripe.String("https://example.com/return"),
		Type:       stripe.String("account_onboarding"),
	}
	accountlink.New(params)

	return nil
}
