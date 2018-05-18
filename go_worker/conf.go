package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type configuration struct {
	Infura           string `json:"infura"`
	Password         string `json:"password"`
	Abi              string `json:"abi"`
	CrowdsaleAddress string `json:"crowdsale_address"`
	Gas              uint64 `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Cmd              string `json:"cmd"`
	Key              string
}

var config *configuration

func confInit() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panic("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = &configuration{}
	err = decoder.Decode(config)
	if err != nil {
		log.Panic("Cannot get configuration from file", err)
	}
	loadKey()
}

func loadKey() {
	bf, err := ioutil.ReadFile("key")
	if err != nil {
		log.Panic("Cannot open key file", err)
	}
	config.Key = string(bf)
}

func env() *configuration {
	return config
}
