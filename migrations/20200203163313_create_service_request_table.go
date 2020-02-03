package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE TABLE SERVICE_REQUESTS ( ID uuid, WORKFLOW_NAME VARCHAR NOT NULL, PRIMARY KEY(ID));")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE SERVICE_REQUESTS;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200203163313_create_service_request_table", up, down, opts)
}
