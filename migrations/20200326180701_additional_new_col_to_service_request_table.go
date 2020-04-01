package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS ADD COLUMN start_time timestamptz NOT NULL DEFAULT now(), ADD COLUMN end_time timestamp, ADD COLUMN total_time_elapsed_ms int, ADD COLUMN steps jsonb;")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS DROP COLUMN start_time,DROP COLUMN end_time,DROP COLUMN total_time_elapsed_ms,DROP COLUMN steps;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200326180701_additional_new_col_to_service_request_table", up, down, opts)
}
