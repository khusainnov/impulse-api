package entity

type Summary struct {
	Planets   []Planets `json:"planets,omitempty"`
	Houses    []Houses  `json:"houses,omitempty"`
	Ascendant float64   `json:"ascendant,omitempty"`
	Midheaven float64   `json:"midheaven,omitempty"`
	Vertex    float64   `json:"vertex,omitempty"`
	Aspects   []Aspects `json:"aspects,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type Planets struct {
	Name       string      `json:"name,omitempty"`
	FullDegree float64     `json:"full_degree,omitempty"`
	NormDegree float64     `json:"norm_degree,omitempty"`
	Speed      float64     `json:"speed,omitempty"`
	IsRetro    string      `json:"is_retro,omitempty"`
	SignID     int         `json:"sign_id,omitempty"`
	Sign       string      `json:"sign,omitempty"`
	House      int         `json:"house,omitempty"`
	Element    string      `json:"element,omitempty"`
	Crest      string      `json:"crest,omitempty"`
	Burred     string      `json:"burred,omitempty"`
	Power      interface{} `json:"power,omitempty"`
}

type Houses struct {
	House  int     `json:"house,omitempty"`
	Sign   string  `json:"sign,omitempty"`
	Degree float64 `json:"degree,omitempty"`
}

type Aspects struct {
	AspectingPlanet   string  `json:"aspecting_planet,omitempty"`
	AspectedPlanet    string  `json:"aspected_planet,omitempty"`
	AspectingPlanetID int     `json:"aspecting_planet_id,omitempty"`
	AspectedPlanetID  int     `json:"aspected_planet_id,omitempty"`
	Type              string  `json:"type,omitempty"`
	Orb               float64 `json:"orb,omitempty"`
	Diff              float64 `json:"diff,omitempty"`
}
