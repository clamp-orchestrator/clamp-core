package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS ADD COLUMN CREATED_AT timestamp DEFAULT NOW();")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS DROP COLUMN CREATED_AT;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200402130000_add_created_at_col_to_service_request_table", up, down, opts)
}
