package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	server   = ""
	port     = ""
	user     = ""
	password = ""
	database = ""
)

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s",
		server, port, user, password, database)
	connString += " sslmode=disable"
	fmt.Println(connString)
	Db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
