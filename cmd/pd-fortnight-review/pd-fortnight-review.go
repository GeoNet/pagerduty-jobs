// pd-fortnight-review lists last 14 days of incidents
package main

import (
	"flag"
	"fmt"
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
	incidentOpts.Until = time.Now().Format(time.RFC3339)
	incidentOpts.ServiceIDs = svcIDs
	incidentOpts.Limit = 100
	log.Printf("opts: %v\n", incidentOpts)
	resp, err := pd.ListIncidents(incidentOpts)

	if err != nil {
		log.Fatalln("Failed to fetch incidents for given filter: " + err.Error())
	}

	fmt.Printf("found %d incidents\n", len(resp.Incidents))
	page := 0
	for _, i := range resp.Incidents {
		fmt.Printf("%v, %d, %s\n", i.CreatedAt, i.IncidentNumber, i.Summary)
	}
	for resp.More {
		for _, i := range resp.Incidents {
			fmt.Printf("%v, %d, %s\n", i.CreatedAt, i.IncidentNumber, i.Summary)
		}
		page++
		fmt.Printf("page:%d\n", page)
		incidentOpts.Offset = incidentOpts.Offset + incidentOpts.Limit
		resp, err = pd.ListIncidents(incidentOpts)
		if err != nil {
			log.Fatalln("Failed to fetch incidents for given filter: " + err.Error())
		}
		log.Printf("paginated? %t", resp.More)
	}
	fmt.Println("printing schedules")

	var scheduleOpts pagerduty.ListSchedulesOptions
	scheduleOpts.Query = "Application"
	respSched, err := pd.ListSchedules(scheduleOpts)
	if err != nil {
		log.Fatalln("failed to read schedules")
	}

	var onCallOpts pagerduty.ListOnCallUsersOptions
	onCallOpts.Since = time.Now().AddDate(0, 0, -14).Format(time.RFC3339)
	onCallOpts.Until = time.Now().Format(time.RFC3339)

	for _, sched := range respSched.Schedules {
		respUser, err := pd.ListOnCallUsers(sched.ID, onCallOpts)
		if err != nil {
			log.Println("failed to read " + sched.ID)
			continue
		}
		log.Printf("Schedule: %s Username:? %s\n", sched.Summary, respUser[0].Name)
	}

}
