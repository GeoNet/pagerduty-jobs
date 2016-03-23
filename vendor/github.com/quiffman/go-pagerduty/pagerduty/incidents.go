package pagerduty

import "net/http"
import "time"

// IncidentsService type
type IncidentsService struct {
	client *Client
}

// Incident type
type Incident struct {
	ID                    string            `json:"id,omitempty"`
	IncidentNumber        int               `json:"incident_number,omitempty"`
	Status                string            `json:"status,omitempty"`
	CreatedOn             time.Time         `json:"created_on,omitempty"`
	Summary               *IncidentSummary  `json:"trigger_summary_data,omitempty"`
	User                  *User             `json:"assigned_to_user,omitempty"`
	SService               *Service          `json:"service,omitempty"` // This is conflicting with the package name on assignment in test. Not sure of the soltuion
	EEscalationPolicy      *EscalationPolicy `json:"escalation_policy,omitempty"` // This is conflicting with the package name on assignment in test. Not sure of the soltuion
	HTMLURL               string            `json:"html_url,omitempty"`
	IncidentKey           string            `json:"incident_key,omitempty"`
	TriggerDetailsHTMLURL string            `json:"trigger_details_html_url,omitempty"`
	TriggerType           string            `json:"trigger_type,omitempty"`
	LastStatusChangeOn    string            `json:"last_status_change_on,omitempty"`
	LastStatusChangeBy    *User             `json:"last_status_change_by,omitempty"`
	NumberOfEscalations   int               `json:"number_of_escalations,omitempty"`
	ResolvedByUser        *User             `json:"resolved_by_user,omitempty"`
	AssignedToUser        *User             `json:"assigned_to_user,omitempty"`
	AssignedTo            []*User           `json:"assigned_to,omitempty"`
}

// Incidents is a list of incidents
type Incidents struct {
	Pagination
	Incidents []Incident
}

// IncidentSummary type
type IncidentSummary struct {
	Subject     string //`json:"subject,omitempty"`
	Description string //`json:"description,omitempty"`
}

// Get returns a single incident by id if found
func (s *IncidentsService) Get(id string) (*Incident, *http.Response, error) {
	incident := new(Incident)

	res, err := s.client.Get("incidents/"+id, incident)
	if err != nil {
		return nil, res, err
	}

	return incident, res, nil
}

// IncidentsOptions provides optional parameters to list requests
type IncidentsOptions struct {
	Pagination
	Status         string `url:"status,omitempty"`
	SortBy         string `url:"sort_by,omitempty"`
	Since          string `url:"since,omitempty"`
	Until          string `url:"until,omitempty"`
	DateRange      string `url:"date_range,omitempty"`
	Service        string `url:"service,omitempty"`
	AssignedToUser string `url:"assigned_to_user,omitempty"`
}

// ListAll returns a list of incidents, by recursively calling List and utilising the pagination response
func (s *IncidentsService) ListAll(opt *IncidentsOptions) ([]Incident, error) {
	var incidents []Incident

	i, _, e := s.List(opt)
	if e == nil {
		incidents = i.Incidents
		opt.Pagination = i.Pagination
		opt.Pagination.Offset += opt.Pagination.Limit

		if opt.Pagination.Offset < opt.Pagination.Total {
			if i, e := s.ListAll(opt); e == nil {
				return append(incidents, i...), nil
			}
		}
	}
	return incidents, e
}

// List returns a list of incidents
func (s *IncidentsService) List(opt *IncidentsOptions) (Incidents, *http.Response, error) {
	var incidents Incidents
	u, err := addOptions("incidents", opt)
	if err != nil {
		return incidents, nil, err
	}


	res, err := s.client.Get(u, &incidents)
	if err != nil {
		return incidents, res, err
	}

	return incidents, res, err
}

// ReassignOptions provides optional parameters to incident reassign
type ReassignOptions struct {
	RequesterID      string `json:"requester_id,omitempty"`
	EscalationPolicy string `json:"escalation_policy,omitempty"`
	EscalationLevel  int    `json:"escalation_level,omitempty"`
	AssignedToUser   string `json:"assigned_to_user,omitempty"`
}

// Reassign an incident according to the options provided
func (s *IncidentsService) Reassign(id string, opt *ReassignOptions) (*http.Response, error) {
	res, err := s.client.Put("incidents/"+id+"/reassign", opt, nil)
	if err != nil {
		return res, err
	}

	return res, err
}
