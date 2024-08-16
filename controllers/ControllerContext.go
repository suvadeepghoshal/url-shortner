package controllers

import "github.com/go-playground/validator/v10"

type ControllerContext struct {
	Validator *validator.Validate
}
