package config

import "github.com/caarlos0/env"

var Config = struct {

	/**
	Examples:

	CLAMP_DB_DBDRIVER     CLAMP_DB_DBCONNECTIONSTR
	=================     ===============================================================
	"postgres"            "host=localhost user=root dbname=clamp password=mypassword"

	*/
	DBDriver        string `env:"CLAMP_DB_DBDRIVER" envDefault:"postgres"`
	DBConnectionStr string `env:"CLAMP_DB_DBCONNECTIONSTR" envDefault:"host=myhost user=root dbname=clamp password=mypassword"`
}{}

func init() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}
