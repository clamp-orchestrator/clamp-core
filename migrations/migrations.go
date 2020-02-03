package migrations

import (
	"log"
	"os"

	"github.com/go-pg/pg/v9"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

const directory = "migrations"

//Migrate executes all the database migrations on the database
func Migrate() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "clamp",
		Database: "clampdev",
		Password: "clamppass",
	})

	err := migrations.Run(db, directory, os.Args)

	if err != nil {
		log.Fatalln(err)
	}
}
