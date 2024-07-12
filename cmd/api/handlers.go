package main

import (
	"net/http"
	"strconv"

	"github.com/mreleftheros/greenlight-api/internal/models"
)

func (app *application) healthGet(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.cfg.env,
		"version":     VERSION,
	}

	jsonRes(w, data, nil)
}

func (app *application) moviesGet(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	mq := &models.MovieQuery{}
	if errors, ok := mq.Validate(&values); !ok {
		errRes(w, errors, nil)
		return
	}

	movies, metadata, err := app.movieModel.GetAll(mq)
	if err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}
	
	jsonRes(w, struct{Movies []*models.Movie`json:"movies"`;Metadata *models.Metadata`json:"metadata"`;}{Movies: movies, Metadata: metadata}, nil)
}

func (app *application) moviesPost(w http.ResponseWriter, r *http.Request) {
	mb := models.MovieBody{}
	if err := jsonBody(r, &mb); err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}

	mv := &models.Movie{}
	if mb.Title != nil {
		mv.Title = *mb.Title
	}
	if mb.Year != nil {
		mv.Year = *mb.Year
	}
	if mb.Runtime != nil {
		mv.Runtime = *mb.Runtime
	}
	if mb.Genres != nil {
		mv.Genres = mb.Genres
	}

	if errors, ok := app.movieModel.Validate(mv); !ok {
		errRes(w, errors, nil)
		return
	}

	if err := app.movieModel.Set(mv); err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil, 500)
		return
	}

	jsonRes(w, mv, nil, 201)
}

func (app *application) moviesIdParamGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		errRes(w, map[string]string{"error": "id not found"}, nil, 404)
		return
	}

	mv, err := app.movieModel.Get(id)
	if err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil, 400)
		return
	}

	jsonRes(w, mv, nil)
}

func (app *application) moviesIdParamPut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		errRes(w, map[string]string{"error": "id not found"}, nil)
		return
	}

	mb := models.MovieBody{}
	if err := jsonBody(r, &mb); err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}

	mv, err := app.movieModel.Get(id)
	if err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}

	if mb.Title != nil {
		mv.Title = *mb.Title
	}
	if mb.Year != nil {
		mv.Year = *mb.Year
	}
	if mb.Runtime != nil {
		mv.Runtime = *mb.Runtime
	}
	if mb.Genres != nil {
		mv.Genres = mb.Genres
	}

	errors, ok := app.movieModel.Validate(mv)
	if !ok {
		errRes(w, errors, nil)
		return
	}

	err = app.movieModel.Update(mv, id)
	if err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}

	jsonRes(w, mv, nil)
}

func (app *application) moviesIdParamDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		errRes(w, map[string]string{"error": "id not found"}, nil)
		return
	}

	if err = app.movieModel.Delete(id); err != nil {
		errRes(w, map[string]string{"error": err.Error()}, nil)
		return
	}

	jsonRes(w, nil, nil, 204)
}
