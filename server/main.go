package main

import (
	"database/sql" // SQL database interactions
	"fmt"          // for formated I/O
	"log"

	_ "github.com/mattn/go-sqlite3"        // SQLite driver
	"github.com/vini464/wizard-duel/share" // shared stuff see share folder
)

var DB *sql.DB

func init_DB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

  // creating tables
  sqlStmt := `
  CREATE TABLE IF NOT EXISTS users (
  username TEXT NOT NULL PRIMARY KEY,
  password TEXT NOT NULL,
  coins INTEGER NOT NULL,

  )
  `
}

func main() {

}
