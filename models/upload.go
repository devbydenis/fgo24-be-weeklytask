package models

import (
	"be-weeklytask-ewallet/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func ChangeProfileImgDB(filename string, userId uuid.UUID) error {
		// conncect to db
	conn, err := utils.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return err
	}

	// jangan lupa tutup kalo udah selesai
	defer func() {
		conn.Conn().Close(context.Background())
	}()

	_, err = conn.Exec(context.Background(),
		`UPDATE profiles SET profile_image = $1 WHERE user_id = $2`,
		filename, userId)

	return err
}