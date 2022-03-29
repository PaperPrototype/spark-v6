package payments

import (
	"main/helpers"

	"github.com/stripe/stripe-go/v72"
)

// we take a 15 percent cut
const PercentageShare float32 = 0.15
const MaxCoursePrice uint16 = 10

func Setup() {
	// set global stripe API key
	stripe.Key = helpers.GetStripeKey()
}
