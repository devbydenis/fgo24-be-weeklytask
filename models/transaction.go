package models

import (
	u "be-weeklytask-ewallet/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	SenderID        *uuid.UUID `json:"sender_id"`
	ReceiverID      *uuid.UUID `json:"receiver_id"`
	TransactionType string     `json:"transaction_type"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	Description     string     `json:"description"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

// TransactionHistory untuk response history
type TransactionHistory struct {
	Transaction
	Direction       string `json:"direction"` // IN atau OUT
	CounterpartName string `json:"counterpart_name"`
}

// Transfer - Transfer uang antar user
func Transfer(senderID, receiverID uuid.UUID, amount float64, description, notes string) error {
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
	// Start transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Cek saldo pengirim
	var senderBalance float64
	err = tx.QueryRow(context.Background(),
		`SELECT balance FROM users WHERE id = $1`,
		senderID).Scan(&senderBalance)
	if err != nil {
		return err
	}

	// Cek saldo cukup
	if senderBalance < amount {
		return pgx.ErrNoRows // Custom error untuk saldo tidak cukup
	}

	// Insert transaksi
	transactionID := uuid.New()
	completedAt := time.Now()

	_, err = tx.Exec(context.Background(),
		`INSERT INTO transactions (id, sender_id, receiver_id, transaction_type, amount, status, description, notes, completed_at)
				VALUES ($1, $2, $3, 'TRANSFER', $4, 'COMPLETED', $5, $6, $7)`,
		transactionID, senderID, receiverID, amount, description, notes, completedAt)
	if err != nil {
		return err
	}

	// Update balance pengirim
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET balance = balance - $1 WHERE id = $2`,
		amount, senderID)
	if err != nil {
		return err
	}

	// Update balance penerima
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET balance = balance + $1 WHERE id = $2`,
		amount, receiverID)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(context.Background())
}

// TopUp - Top up saldo user
func TopUp(userID uuid.UUID, amount float64, description string) error {
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
	
	// Start transaction

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert transaksi
	transactionID := uuid.New()
	completedAt := time.Now()

	_, err = tx.Exec(context.Background(),
		`INSERT INTO transactions (id, receiver_id, transaction_type, amount, status, description, completed_at)
				VALUES ($1, $2, 'TOP_UP', $3, 'COMPLETED', $4, $5)`,
		transactionID, userID, amount, description, completedAt)
	if err != nil {
		return err
	}

	// Update balance
	_, err = tx.Exec(context.Background(),
		`UPDATE users SET balance = balance + $1 WHERE id = $2`,
		amount, userID)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(context.Background())
}

// GetHistory - Get transaction history
func GetHistory(userID uuid.UUID, limit int) ([]TransactionHistory, error) {
		// conncect to db
	conn, err := u.DBConnect()
	if err != nil {
		fmt.Println("IsEmailExist error connet to db:", err)
		return nil, err
	}

	// jangan lupa tutup kalo udah selesai
	defer func() {
		conn.Conn().Close(context.Background())
	}()

	rows, err := conn.Query(context.Background(),
		`SELECT 
            t.id, t.sender_id, t.receiver_id, t.transaction_type, t.amount, 
            t.status, t.description, t.notes, t.created_at, t.completed_at,
            CASE 
                WHEN t.sender_id = $1 THEN 'OUT'
                WHEN t.receiver_id = $1 THEN 'IN'
            END as direction,
            CASE 
                WHEN t.sender_id = $1 THEN pr.full_name
                WHEN t.receiver_id = $1 THEN ps.full_name
                ELSE ''
            END as counterpart_name
				FROM transactions t
				LEFT JOIN profiles ps ON t.sender_id = ps.user_id
				LEFT JOIN profiles pr ON t.receiver_id = pr.user_id
				WHERE t.sender_id = $1 OR t.receiver_id = $1
				ORDER BY t.created_at DESC
				LIMIT $2`,
		userID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []TransactionHistory
	for rows.Next() {
		var h TransactionHistory
		err = rows.Scan(
			&h.ID, &h.SenderID, &h.ReceiverID, &h.TransactionType, &h.Amount,
			&h.Status, &h.Description, &h.Notes, &h.CreatedAt, &h.CompletedAt,
			&h.Direction, &h.CounterpartName)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}
