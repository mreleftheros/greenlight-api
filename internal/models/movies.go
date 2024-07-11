package models

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
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
	Version int `json:"-"`
}

type MovieBody struct {
	Title   *string   `json:"title"`
	Year    *int      `json:"year"`
	Runtime *int      `json:"runtime"`
	Genres  []string `json:"genres"`
}

type MovieQuery struct {
	Title string
	Genres []string
	Page int
	PageSize int
	Sort string
}

type MovieModel struct {
	Db *pgxpool.Pool
}

func (mq *MovieQuery) Validate(values *url.Values) (map[string]string, bool) {
	errors := make(map[string]string)
	sortValues := [8]string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	title := values.Get("title")
	if title != "" {
		mq.Title = title
	} else {
		mq.Title = ""
	}

	genres := values.Get("genres")
	if genres != "" {
		mq.Genres = strings.Split(genres, ",")
	} else {
		mq.Genres = []string{}
	}

	page := values.Get("page")
	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			errors["pageError"] = err.Error()
		} else if p < 1 {
			errors["pageError"] = "page must be greater than zero"
		} else if p > 10000000 {
			errors["pageError"] = "page must be less or equal than 10 million"
		}
		mq.Page = p
	} else {
		mq.Page = 1
	}

	pageSize := values.Get("page_size")
	if pageSize != "" {
		ps, err := strconv.Atoi(pageSize)
		if err != nil {
			errors["pageSizeError"] = err.Error()
		} else if ps < 1 {
			errors["pageSizeError"] = "page_size must be greater than zero"
		} else if ps > 100 {
			errors["pageSizeError"] = "page_size must be less or equal than 100"
		}
		mq.PageSize = ps
	} else {
		mq.PageSize = 20
	}

	sort := values.Get("sort")
	if sort != "" {
		matched := false

		for _, item := range sortValues {
			if item == sort {
				mq.Sort = sort
				matched = true
				break
			}
		}
		if !matched {
			errors["sortError"] = "Sort value is not valid"
		}
	} else {
		mq.Sort = "-id"
	}

	if len(errors) > 0 {
		errors["error"] = "Validation failed"
		return errors, false
	}

	return nil, true
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

func (m *MovieModel) GetAll(mq *MovieQuery) ([]*Movie, error) {
	movies := []*Movie{}

	s := ""
	c := strings.TrimPrefix(mq.Sort, "-")
	s += c
	if strings.HasPrefix(mq.Sort, "-") {
		s += " DESC"
	}
	if c != "id" {
		s += ", id DESC"
	}

	stmt := fmt.Sprintf("SELECT id, title, year, runtime, genres, created FROM movies WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (genres @> $2 OR $2 = '{}') ORDER BY %s LIMIT $3 OFFSET $4;", s)

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	rows, err := m.Db.Query(ctx, stmt, mq.Title, mq.Genres, mq.PageSize, (mq.Page - 1) * mq.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := Movie{}

		if err := rows.Scan(&m.Id, &m.Title, &m.Year, &m.Runtime, &m.Genres, &m.Created); err != nil {
			return nil, err
		}

		movies = append(movies, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *MovieModel) Set(mv *Movie) error {
	stmt := "INSERT INTO movies (title, year, runtime, genres) VALUES($1, $2, $3, $4) RETURNING id, created;"

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	if err := m.Db.QueryRow(ctx, stmt, mv.Title, mv.Year, mv.Runtime, mv.Genres).Scan(&mv.Id, &mv.Created); err != nil {
		return err
	}

	return nil
}

func (m *MovieModel) Get(id int) (*Movie, error) {
	stmt := "SELECT id, title, year, runtime, genres, created, version FROM movies WHERE id = $1;"

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	
	mv := &Movie{}
	if err := m.Db.QueryRow(ctx, stmt, id).Scan(&mv.Id, &mv.Title, &mv.Year, &mv.Runtime, &mv.Genres, &mv.Created, &mv.Version); err != nil {
		return nil, err
	}
	
	return mv, nil
}

func (m *MovieModel) Update(mv *Movie, id int) error {
	stmt := "UPDATE movies SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 WHERE id = $5 AND version = $6 RETURNING version;"

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.Db.QueryRow(ctx, stmt, mv.Title, mv.Year, mv.Runtime, mv.Genres, id, mv.Version).Scan(&mv.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m *MovieModel) Delete(id int) error {
	stmt := "DELETE FROM movies WHERE id = $1;"

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	result, err := m.Db.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return errors.New("row not found")
	}

	return nil
}