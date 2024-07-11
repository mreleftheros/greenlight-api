package models

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNoRows = errors.New("no rows")