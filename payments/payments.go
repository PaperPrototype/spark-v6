package payments

import (
	"main/helpers"

	"github.com/stripe/stripe-go/v72"
)

// we take a 25 percent cut
const PercentageShare float32 = 0.20

// 1000 => $10 USD
const MaxCoursePrice uint64 = 1000

// if something goes wrong it is our fault.
// The customer should still be able to get what they are needing and we suffer the loss.
const (
	DescStripeConnectionNotSetup      string = "Stripe connection not setup. Gifting course for free."
	DescStripeConnectionNotSetupError string = "An error occured. Is your stripe connection setup? Gifting course for free."
	DescStripeChargesNotEnabled       string = "Charges are not enabled for this stripe account. Could not accept payment. Gifting course for free."
	DescStripeChargesNotEnabledError  string = "An error occured. Are charges enabled for this stripe account? Gifting course for free."
)

func Setup() {
	// set global stripe API key
	stripe.Key = helpers.GetStripeKey()
}
