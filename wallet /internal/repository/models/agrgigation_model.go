package models

type Day struct {
	Date string `json:"date"`
}

type Week struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type Month struct {
	Date_from string `json:"dateFrom"`
	Date_to   string `json:"dateTo"`
}
