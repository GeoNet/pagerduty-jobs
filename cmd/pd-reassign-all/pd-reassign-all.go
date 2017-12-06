// pd-reassign-all reassigns all Pagerduty incidents from one user to another, either directly or via re-escalating.
package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/GeoNet/pagerduty-jobs/internal/finduser"
	"github.com/quiffman/go-pagerduty/pagerduty"
)

var (
	InvalidValueError = errors.New("Invalid Value.")
)

var (
	subdomain string = os.Getenv("PD_SUBDOMAIN")
	apiKey    string = os.Getenv("PD_APIKEY")
	fromUser  string
	toUser    string
	toLevel   int = -1
)

func init() {
	flag.StringVar(&subdomain, "subdomain", subdomain, "The subdomain to be used for the Pagerduty API.")
	flag.StringVar(&apiKey, "api-key", apiKey, "The api-key authorized for calls to the Pagerduty API.")

	flag.StringVar(&fromUser, "from-user", fromUser,
		`The user who is currently assigned incidents that should be reassigned elsewhere.
		Specify a PagerDuty user ID, or alternatively a name or email to query.`)

	flag.StringVar(&toUser, "to-user", toUser, "Another user ID (or name or email to query for), to whom the incidents should be assigned.")
	flag.IntVar(&toLevel, "to-level", toLevel, "The escalation level to reset this incident to.")

	flag.Parse()

	if subdomain == "" || apiKey == "" {
		log.Fatalln("PagerDuty subdomain and API token are required.")
	}

	if fromUser == "" {
		log.Fatalln("Please specify a -from-user")
	}

	if (toUser == "" && toLevel == -1) || !(toUser != "" || toLevel != -1) {
		log.Fatalln("Please specify one of -to-user or -to-level.")
	}
}

func main() {
	pd := pagerduty.New(subdomain, apiKey)

	us := finduser.Service{*pd.Users}
	u, err := us.FindAndValidate(fromUser)
	if err != nil {
		log.Fatalln("Failed to validate user: " + err.Error())
	}
	log.Printf("Found from-user: %v, id: %v\n", u.Name, u.ID)

	var rOpts = pagerduty.ReassignOptions{RequesterID: u.ID}
	if toUser != "" {
		toU, err := us.FindAndValidate(toUser)
		if err != nil {
			log.Fatalln("Failed to validate user: " + err.Error())
		}

		log.Printf("Found to-user: %v, id: %v\n", toU.Name, toU.ID)
		rOpts.AssignedToUser = toU.ID
	}
	if toLevel != -1 {
		rOpts.EscalationLevel = toLevel
	}

	incidents, err := pd.Incidents.ListAll(&pagerduty.IncidentsOptions{AssignedToUser: u.ID})

	if err != nil {
		log.Fatalln("Failed to fetch incidents for user: " + err.Error())
	} else {
		for _, incident := range incidents {
			log.Printf("incident: %v status: %v\n", incident.ID, incident.Status)
			_, err := pd.Incidents.Reassign(incident.ID, &rOpts)
			if err != nil {
				log.Fatalln("Failed to reassign incident: " + err.Error())
			}
		}
	}
}
