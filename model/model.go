package model

import "time"

type Currency struct {
	ID    int
	Title string
	Code  string
	Value float64
	ADate time.Time
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
