package postgres

import (
	"app/config"
	"app/pkg/logger"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	conn *SQLDB
)

// GetConn ...
func GetConn() *SQLDB {
	connMx := &sync.Mutex{}
	connMx.Lock()
	defer connMx.Unlock()
	if conn == nil {
		var err error

		conn, err = Init()
		if err != nil {
			logger.Err("Error initializing DB", err)
			return nil
		}
	}

	return conn
}

// SQLDB ...
type SQLDB struct {
	Instance *sqlx.DB
}

func doInit() (*SQLDB, error) {
	addr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Get("HOST"), config.Get("DBPORT"), config.Get("USER"), config.Get("PASSWORD"), config.Get("NAME"))

	instance, err := sqlx.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db := SQLDB{
		Instance: instance,
	}

	return &db, nil
}

// Init inits DB
func Init() (*SQLDB, error) {
	var dbsess *SQLDB
	var err error

	dbsess, err = doInit()
	if err != nil {
		logger.Err("Failed to connect to Postgres", err)
		return nil, err
	}

	logger.Info("Connected to Postgres database")
	return dbsess, nil
}

// Close DB connection
func (db *SQLDB) Close() error {
	return db.Instance.Close()
}
