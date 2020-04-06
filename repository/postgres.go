package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"log"
	"sync"
)

const LogSQLQueries bool = true

var (
	singletonDB   *pg.DB
	singletonOnce sync.Once
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	log.Printf("[PSQL] Query: %v, Error: %v", query, err)
	return nil
}

func connectDB() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	if LogSQLQueries {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

func CloseDB() {
	if singletonDB != nil {
		log.Println("Disconnecting from DB")
		singletonDB.Close()
	}
}

// GetDB gets the db singleton
func GetDB() *pg.DB {
	singletonOnce.Do(func() {
		fmt.Println("Connecting to DB")
		singletonDB = connectDB()
	})
	return singletonDB
}
