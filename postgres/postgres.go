package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	driver string = "postgres"
)

type DB struct {
	conn *sql.DB
}

func DefaultConnection() DB {
	var d DB
	d.conn = createDefaultConnection()
	return d
}

func (d *DB) MustPrepare(statement string) *sql.Stmt {
	prep, err := d.conn.Prepare(statement)
	if err != nil {
		log.Fatalf("failed to prepare SQL %s:%v", statement, err)
	}
	return prep
}

func createDefaultConnection() *sql.DB {

	user := "postgres"
	password := "postgres"
	host := "localhost"
	port := 5400
	database := "roulette"

	url := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err := sql.Open(driver, url)
	if err != nil {
		log.Fatalf("Error connecting to the db %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database %s", err)

	}
	return db

}

func CreateCustomDB(db *sql.DB) DB {
	return DB{db}
}
