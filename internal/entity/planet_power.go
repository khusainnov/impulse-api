package entity

type PlanetPower struct {
	Params []Parameters `json:"planetPower,omitempty"`
}

type Parameters struct {
	Planet    string `json:"planet,omitempty"`
	House     string `json:"house (cell),omitempty"`
	Elevation string `json:"elevation (exaltation),omitempty"`
	Exile     string `json:"exile,omitempty"`
	Fall      string `json:"fall,omitempty"`
}
