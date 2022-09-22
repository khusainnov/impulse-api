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
	Planets    []Planets      `json:"planets,omitempty"`
	House      int            `json:"house,omitempty"`
	Sign       string         `json:"sign,omitempty"`
	Upr        string         `json:"upr,omitempty"`
	Power      interface{}    `json:"power"`
	Aspects    []Aspects      `json:"aspects,omitempty"`
	AllElems   map[string]int `json:"allElems,omitempty"`
	PrevVal    PrevVal        `json:"prevVal,omitempty"`
	TestElems  string         `json:"testElems,omitempty"`
	AllCrests  map[string]int `json:"allCrests,omitempty"`
	PrevCrest  PrevCrest      `json:"prevCrest,omitempty"`
	TestCrests string         `json:"testCrests,omitempty"`
	RespMsg    string         `json:"respMsg,omitempty"`
}
