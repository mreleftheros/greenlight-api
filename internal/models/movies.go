package models

import "time"

type Movie struct {
	Id      int       `json:"id"`
	Created time.Time `json:"created"`
	Title   string    `json:"title"`
	Year    int       `json:"year"`
	Runtime int       `json:"runtime"`
	Genres  []string  `json:"genres"`
}
