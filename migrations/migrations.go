package migrations

import (
	"clamp-core/repository"
	"os"

	"github.com/go-pg/pg/v9"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	log "github.com/sirupsen/logrus"
)

const directory = "migrations"

//Migrate executes all the database migrations on the database
func Migrate() {
	options := repository.GetPostgresOptions()
	db := pg.Connect(options)

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalf("migration failed: %s", err)
	}
}
