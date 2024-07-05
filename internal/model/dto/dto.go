package dto

import "time"

type UserInfoRequest struct {
	PassportSerie  string
	PassportNumber string
}

type UserInfoResponse struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Adress     string `json:"adress"`
}

type AddUserRequest struct {
	Passport string `json:"passportNumber"`
}

type UpdateUserRequest struct {
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	Passport   string `json:"passport,omitempty"`
	Adress     string `json:"adress,omitempty"`
}

type TaskRequest struct {
	UserID int    `json:"user_id"`
	TaskID string `json:"task_id"`
}

type TaskResponse struct {
	TaskID    string    `json:"task_id"`
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`
}

type Error struct {
	Message string `json:"error"`
}

type User struct {
	UserID     int    `json:"user_id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Passport   string `json:"passport"`
	Adress     string `json:"adress"`
}

type GetUsersResponse struct {
	Users []User `json:"users"`
	Next  string `json:"next_page,omitempty"`
}
