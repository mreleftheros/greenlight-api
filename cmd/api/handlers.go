package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status: available\n")
	fmt.Fprintf(w, "environment: %s\n", app.cfg.env)
	fmt.Fprintf(w, "version: %s\n", VERSION)
}
