// pd-list-incidents lists incidents according to a given filter set of incident options, e.g.
// ./pd-list-incidents -authtoken="y_token" -filter='{"Statuses":["triggered"]}'
// ./pd-list-incidents -authtoken="y_token" -filter='{"DateRange":"all"}'
// ./pd-list-incidents -authtoken="y_token" -filter='{"DateRange":"all","ServiceIDs":[PDAB123,PDABXYZ"]}'
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var (
	authToken string = os.Getenv("PD_AUTHTOKEN")
	filter    string
)

func init() {
	flag.StringVar(&authToken, "authtoken", authToken, "The authorization token for calls to the Pagerduty API.")

	flag.StringVar(&filter, "filter", filter, "The incident options to filter incidents.")

	flag.Parse()

	if authToken == "" {
		log.Fatalln("PagerDuty auth token is required.")
	}
}

func main() {
	pd := pagerduty.NewClient(authToken)

	var opts pagerduty.ListIncidentsOptions
	if err := json.Unmarshal([]byte(filter), &opts); err != nil {
		log.Fatalln("Failed to parse filter: " + err.Error())
	}

	log.Printf("opts: %v\n", opts)
	resp, err := pd.ListIncidents(opts)

	if err != nil {
		log.Fatalln("Failed to fetch incidents for given filter: " + err.Error())
	}

	log.Printf("found %d incidents\n", len(resp.Incidents))
	for _, i := range resp.Incidents {
		log.Printf("%v, %d, %s\n", i.CreatedAt, i.IncidentNumber, i.Summary)
	}

}
