package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"status":      "available",
			"environment": app.cfg.env,
			"version":     VERSION,
		}

		jsonRes(w, data, nil)
	})

	mux.HandleFunc("POST /v1/users/signup", app.usersSignupPost)
	mux.HandleFunc("POST /v1/users/login", app.usersLoginPost)
	mux.HandleFunc("POST /v1/users/logout", app.usersLogoutPost)
	mux.HandleFunc("GET /v1/users/me", app.usersMeGet)

	mux.HandleFunc("GET /v1/movies", app.moviesGet)
	mux.HandleFunc("POST /v1/movies", app.moviesPost)
	mux.HandleFunc("GET /v1/movies/{id}", app.moviesIdParamGet)
	mux.HandleFunc("PUT /v1/movies/{id}", app.moviesIdParamPut)
	mux.HandleFunc("DELETE /v1/movies/{id}", app.moviesIdParamDelete)

	return app.recoverPanic(mux)
}
