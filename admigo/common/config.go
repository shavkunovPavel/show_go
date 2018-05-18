package common

import (
	"encoding/json"
	"os"
)

type configuration struct {
	Lang         string      `json:"lang"`
	Debug        bool        `json:"debug"`
	Address      string      `json:"address"`
	Port         int         `json:"port"`
	Rpc          string      `json:"rpc"`
	ReadTimeout  int64       `json:"readTimeout"`
	WriteTimeout int64       `json:"writeTimeout"`
	IdleTimeout  int64       `json:"idleTimeout"`
	Static       string      `json:"static"`
	Db           *dbConfig   `json:"db"`
	Mail         *mailConfig `json:"mail"`
	Redirects    *[]redirect `json:"redirects,omitempty"`
}
type dbConfig struct {
	Driver   string `json:driver`
	User     string `json:user`
	Password string `json:password`
	Dbname   string `json:dbname`
	Sslmode  string `json:sslmode`
}
type mailConfig struct {
	From     string `json:"from"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	GotoUrl  string `json:"gotourl"`
}

type redirect struct {
	Prefix   string `json:"prefix"`
	Protocol string `json:"protocol"`
	ReqUri   string `json:"req_uri"`
	Port     int    `json:"port"`
}

var config *configuration

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		Danger("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = &configuration{}
	err = decoder.Decode(config)
	if err != nil {
		Danger("Cannot get configuration from file", err)
	}
}

func Env() *configuration {
	return config
}
