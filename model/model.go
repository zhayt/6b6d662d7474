package model

import "time"

type Currency struct {
	ID    int       `json:"-" db:"ID"`
	Title string    `json:"title" db:"TITLE"`
	Code  string    `json:"code" db:"CODE"`
	Value float64   `json:"value" db:"VALUE"`
	ADate time.Time `json:"ADate" db:"A_DATE"`
}
