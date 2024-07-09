package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

const VERSION = "1.0.0"

type config struct {
	addr string
	env  string
}

type application struct {
	cfg      *config
	infogLog *log.Logger
	errLog   *log.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", "localhost:3000", "API server address")
	flag.StringVar(&cfg.env, "env", "DEVELOPMENT", "Environment (DEVELOPMENT|STAGING|PRODUCTION)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		cfg:      &cfg,
		infogLog: infoLog,
		errLog:   errLog,
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     errLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	infoLog.Printf("Listening on %s\n", cfg.addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
