package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var DB *sql.DB
var Redis *redis.Client

func Connect() {
	// ** DATABASE SETTINGS & CONNECT**
	//load environment variables
	psql_host := os.Getenv("POSTGRE_HOST")
	psql_user := os.Getenv("POSTGRE_USER")
	psql_password := os.Getenv("POSTGRE_PASSWORD")
	psql_dbname := os.Getenv("POSTGRE_DBNAME")
	psql_port := os.Getenv("POSTGRE_PORT")
	//redis_string := os.Getenv("REDIS_SECRET")

	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_dbname := os.Getenv("REDIS_DBNAME")

	redisAddr := fmt.Sprintf("%s:%s", redis_host, redis_port)
	dbname_number, _ := strconv.Atoi(redis_dbname)
	// set redis options
	opt := &redis.Options{
		Addr:     redisAddr,
		Password: redis_password,
		DB:       dbname_number,
	}
	// create redis connection
	redis := redis.NewClient(opt)

	// test redis server connection
	statusCmd := redis.Ping(context.Background())
	if statusCmd.Err() != nil {
		log.Fatal(statusCmd.Err().Error())
	}

	Redis = redis

	fmt.Println("Successful connect to redis database")

	// connect to postgres
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", psql_host, psql_user, psql_password, psql_dbname, psql_port)

	db, err := sql.Open("postgres", connection)

	DB = db

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successful connect to psql database")
}
