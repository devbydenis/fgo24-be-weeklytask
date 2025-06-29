package models

import (
	"be-weeklytask-ewallet/utils"
	"context"
)

type RegisterRequest struct {
	Email 					string `json:"email" form:"email"`
	Password 				string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	Pin							string `json:"pin" form:"pin"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func InsertUserToDB(email string, password string, pin string) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer func(){
		conn.Conn().Close(context.Background())
	}()

	_, err = conn.Exec(context.Background(), "INSERT INTO users (email, password, pin_hash) VALUES ($1, $2, $3)", email, password, pin)
	return err
}