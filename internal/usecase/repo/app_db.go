package repo

import (
	"app/internal/entity"
	"app/pkg/logger"
	"app/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func GetUserBalance(id int) (*entity.User, error) {

	logger.Info("Trying to get balance")

	user, err := getUserById(id)
	if err != nil {
		return nil, err
	}

	logger.Info("Success of geting Balance")

	return user, nil
}

func GetUserTransanctions(id int) ([]*entity.Transanction, error) {

	logger.Info("Trying to get history")

	query := fmt.Sprintf(`SELECT 
			user_id,		
			amount, 
			operation,
			date
		FROM transanctions
		WHERE user_id=$1`)
	rows, err := postgres.GetConn().Instance.Queryx(query, id)
	if err != nil {
		return nil, err
	}

	trans := []*entity.Transanction{}

	for rows.Next() {
		tr := entity.Transanction{}

		if err := rows.StructScan(&tr); err != nil {
			return nil, err
		}

		trans = append(trans, &tr)
	}

	logger.Info("Success of geting history")

	return trans, nil
}

func ChangeBalance(id int, amount float64) error {

	var operation string

	user, err := getUserById(id)

	if err == sql.ErrNoRows {
		user, err = insertNewUser(id)
	}

	if err != nil {
		return err
	}

	if amount > 0 {
		operation = "accrual"
	} else {
		operation = "debit"
	}
	tx, err := postgres.GetConn().Instance.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = insertTransanction(tx, user, amount, operation); err != nil {
		return err
	}

	amount = amount + user.Balance

	if err = updateBalance(tx, user.Id, amount); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func TransanctionBetweenUsers(id1, id2 int, amount float64) error {
	operation := "transfer"

	user1, err := getUserById(id1)
	if err != nil {
		return err
	}
	user2, err := getUserById(id2)
	if err != nil {
		return err
	}
	tx, err := postgres.GetConn().Instance.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	amount1 := user1.Balance - amount
	err = updateBalance(tx, user1.Id, amount1)
	if err != nil {
		return err
	}
	amount2 := user2.Balance + amount
	err = updateBalance(tx, user2.Id, amount2)
	if err != nil {
		return err
	}
	err = insertTransanction(tx, user1, -amount, operation)
	if err != nil {
		return err
	}
	err = insertTransanction(tx, user2, amount, operation)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getUserById(id int) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(`SELECT 
			id, 
			balance
		FROM users
		WHERE id=$1`)
	if err := postgres.GetConn().Instance.Get(&user, query, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func insertNewUser(id int) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(`INSERT INTO %s (
		id,
		balance
	  )
	  VALUES ($1, $2)`, "users")
	if _, err := postgres.GetConn().Instance.Exec(query, id, 0.00); err != nil {
		return nil, err
	}
	user.Id = id
	return &user, nil
}

func updateBalance(tx *sql.Tx, id int, amount float64) error {
	query := fmt.Sprintf(`UPDATE users 
	SET balance = $1
	WHERE id = $2;`)
	if _, err := tx.Exec(query, amount, id); err != nil {
		return err
	}
	return nil
}

func insertTransanction(tx *sql.Tx, user *entity.User, amount float64, operation string) error {
	query := fmt.Sprintf(`INSERT INTO %s (
		user_id,
		amount,
		operation,
		date
	  )
	  VALUES ($1, $2, $3, $4)`, "Transanctions")
	if _, err := tx.Exec(query, user.Id, amount, operation, time.Now()); err != nil {
		return err
	}

	return nil
}
