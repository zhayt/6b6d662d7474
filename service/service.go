package service

import (
	"encoding/xml"
	"errors"
)

type Rates struct {
	XMLName xml.Name
	Items   []Item `xml:"item"`
}

type Item struct {
	FullName    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
}

var ErrUserStupid = errors.New("user error")
