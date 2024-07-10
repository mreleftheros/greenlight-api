package models

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Year    int       `json:"year"`
	Runtime int       `json:"runtime"`
	Genres  []string  `json:"genres"`
	Created time.Time `json:"created"`
}

type MovieBody struct {
	Title   *string   `json:"title"`
	Year    *int      `json:"year"`
	Runtime *int      `json:"runtime"`
	Genres  []string `json:"genres"`
}

type MovieModel struct {
	Db *pgxpool.Pool
}

func (m *MovieModel) Validate(mv *Movie) (map[string]string, bool) {
	errors := make(map[string]string)

	if strings.TrimSpace(mv.Title) == "" {
		errors["titleError"] = "Title cannot be empty"
	} else if utf8.RuneCountInString(mv.Title) > 100 {
		errors["titleError"] = "Title length must be up to 100 characters"
	}

	if mv.Year < 0 {
		errors["yearError"] = "Year must be positive"
	} else if mv.Year == 0 {
		errors["yearError"] = "Year cannot be empty"
	} else if mv.Year < 1888 {
		errors["yearError"] = "Year must be minimum 1888"
	}

	if mv.Runtime < 0 {
		errors["runtimeError"] = "Runtime must be positive"
	} else if mv.Runtime == 0 {
		errors["runtimeError"] = "Runtime cannot be empty"
	}

	if len(mv.Genres) > 5 {
		errors["genresError"] = "No more than 5 genres"
	}

	genres := make(map[string]bool)
	for _, genre := range mv.Genres {
		genres[genre] = true
	}

	if len(genres) != len(mv.Genres) {
		if _, prs := errors["genresError"]; !prs {
			errors["genresError"] = "Genres must be unique"
		} else {
			errors["genresError"] += " and Genres must be unique"
		}
	}

	if len(errors) > 0 {
		errors["error"] = "Validation failed"

		return errors, false
	}

	return errors, true
}

func (m *MovieModel) Set(mv *Movie) error {
	stmt := "INSERT INTO movies (title, year, runtime, genres) VALUES($1, $2, $3, $4) RETURNING id, created;"

	if err := m.Db.QueryRow(context.Background(), stmt, mv.Title, mv.Year, mv.Runtime, mv.Genres).Scan(&mv.Id, &mv.Created); err != nil {
		return err
	}

	return nil
}

func (m *MovieModel) Get(id int) (*Movie, error) {
	stmt := "SELECT id, title, year, runtime, genres, created FROM movies WHERE id = $1;"
	
	mv := &Movie{}
	if err := m.Db.QueryRow(context.Background(), stmt, id).Scan(&mv.Id, &mv.Title, &mv.Year, &mv.Runtime, &mv.Genres, &mv.Created); err != nil {
		return nil, err
	}
	
	return mv, nil
}

func (m *MovieModel) Update(mv *Movie, id int) error {
	stmt := "UPDATE movies SET title = $1, year = $2, runtime = $3, genres = $4 WHERE id = $5;"

	result, err := m.Db.Exec(context.Background(), stmt, mv.Title, mv.Year, mv.Runtime, mv.Genres, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() < 1 {
		return errors.New("update failed")
	}

	return nil
}

func (m *MovieModel) Delete(id int) error {
	stmt := "DELETE FROM movies WHERE id = $1;"

	result, err := m.Db.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return errors.New("row not found")
	}

	return nil
}