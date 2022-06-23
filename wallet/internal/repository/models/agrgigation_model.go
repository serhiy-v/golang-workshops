package models

type Day struct {
	Date string `json:"date"`
}

type Week struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type Month struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}
