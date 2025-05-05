// pd-fortnight-review lists last 14 days of incidents
package main

import (
	"flag"
	"log"
	"os"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var (
	authToken string = os.Getenv("PD_AUTHTOKEN")
	teamID    string = "PXRFYRO" //app support id
)

func init() {
	flag.StringVar(&authToken, "authtoken", authToken, "The authorization token for calls to the Pagerduty API.")

	flag.StringVar(&teamID, "teamID", teamID, "pagerduty team identifier to filter incidents from services")

	flag.Parse()

	if authToken == "" {
		log.Fatalln("PagerDuty auth token is required.")
	}
}

func main() {
	pd := pagerduty.NewClient(authToken)

	var serviceOpts pagerduty.ListServiceOptions
	serviceOpts.TeamIDs = []string{teamID}

	respService, err := pd.ListServices(serviceOpts)
	if err != nil {
		log.Fatalln("Failed to fetch services for given team filter: " + err.Error())
	}
	svcIDs := make([]string, len(respService.Services))
	for i, svc := range respService.Services {
		svcIDs[i] = svc.ID
	}

	var incidentOpts pagerduty.ListIncidentsOptions
	incidentOpts.Since = time.Now().AddDate(0, 0, -14).Format(time.RFC3339)
	incidentOpts.ServiceIDs = svcIDs
	log.Printf("opts: %v\n", incidentOpts)
	resp, err := pd.ListIncidents(incidentOpts)

	if err != nil {
		log.Fatalln("Failed to fetch incidents for given filter: " + err.Error())
	}

	log.Printf("found %d incidents\n", len(resp.Incidents))
	for _, i := range resp.Incidents {
		log.Printf("%v, %d, %s\n", i.CreatedAt, i.IncidentNumber, i.Summary)
	}

}

//https://geonet.pagerduty.com/teams/PXRFYRO/users

//loop through teams services
