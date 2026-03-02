package models

import (
	"errors"
)

// define a custom error which can be used to check for when a specific record is not found in the database
// to avoid being tied to the sql package and its error types, we can define our own error and return that instead when a record is not found
var ErrNoRecord = errors.New("models: no matching record found")