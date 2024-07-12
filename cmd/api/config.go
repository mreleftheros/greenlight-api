package main

import "flag"

type config struct {
	addr string
	env  string
	db   string
}

func (app *application) NewConfig() {
	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", "localhost:3000", "API server address")
	flag.StringVar(&cfg.env, "env", "DEVELOPMENT", "Environment (DEVELOPMENT|STAGING|PRODUCTION)")
	flag.StringVar(&cfg.db, "db", "postgresql://postgres@localhost/greenlight", "database string")
	flag.Parse()

	app.cfg = &cfg
}
