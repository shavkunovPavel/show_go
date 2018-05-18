package model

import (
	"admigo/common"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
)

var Db *sql.DB

func loadDriver() {
	var err error
	c := common.Env()
	connect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", c.Db.User, c.Db.Password, c.Db.Dbname, c.Db.Sslmode)
	Db, err = sql.Open(c.Db.Driver, connect)
	if err != nil {
		DangerDB("Error in loadDriver", err)
	}
	return
}

func FormToJson(body io.Reader, v interface{}) {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		DangerDB("Cannot parse form", err)
	}
}

/*
 * helper function returns format string (%s\n%s)
 * for writing many rows sql query
 */
func GetFormat(rows int) (format string) {
	if rows == 0 {
		return
	}
	for i := 0; i < rows; i++ {
		if len(format) > 0 {
			format += "\n"
		}
		format += "%s"
	}
	return
}
