package router2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/db"
	"main/helpers"
	"main/payments"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/webhook"
)

func postStripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := helpers.GetStripeWebhook()
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		c.Writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		log.Println("payment_intent.succeeded")
	case "checkout.session.completed":
		log.Println("checkout.session.completed")

		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		buyRelease, err1 := db.GetBuyRelease(checkoutSession.ID)
		if err1 != nil {
			log.Println("router2/stripe_webhooks.go ERROR getting buy release in postStripeWebhook:", err1)
			// if we don't find thebuy release its because we were testing with the CLI and no buy release was created
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		release, err3 := db.GetAnyRelease(buyRelease.ReleaseID)
		if err3 != nil {
			log.Println("router2/stripe_webhooks.go ERROR getting release in postStripeWebhook:", err3)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// AmountPayed * PercentageShare
		// we took PercentageShare percent
		ourCut := payments.CalculateCut(release.Price)
		authorsCut := buyRelease.AmountPaying - uint64(ourCut)

		purchase := db.Purchase{
			UserID:          buyRelease.UserID,
			ReleaseID:       buyRelease.ReleaseID,
			StripeSessionID: buyRelease.StripeSessionID,
			StripePaymentID: buyRelease.StripePaymentID,
			CourseID:        release.CourseID,
			CreatedAt:       time.Now(),
			AmountPaid:      buyRelease.AmountPaying,
			AuthorsCut:      authorsCut,
		}
		err2 := db.CreatePurchase(&purchase)
		if err2 != nil {
			log.Println("router2/stripe_webhooks.go ERROR creating purchase in postStripeWebhook:", err2)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ownership := db.Ownership{
			UserID:    buyRelease.UserID,
			CourseID:  release.CourseID,
			ReleaseID: buyRelease.ReleaseID,
			Completed: false,
		}
		err5 := db.CreateOwnership(&ownership)
		if err5 != nil {
			log.Println("router2/stripe_webhooks.go ERROR creating ownership in postStripeWebhook:", err5)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if purchase creation succeeded
		// delete buyRelease to prevent user from re-buying for free
		err4 := db.DeleteBuyRelease(buyRelease.StripeSessionID)
		if err4 != nil {
			log.Println("failed to delete buyRelease (may have timed out and auto deleted):", err4)
		}

		log.Println("Payment success!")
		c.Writer.WriteHeader(http.StatusOK)
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	c.Writer.WriteHeader(http.StatusOK)
}
