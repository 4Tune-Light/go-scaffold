package greeter

import "errors"

var (
	ErrNameRequired = errors.New("name is required")
	ErrNotFound     = errors.New("greeting not found")
)
