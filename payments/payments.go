package payments

import (
	"main/helpers"

	"github.com/stripe/stripe-go/v72"
)

// we take a 25 percent cut
const PercentageShare float32 = 0.25

// 1000 => $10 USD
const MaxCoursePrice uint64 = 1000

const DescStripeConnectionNotSetup string = "Stripe connection not setup. Gifting course as free."

func Setup() {
	// set global stripe API key
	stripe.Key = helpers.GetStripeKey()
}
