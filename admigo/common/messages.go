package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type message struct {
	Confirmed         string `json:"confirmed"`
	Regfin            string `json:"regfin"`
	Required          string `json:"required"`
	EmailNotFound     string `json:"email_not_found"`
	IncorrectPassword string `json:"incorrect_password"`
	Welcome           string `json:"welcome_logged"`
	UserNotFound      string `json:"user_not_found"`
	InsufRights       string `json:"insuf_rights"`
	UserNotConfirmed  string `json:"user_not_confirm"`
}

var msgs *message

func loadMessages() {
	file, err := os.Open(fmt.Sprintf("lang/%s/messages.json", Env().Lang))
	if err != nil {
		Danger("Cannot open messages file", err)
	}
	decoder := json.NewDecoder(file)
	msgs = &message{}
	err = decoder.Decode(msgs)
	if err != nil {
		Danger("Cannot get messages from file", err)
	}
}

func Mess() *message {
	return msgs
}
