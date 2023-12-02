package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func Connect() *sql.DB {
	// ** DATABASE SETTINGS & CONNECT**
	//load environment variables
	godotenv.Load()
	psql_host := os.Getenv("POSTGRE_HOST")
	psql_user := os.Getenv("POSTGRE_USER")
	psql_password := os.Getenv("POSTGRE_PASSWORD")
	psql_dbname := os.Getenv("POSTGRE_DBNAME")
	psql_port := os.Getenv("POSTGRE_PORT")

	// connect to postgres
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", psql_host, psql_user, psql_password, psql_dbname, psql_port)

	db, err := sql.Open("postgres", connection)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
