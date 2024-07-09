package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mreleftheros/greenlight-api/internal/models"
)

func (app *application) healthGet(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.cfg.env,
		"version":     VERSION,
	}

	err := jsonRes(w, data, nil)
	if err != nil {
		app.errLog.Println(err)
		http.Error(w, "Server encountered a problem", 500)
		return
	}
}

func (app *application) moviesGet(w http.ResponseWriter, r *http.Request) {
	errRes(w, "oh no", nil)
}

func (app *application) moviesPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create new movie")
}

func (app *application) moviesIdParamGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := models.Movie{
		Id:      id,
		Created: time.Now(),
		Title:   "Casablanca",
		Year:    0,
		Runtime: 122,
		Genres:  []string{"Comedy", "Drama", "Crime"},
	}

	jsonRes(w, movie, nil)
}

func (app *application) moviesIdParamPut(w http.ResponseWriter, r *http.Request) {}

func (app *application) moviesIdParamDelete(w http.ResponseWriter, r *http.Request) {}
