package repository

import (
	"clamp-core/config"
	"context"
	"strings"
	"sync"

	"github.com/go-pg/pg/v9"

	log "github.com/sirupsen/logrus"
)

var singletonOnce sync.Once

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	log.Infof("[PSQL] Query: %v, Error: %v", query, err)
	return nil
}

func connectDB() (db *pg.DB) {
	db = pg.Connect(GetPostgresOptions())
	if config.ENV.EnableSQLQueriesLog {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

// GetPostgresOptions returns connection details for postgres DB
func GetPostgresOptions() *pg.Options {
	connStr := config.ENV.DBConnectionStr
	connArr := strings.Split(connStr, " ")
	var host, user, password, dbName string
	for _, conn := range connArr {
		connMap := strings.Split(conn, "=")
		switch connMap[0] {
		case "host":
			host = connMap[1]
		case "user":
			user = connMap[1]
		case "password":
			password = connMap[1]
		case "dbname":
			dbName = connMap[1]
		}
	}
	return &pg.Options{
		Addr:     host,
		User:     user,
		Password: password,
		Database: dbName,
	}
}

type postgresDB struct {
	db *pg.DB
}

var pgDB postgresDB

func init() {
	pgDB = postgresDB{}
}

func (p *postgresDB) GetDB() *pg.DB {
	singletonOnce.Do(func() {
		log.Info("Connecting to DB")
		p.db = connectDB()
	})
	return p.db
}

func (p *postgresDB) Ping() error {
	_, err := p.GetDB().Exec("SELECT 1")
	return err
}

func (p *postgresDB) closeDB() {
	if p.db != nil {
		log.Info("Disconnecting from DB")
		p.db.Close()
	}
}
