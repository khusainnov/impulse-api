package entity

type ClientMessage struct {
	Success bool   `json:"success,omitempty"`
	Data    []Data `json:"data,omitempty"`
}

type Data struct {
	Text string `json:"text,omitempty"`
}
