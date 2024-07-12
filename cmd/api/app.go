package main

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mreleftheros/greenlight-api/internal/models"
)

type application struct {
	cfg        *config
	infogLog   *log.Logger
	errLog     *log.Logger
	db         *pgxpool.Pool
	movieModel *models.MovieModel
}

func (app *application) initLoggers() {
	app.infogLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func (app *application) initModels() {
	app.movieModel = &models.MovieModel{Db: app.db}
}
