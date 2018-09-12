// pd-reassign-all reassigns all Pagerduty incidents from one user to another, either directly or via re-escalating.
package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/GeoNet/pagerduty-jobs/internal/finduser"
	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var (
	InvalidValueError = errors.New("Invalid Value.")
)

var (
	authToken string = os.Getenv("PD_AUTHTOKEN")
	fromUser  string
	toUser    string
	toLevel   int = -1
)

func init() {
	flag.StringVar(&authToken, "authtoken", authToken, "The authorization token for calls to the Pagerduty API.")

	flag.StringVar(&fromUser, "from-user", fromUser,
		`The user who is currently assigned incidents that should be reassigned elsewhere.
		Specify a PagerDuty user ID, or alternatively a name or email to query.`)

	flag.StringVar(&toUser, "to-user", toUser, "Another user ID (or name or email to query for), to whom the incidents should be assigned.")
	flag.IntVar(&toLevel, "to-level", toLevel, "The escalation level to reset this incident to.")

	flag.Parse()

	if authToken == "" {
		log.Fatalln("PagerDuty auth token is required.")
	}

	if fromUser == "" {
		log.Fatalln("Please specify a -from-user")
	}

	if (toUser == "" && toLevel == -1) || !(toUser != "" || toLevel != -1) {
		log.Fatalln("Please specify one of -to-user or -to-level.")
	}
}

func main() {
	pd := pagerduty.NewClient(authToken)

	us := finduser.Client{*pd}
	u, err := us.FindAndValidate(fromUser)
	if err != nil {
		log.Fatalln("Failed to validate user: " + err.Error())
	}
	log.Printf("Found from-user: %v, id: %v\n", u.Name, u.ID)

	var manageIncident pagerduty.Incident
	manageIncident.Type = "incident_reference"
	if toUser != "" {
		toU, err := us.FindAndValidate(toUser)
		if err != nil {
			log.Fatalln("Failed to validate user: " + err.Error())
		}

		log.Printf("Found to-user: %v, id: %v\n", toU.Name, toU.ID)
		manageIncident.Assignments = []pagerduty.Assignment{
			pagerduty.Assignment{
				Assignee: pagerduty.APIObject{
					ID:   toU.ID,
					Type: "user_reference",
				},
			},
		}
	}
	if toLevel != -1 {
		manageIncident.EscalationLevel = uint(toLevel)
	}

	resp, err := pd.ListIncidents(pagerduty.ListIncidentsOptions{UserIDs: []string{u.ID}})

	if err != nil {
		log.Fatalln("Failed to fetch incidents for user: " + err.Error())
	}

	manage := make([]pagerduty.Incident, len(resp.Incidents))
	for i, incident := range resp.Incidents {
		log.Printf("incident: %v status: %v\n", incident.ID, incident.Status)
		manage[i] = manageIncident
		manage[i].ID = incident.ID
	}

	err = pd.ManageIncidents(u.Email, []pagerduty.Incident(manage))
	if err != nil {
		log.Fatalln("Failed to manage incidents: " + err.Error())
	}

}
