// pd-find-user finds and validates a user by ID, name or email.
package main

import (
	"flag"
	"log"

	"github.com/GeoNet/pagerduty-jobs/finduser"
	"github.com/quiffman/go-pagerduty/pagerduty"
)

var (
	subdomain string
	apiKey    string
	user      string
)

func init() {
	flag.StringVar(&subdomain, "subdomain", subdomain, "The subdomain to be used for the Pagerduty API.")
	flag.StringVar(&apiKey, "api-key", apiKey, "The api-key authorized for calls to the Pagerduty API.")

	flag.StringVar(&user, "user", user, "The user ID, name or email to find and validate.")

	flag.Parse()

	if subdomain == "" || apiKey == "" {
		log.Fatalln("PagerDuty subdomain and API token are required.")
	}

	if user == "" {
		log.Fatalln("Please specify a -user")
	}
}

func main() {
	pd := pagerduty.New(subdomain, apiKey)

	us := finduser.Service{*pd.Users}
	u, err := us.FindAndValidate(user)
	if err != nil {
		log.Fatalln("Failed to validate user: " + err.Error())
	}
	log.Printf("Found from-user: %v, id: %v\n", u.Name, u.ID)
}
