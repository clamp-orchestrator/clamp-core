package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS ADD COLUMN PAYLOAD JSONB;")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE SERVICE_REQUESTS DROP COLUMN PAYLOAD;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200408113700_add_payload_col_to_service_request_table", up, down, opts)
}
