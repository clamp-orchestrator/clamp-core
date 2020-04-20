package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS ADD COLUMN STEP_ID INT NOT NULL DEFAULT 0;")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS DROP COLUMN STEP_ID;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200420115810_add_step_id_col_to_steps_status_table", up, down, opts)
}
