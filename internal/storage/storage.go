package storage

import (
	"time"
)

// User structure describing the user
type User struct {
	Id        int       `db:"id"`
	Login     string    `json:"login" form:"login" db:"username"`
	Password  string    `json:"passwd" form:"password" db:"password"`
	Email     string    `json:"email" form:"email" db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Status    bool      `db:"status"`
}

type Card struct {
	Id         int    `db:"id" json:"id"`
	Bank       string `db:"bank" json:"bank"`
	LoginOwner string `db:"login_owner" json:"login_owner"`
	Number     string `db:"number" json:"number"`
	DataEnd    string `db:"data_end" json:"data_end"`
	SecretCode string `db:"secret_code" json:"secret_code"`
	Owner      string `db:"owner" json:"owner"`
}

type Password struct {
	Id         int    `db:"id" json:"id"`
	Service    string `db:"service" json:"service"`
	LoginOwner string `db:"login_owner" json:"login_owner"`
	Login      string `db:"login" json:"login"`
	Password   string `db:"password" json:"password"`
}

type BinaryData struct {
	Id         int    `db:"id" json:"id"`
	Title      string `db:"title" json:"title"`
	LoginOwner string `db:"login_owner" json:"login_owner"`
	Data       []byte `db:"data" json:"data"`
}

type UserDate struct {
	User       `json:"user"`
	Cards      []Card       `json:"cards"`
	Passwords  []Password   `json:"passwords"`
	BinaryData []BinaryData `json:"binary_data"`
}
