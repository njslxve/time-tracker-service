package entity

import "time"

type User struct {
	ID         string
	UserID     int
	Passport   string
	Name       string
	Surmame    string
	Patronymic string
	Adress     string
}

type Task struct {
	ID        string
	TaskID    string
	UserID    int
	StartTime time.Time
	EndTime   time.Time
	Duration  int
}

type FilterOptions struct {
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	Adress     string `json:"adress,omitempty"`
}

type PaginationOptions struct {
	Limit int
	Next  string
}

type TokenData struct {
	Params  []byte
	Old     int
	IsAlive bool
}
