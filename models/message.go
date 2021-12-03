package models

type Message struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Time   int64  `json:"time"`
}
