package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE UNIQUE INDEX WORKFLOW_NAME_INDEX ON WORKFLOWS (NAME);")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP INDEX WORKFLOW_NAME_INDEX;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200402131700_create_workflow_name_unique_index", up, down, opts)
}
