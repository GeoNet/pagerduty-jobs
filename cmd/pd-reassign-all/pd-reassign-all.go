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
	toLevel   uint = 0
)

func init() {
	flag.StringVar(&authToken, "authtoken", authToken, "The authorization token for calls to the Pagerduty API.")

	flag.StringVar(&fromUser, "from-user", fromUser,
		`The user who is currently assigned incidents that should be reassigned elsewhere.
		Specify a PagerDuty user ID, or alternatively a name or email to query.`)

	flag.StringVar(&toUser, "to-user", toUser, "Another user ID (or name or email to query for), to whom the incidents should be assigned.")
	flag.UintVar(&toLevel, "to-level", toLevel, "The escalation level to reset this incident to.")

	flag.Parse()

	if authToken == "" {
		log.Fatalln("PagerDuty auth token is required.")
	}

	if fromUser == "" {
		log.Fatalln("Please specify a -from-user")
	}

	if (toUser == "" && toLevel == 0) || !(toUser != "" || toLevel != 0) {
		log.Fatalln("Please specify one of -to-user or -to-level.")
	}
}

func main() {
	pd := pagerduty.NewClient(authToken)

	us := finduser.Client{Client: *pd}
	u, err := us.FindAndValidate(fromUser)
	if err != nil {
		log.Fatalln("Failed to validate user: " + err.Error())
	}
	log.Printf("Found from-user: %v, id: %v\n", u.Name, u.ID)

	if toUser != "" {
		toU, err := us.FindAndValidate(toUser)
		if err != nil {
			log.Fatalln("Failed to validate user: " + err.Error())
		}

		log.Printf("Found to-user: %v, id: %v\n", toU.Name, toU.ID)
	}

	resp, err := pd.ListIncidents(pagerduty.ListIncidentsOptions{UserIDs: []string{u.ID}})

	if err != nil {
		log.Fatalln("Failed to fetch incidents for user: " + err.Error())
	}

	manage := make([]pagerduty.ManageIncidentsOptions, len(resp.Incidents))
	for i, incident := range resp.Incidents {
		log.Printf("incident: %v status: %v\n", incident.Id, incident.Status)
		log.Printf("Escap:%v\n", incident.EscalationPolicy.ID)

		// looking for current on-call user
		if toUser == "" {
			res, err := pd.ListOnCalls(pagerduty.ListOnCallOptions{EscalationPolicyIDs: []string{incident.EscalationPolicy.ID}})
			if err != nil {
				continue
			}

			var userid string
			for _, users := range res.OnCalls {
				toUser = users.User.ID
				if userid == u.ID {
					continue
				}
				if toLevel != 0 && users.EscalationLevel != toLevel {
					continue
				}
				log.Printf("user oncall:%v\n", toUser)
				log.Printf("user oncall:%v\n", users.EscalationLevel)
			}
			log.Printf("incident:%v\n", incident.Id)

			manage[i].Assignments = []pagerduty.Assignee{
				{
					Assignee: pagerduty.APIObject{
						ID:   toUser,
						Type: "user_reference",
					},
				},
			}
		}

		if toUser == "" {
			log.Fatalln("Could not find the correct User.")
		}

		manage[i].Assignments = []pagerduty.Assignee{
			{
				Assignee: pagerduty.APIObject{
					ID:   toUser,
					Type: "user_reference",
				},
			},
		}
		manage[i].ID = incident.Id
		manage[i].Type = "incident_reference"
	}

	_, err = pd.ManageIncidents(u.Email, (manage))
	if err != nil {
		log.Fatalln("Failed to manage incidents: " + err.Error())
	}

}
