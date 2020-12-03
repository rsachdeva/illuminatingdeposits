package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/dbconn"
)

func Db(cfg AppConfig) (*sqlx.DB, error) {
	db, err := dbconn.Open(dbconn.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
