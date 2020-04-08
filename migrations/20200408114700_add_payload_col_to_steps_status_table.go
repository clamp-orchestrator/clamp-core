package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS ADD COLUMN PAYLOAD JSONB;")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS DROP COLUMN PAYLOAD;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200408114700_add_payload_col_to_steps_status_table", up, down, opts)
}
