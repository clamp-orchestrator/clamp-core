package repository

import (
	"clamp-core/config"
	"fmt"
)

func Ping() error {
	switch config.ENV.DBDriver {
	case "postgres":
		return pgDB.Ping()
	}
	return fmt.Errorf("Unsupported Database %s", config.ENV.DBDriver)
}
