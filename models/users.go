package models

import (
	u "be-weeklytask-ewallet/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Pin		    string  `json:"pin_hash"`
	Balance   float32 `json:"balance"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type Profile struct {
	Id           string `json:"id"`
	UserId       string `json:"user_id"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	ProfileImage string `json:"profile_image"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type IsEmailExistType struct {
	Email string `json:"email"`
}

func IsEmailExist(email string) bool {
	// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return false
	}
	// check if email exist
	rows, err := conn.Query(context.Background(), "SELECT email FROM users WHERE email = $1", email)
	if err != nil {
		fmt.Println("IsEmailExist error query:", err)
		return false
	}
	// collect row and map to struxt
	users, err := pgx.CollectRows[IsEmailExistType](rows, pgx.RowToStructByName)
	if err != nil {
		fmt.Println("IsEmailExist error collect row:", err)
		return false
	}

	fmt.Println("IsEmailExist users:", users)
	if len(users) > 0 {
		return true
	}

	return false
}
