package models

type MessageContainer struct {
	Messages map[string]Message `json:"messages"`
}
