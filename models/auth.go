package models

import (
	u "be-weeklytask-ewallet/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	Pin      string `json:"pin"`
}

func InsertUserToDB(email string, password string, pin string, userUUID uuid.UUID) error {
	conn, err := u.DBConnect()
	if err != nil {
		return err
	}
	defer func(){
		conn.Conn().Close(context.Background())
	}()
	
	fmt.Println("InsertUserToDB userUUID:", userUUID)

	_, err = conn.Exec(context.Background(), `
		INSERT INTO users (id, email, password, pin) VALUES ($1, $2, $3, $4);
	`, userUUID, email, password, pin,
)

	_, err = conn.Exec(context.Background(), `
		INSERT INTO profiles (user_id) VALUES ($1);
	`, userUUID,
)

	return err
}

func MatchUserInDatabase(email string, password string, pin string) bool {
	// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("MatchUserInDatabase error connet to db:", err)
		return false
	}

	// jangan lupa tutup kalo udah selesai
	defer func(){
		conn.Conn().Close(context.Background())
	}()

	// check if email exist
	rows, err := conn.Query(context.Background(), "SELECT email, password, pin FROM users WHERE email = $1 AND password = $2 AND pin = $3", email, password, pin)
	if err != nil {
		fmt.Println("MatchUserInDatabase error query:", err)
		return false
	}

	// collect row and map to struxt
	users, err := pgx.CollectRows[LoginRequest](rows, pgx.RowToStructByName)
	if err != nil {
		fmt.Println("MatchUserInDatabase error collect row:", err)
		return false
	}
	
	fmt.Println("MatchUserInDatabase users:", users)
	if len(users) == 0 {
		return false
	}

	return true
}