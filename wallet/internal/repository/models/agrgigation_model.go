package models

type Day struct {
	Date    string `bson:"_id"`
	Income  int
	Outcome int
}

type Week struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type Month struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}
