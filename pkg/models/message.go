package models

import "time"

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
	Time time.Time `json:"time"`
}
