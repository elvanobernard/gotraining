package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func GetDb(user, password, net, address, dbName string) *sql.DB {
	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  net,
		Addr:                 address,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	return db
}
