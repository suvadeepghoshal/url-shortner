package controllers

import (
	"github.com/go-playground/validator/v10"
)

// ControllerContext Passing the validator instance to all the controllers, so that only one reference of the validator is created and can be used all the time
type ControllerContext struct {
	Validator *validator.Validate
}
