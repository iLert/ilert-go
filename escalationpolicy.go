package ilert

// EscalationPolicy definition https://api.ilert.com/api-docs/#!/Escalation_Policies
type EscalationPolicy struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	EscalationRules []EscalationRule `json:"escalationRules"`
	Repeating       bool             `json:"repeating"`
	Frequency       int              `json:"frequency"`
}

// EscalationRule definition
type EscalationRule struct {
	User              UserShort `json:"user"`
	Schedule          Schedule  `json:"schedule"`
	EscalationTimeout int       `json:"escalationTimeout"`
	Type              string    `json:"type"`
}

// Schedule definition
type Schedule struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	ScheduleMemberIDs []int64 `json:"scheduleMemberIds"`
}
