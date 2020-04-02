package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE TABLE STEPS_STATUS ( ID SERIAL NOT NULL, SERVICE_REQUEST_ID UUID REFERENCES SERVICE_REQUESTS(ID), WORKFLOW_NAME VARCHAR(20) NOT NULL, " +
			"STATUS VARCHAR NOT NULL DEFAULT 'STARTED'," +
			"CREATED_AT timestamp DEFAULT NOW(), STEP_NAME VARCHAR(100) NOT NULL, REASON VARCHAR(500), PRIMARY KEY(ID));")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE STEPS_STATUS;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200402131000_create_steps_status_table", up, down, opts)
}
