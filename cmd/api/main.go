package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mreleftheros/greenlight-api/internal/models"
)

const VERSION = "1.0.0"

type config struct {
	addr string
	env  string
	db string
}

type application struct {
	cfg      *config
	infogLog *log.Logger
	errLog   *log.Logger
	movieModel *models.MovieModel
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", "localhost:3000", "API server address")
	flag.StringVar(&cfg.env, "env", "DEVELOPMENT", "Environment (DEVELOPMENT|STAGING|PRODUCTION)")
	flag.StringVar(&cfg.db, "db", "postgresql://postgres@localhost/greenlight", "database string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbPool, err := pgxpool.New(context.Background(), cfg.db)
	if err != nil {
		errLog.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()
	infoLog.Print("Created connection pool successfully")

	if err = dbPool.Ping(context.Background()); err != nil {
		errLog.Fatal(err)
	}

	app := &application{
		cfg:      &cfg,
		infogLog: infoLog,
		errLog:   errLog,
		movieModel: &models.MovieModel{Db: dbPool},
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     errLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	infoLog.Printf("Starting %s on %s\n", cfg.env, cfg.addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
