package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"sync"
)

var once sync.Once
var validate *validator.Validate

// Validate returns the Validate instance
func Validate() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}
