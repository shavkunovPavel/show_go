package model

import (
	"fmt"
	"log"
	"os"
)

var logger_db *log.Logger

func init() {
	file, err := os.OpenFile("admigo_db.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open admigo_db.log file", err)
	}
	err = file.Truncate(0)
	if err != nil {
		log.Fatalln("Failed to truncate admigo_db.log file", err)
	}
	logger_db = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	afterInit()
}

func afterInit() {
	loadDriver()
}

func DangerDB(step string, args ...interface{}) {
	pref := fmt.Sprintf("%s [%s] ", "ER", step)
	logger_db.SetPrefix(pref)
	logger_db.Println(args...)
}
