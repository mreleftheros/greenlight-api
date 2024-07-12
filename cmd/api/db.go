package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) initDb() error {
	db, err := pgxpool.New(context.Background(), app.cfg.db)
	if err != nil {
		return err
	}
	defer db.Close()

	app.infogLog.Print("Database connected")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return err
	}

	app.db = db

	return nil
}
