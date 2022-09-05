package entity

type HouseUpr struct {
	Hoe []MainStruct `json:"houseUsp,omitempty"`
}

type MainStruct struct {
	House     int    `json:"house,omitempty"`
	Sign      string `json:"sign,omitempty"`
	Upr       string `json:"upr,omitempty"`
	Soupr     string `json:"soupr,omitempty"`
	HouseSymb string `json:"house_symb,omitempty"`
}

type ResponseUpr struct {
	House   int       `json:"house,omitempty"`
	Sign    string    `json:"sign,omitempty"`
	Upr     string    `json:"upr,omitempty"`
	Aspects []Aspects `json:"aspects,omitempty"`
}
