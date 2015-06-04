// pd-list-incidents lists incidents according to a given filter set of incident options, e.g.
// ./pd-list-incidents -subdomain="subdomain" -api-key="api-key" -filter='{"Status":"triggered"}'
// ./pd-list-incidents -subdomain="subdomain" -api-key="api-key" -filter='{"DateRange":"all"}'
// ./pd-list-incidents -subdomain="subdomain" -api-key="api-key" -filter='{"DateRange":"all","Service":"PDAB123,PDABXYZ"}'
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/quiffman/go-pagerduty/pagerduty"
)

var (
	subdomain string
	apiKey    string
	filter    string
)

func init() {
	flag.StringVar(&subdomain, "subdomain", subdomain, "The subdomain to be used for the Pagerduty API.")
	flag.StringVar(&apiKey, "api-key", apiKey, "The api-key authorized for calls to the Pagerduty API.")

	flag.StringVar(&filter, "filter", filter, "The incident options to filter incidents.")

	flag.Parse()

	if subdomain == "" || apiKey == "" {
		log.Fatalln("PagerDuty subdomain and API token are required.")
	}
}

func main() {
	pd := pagerduty.New(subdomain, apiKey)

	var opts pagerduty.IncidentsOptions
	if err := json.Unmarshal([]byte(filter), &opts); err != nil {
		log.Fatalln("Failed to parse filter: " + err.Error())
	}

	incidents, err := pd.Incidents.ListAll(&opts)

	fmt.Printf("found %d incidents\n", len(incidents))
	if err != nil {
		log.Fatalln("Failed to fetch incidents for given filter: " + err.Error())
	} else {
		for _, i := range incidents {
			fmt.Printf("%v, %d, %s\n", i.CreatedOn.Format(time.RFC3339), i.IncidentNumber, i.Summary.Description)
		}
	}

}
