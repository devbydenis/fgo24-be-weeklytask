package models

import (
	u "be-weeklytask-ewallet/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	Id        string  		`json:"id"`
	Email     string  		`json:"email"`
	Password  string  		`json:"password"`
	Pin       string  		`json:"pin"`
	Balance   float32 		`json:"balance"`
	IsActive  bool    		`json:"is_active"`
	CreatedAt time.Time  	`json:"created_at"`
	UpdatedAt time.Time  	`json:"updated_at"`
}

type Profile struct {
	Id           string 	 `json:"id"`
	UserId       string 	 `json:"user_id"`
	FullName     *string 	 `json:"full_name"`
	Phone        *string 	 `json:"phone"`
	ProfileImage *string 	 `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type IsEmailExistType struct {
	Email string `json:"email"`
}

type GetProfileFromDbType struct {
	FullName     string  `json:"full_name,omitempty"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	ProfileImage string  `json:"profile_image"`
	Balance      float32 `json:"balance"`
}

type UserWithProfile struct {
	User
	Profile Profile `json:"profile"`
}

func IsEmailExist(email string) bool {
	// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return false
	}
	// jangan lupa tutup kalo udah selesai
	defer func() {
		conn.Conn().Close(context.Background())
	}()

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


// func GetProfileFromDb(user_id string) (GetProfileFromDbType, error) {
// 	// conncect to db
// 	conn, err := u.DBConnect()
// 	if err != nil {
// 		fmt.Println("IsEmailExist error connet to db:", err)
// 		return GetProfileFromDbType{}, err
// 	}

// 	// jangan lupa tutup kalo udah selesai
// 	defer func() {
// 		conn.Conn().Close(context.Background())
// 	}()

// 	// check if email exist
// 	rows, err := conn.Query(
// 		context.Background(),
// 		`
// 			SELECT 
// 				p.full_name,
// 				p.phone,
// 				p.profile_image,
// 				u.email,
// 				u.balance
// 			FROM profiles p
// 			JOIN users u ON p.user_id = u.id
// 			WHERE u.id = $1
// 		`,
// 		user_id)
// 	if err != nil {
// 		fmt.Println("GetProfileFromDB error query:", err)
// 		return GetProfileFromDbType{}, err
// 	}

// 	// collect row and map to struxt
// 	profiles, err := pgx.CollectRows[GetProfileFromDbType](rows, pgx.RowToStructByName)
// 	if err != nil {
// 		fmt.Println("GetProfileFromDB error collect row:", err)
// 		return GetProfileFromDbType{}, err
// 	}

// 	fmt.Println("GetProfileFromDB users:", profiles)
// 	if len(profiles) > 0 {
// 		return profiles[0], nil
// 	}

// 	return GetProfileFromDbType{}, nil
// } 

func GetProfileFromDb(id uuid.UUID) (*UserWithProfile, error) {
	var user User
	var profile Profile

	// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return &UserWithProfile{}, err
	}

	// jangan lupa tutup kalo udah selesai
	defer func() {
		conn.Conn().Close(context.Background())
	}()

	err = conn.QueryRow(context.Background(),
		`
		SELECT u.id, u.email, u.balance, u.is_active, u.created_at, 
					p.id, p.full_name, p.phone, p.profile_image
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE u.id = $1`,
		id).Scan(
		&user.Id, &user.Email, &user.Balance, &user.IsActive, &user.CreatedAt,
		&profile.Id, &profile.FullName, &profile.Phone, &profile.ProfileImage)
	
	if err != nil {
		return nil, err
	}

	profile.UserId = user.Id

	return &UserWithProfile{
		User:    user,
		Profile: profile,
	}, nil
}


// UpdateBalance - Update balance user
func UpdateBalance(userID uuid.UUID, amount float64) error {
		// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return err
	}

	// jangan lupa tutup kalo udah selesai
	defer func() {
		conn.Conn().Close(context.Background())
	}()

	_, err = conn.Exec(context.Background(),
		`UPDATE users SET balance = balance + $1 WHERE id = $2`,
		amount, userID)

	return err
}
