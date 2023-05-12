package service

import "encoding/xml"

type Rates struct {
	XMLName xml.Name
	Items   []Item `xml:"item"`
}

type Item struct {
	FullName    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
}
