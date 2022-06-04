package ilert

// Status-Page definition https://api.ilert.com/api-docs/#!/Status-Pages
type StatusPage struct {
	ID                        int64     `json:"id"`
	Name                      string    `json:"name"`
	Domain                    string    `json:"domain"`
	Subdomain                 string    `json:"subdomain"`
	CustomCss                 string    `json:"customCss"`
	FaviconUrl                string    `json:"faviconUrl"`
	LogoUrl                   string    `json:"logoUrl"`
	Visibility                string    `json:"visibility"`
	HiddenFromSearch          string    `json:"hiddenFromSearch"`
	ShowSubscribeAction       string    `json:"showSubscribeAction"`
	ShowIncidentHistoryOption string    `json:"showIncidentHistoryOption"`
	PageTitle                 string    `json:"pageTitle"`
	PageDescription           string    `json:"pageDescription"`
	LogoRedirectUrl           string    `json:"logoRedirectUrl"`
	Activated                 string    `json:"activated"`
	Status                    string    `json:"status"`
	Teams                     []TeamRel `json:"teams"`
}

//TeamRel defines teams
type TeamRel struct {
}

// StatusPage defines status-page visibility
var StatusPageVisibility = struct {
	Public  string
	Private string
}{
	Public:  "PUBLIC",
	Private: "PRIVATE",
}

var StatusPageStatus = struct {
	Operational      string
	UnderMaintenance string
	Degraded         string
	PartialOutage    string
	MajorOutage      string
}{
	Operational:      "OPERATIONAL",
	UnderMaintenance: "UNDER_MAINTENANCE",
	Degraded:         "DEGRADED",
	PartialOutage:    "PARTIAL_OUTAGE",
	MajorOutage:      "MAJOR_OUTAGE",
}
