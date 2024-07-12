package main

import (
	"net/http"
	"time"
)

func (app *application) initServer() error {

	srv := &http.Server{
		Addr:         app.cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     app.errLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.infogLog.Printf("Starting %s on %s\n", app.cfg.env, app.cfg.addr)
	return srv.ListenAndServe()
}
