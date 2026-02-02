package db

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	db *pgxpool.Pool
}

var (
	once sync.Once
	err  error = nil
)

// this function that creates the db pool can be used once
func NewDb() (*Db, error) {

	var database *Db = nil

	once.Do(func() {

		strconn := os.Getenv("PG_DEV_URI")
		if strconn == "" {

			err = errors.New("the connection string is not set")
			return

		}

		config, currErr := pgxpool.ParseConfig(strconn)
		if currErr != nil {
			err = currErr
			return

		}

		config.ConnConfig.ConnectTimeout = 5 * time.Second

		db, currErr := pgxpool.NewWithConfig(context.Background(), config)

		if currErr != nil {
			err = currErr
			return

		}

		database = &Db{

			db: db,
		}

	})
	return database, err

}

func (d *Db) GetDb() *pgxpool.Pool {
	return d.db

}

func (d *Db) Ping(ctx context.Context) error {
	return d.db.Ping(ctx)
}

func (d *Db) Close() {
	d.db.Close()
}
