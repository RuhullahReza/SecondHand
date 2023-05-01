package db

import (
	"fmt"
	"log"
	"time"

	"github.com/RuhullahReza/SecondHand/util/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


func NewPostgresConnection() *sqlx.DB {

	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		config.DatabaseHost(), 
		config.DatabasePort(),
		config.DatabaseUsername(),
		config.DatabasePassword(),
		config.DatabaseName(),
		config.DatabaseMode(),
	)

	db, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func TestPostgresConnection() *sqlx.DB {

	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		config.DatabaseHost(), 
		config.DatabasePort(),
		config.DatabaseUsername(),
		config.DatabasePassword(),
		config.DatabaseName(),
		config.DatabaseMode(),
	)

	db, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
