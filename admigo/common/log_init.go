package common

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("admigo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	err = file.Truncate(0)
	if err != nil {
		log.Fatalln("Failed to truncate log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	afterInit()
}

func afterInit() {
	loadConfig()
	loadMessages()
	loadMenu()
}

func Version() string {
	return "0.1"
}

func Danger(step string, args ...interface{}) {
	pref := fmt.Sprintf("%s [%s] ", "ER", step)
	logger.SetPrefix(pref)
	logger.Println(args...)
}
