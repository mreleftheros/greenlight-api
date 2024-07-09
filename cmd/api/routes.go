package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthGet)

	mux.HandleFunc("GET /v1/movies", app.moviesGet)
	mux.HandleFunc("POST /v1/movies", app.moviesPost)
	mux.HandleFunc("GET /v1/movies/{id}", app.moviesIdParamGet)
	mux.HandleFunc("PUT /v1/movies/{id}", app.moviesIdParamPut)
	mux.HandleFunc("DELETE /v1/movies/{id}", app.moviesIdParamDelete)

	return mux
}
