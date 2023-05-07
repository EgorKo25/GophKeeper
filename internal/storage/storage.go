package storage

import "time"

// User structure describing the user
type User struct {
	Login     string `json:"login" form:"login"`
	Passwd    string `json:"passwd" form:"password"`
	Email     string `json:"email" form:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
