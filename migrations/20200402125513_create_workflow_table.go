package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE TABLE WORKFLOWS (" +
			"ID SERIAL NOT NULL, " +
			"NAME VARCHAR(20) NOT NULL, " +
			"ENABLED BOOLEAN DEFAULT TRUE, " +
			"DESCRIPTION VARCHAR(100), " +
			"CREATED_AT timestamp DEFAULT NOW(), " +
			"UPDATED_AT timestamp, " +
			"STEPS JSONB, " +
			"PRIMARY KEY(ID));")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE WORKFLOWS;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200402125513_create_workflow_table", up, down, opts)
}
