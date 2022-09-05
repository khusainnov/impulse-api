package entity

type CheckVars struct {
	CheckAspectingID int         `json:"check_aspecting_id,omitempty"`
	CheckType        interface{} `json:"check_type,omitempty"`
	CheckAspectedID  int         `json:"check_aspected_id,omitempty"`
	Soed             string      `json:"soed,omitempty"`
	Body             string      `json:"body,omitempty"`
}
