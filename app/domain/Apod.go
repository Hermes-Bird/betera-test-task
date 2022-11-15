package domain

import "time"

type Apod struct {
	Id   int
	Url  string
	Date time.Time
}

func NewApod(url string) Apod {
	return Apod{
		Url:  url,
		Date: time.Now(),
	}
}
