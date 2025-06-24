package models

type Message struct {
	Text       string      `json:"text"`
	ButtonRows []ButtonRow `json:"button_rows"`
}
