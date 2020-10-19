package data

import (
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DB struct {
	l  *logrus.Logger
	R  *redis.Client
	Ms *sql.DB
}

func New(l *logrus.Logger, r *redis.Client, ms *sql.DB) (*DB, error) {
	return &DB{l, r, ms}, nil
}

func (db DB) Close() error {
	if err := db.Ms.Close(); err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}
