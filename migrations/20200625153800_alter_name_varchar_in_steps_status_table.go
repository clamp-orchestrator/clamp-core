package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS ALTER COLUMN WORKFLOW_NAME TYPE varchar(40);")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE STEPS_STATUS ALTER COLUMN WORKFLOW_NAME TYPE varchar(20);")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200625153800_alter_name_varchar_in_steps_status_table", up, down, opts)
}
