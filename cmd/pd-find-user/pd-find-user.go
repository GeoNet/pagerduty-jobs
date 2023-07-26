// pd-find-user finds and validates a user by ID, name or email.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/GeoNet/pagerduty-jobs/internal/finduser"
	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var (
	authToken string = os.Getenv("PD_AUTHTOKEN")
	user      string
)

func init() {
	flag.StringVar(&authToken, "authtoken", authToken, "The authorization token for calls to the Pagerduty API.")

	flag.StringVar(&user, "user", user, "The user ID, name or email to find and validate.")

	flag.Parse()

	if authToken == "" {
		log.Fatalln("PagerDuty auth token is required.")
	}

	if user == "" {
		log.Fatalln("Please specify a -user")
	}
}

func main() {
	pd := pagerduty.NewClient(authToken)

	us := finduser.Client{Client: *pd}
	u, err := us.FindAndValidate(context.Background(), user)
	if err != nil {
		log.Fatalln("Failed to validate user: " + err.Error())
	}
	log.Printf("Found from-user: %v, id: %v\n", u.Name, u.ID)
}
