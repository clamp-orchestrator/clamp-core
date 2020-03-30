package repository

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"sync"
)

var (
	singletonDB   *pg.DB
	singletonOnce sync.Once
)

func connectDB() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	return db
}

func CloseDB() {
	if singletonDB != nil {
		fmt.Println("Disconnecting from DB")
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
