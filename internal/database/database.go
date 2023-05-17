package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/EgorKo25/GophKeeper/internal/storage"
)

var (
	ErrRace = errors.New("the resource is busy")
)

// ManagerDB structure for managing database
type ManagerDB struct {
	Db *sqlx.DB
}

// NewManagerDB constructor
func NewManagerDB(address string) (*ManagerDB, error) {

	ctx := context.Background()

	db, err := sqlx.Open("pgx", address)
	if err != nil {
		return nil, err
	}

	err = createTable(ctx, db)
	if err != nil {
		return nil, err
	}

	return &ManagerDB{
		Db: db,
	}, nil
}

// createTable creates tables for the database operation
func createTable(ctx context.Context, db *sqlx.DB) error {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS 
    users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status BOOLEAN);`,

		`CREATE TABLE IF NOT EXISTS
	passwords (
	id SERIAL PRIMARY KEY, 
	service VARCHAR(255),
	login_owner VARCHAR(255) NOT NULL,
	login VARCHAR(255) NOT NULL,
	password VARCHAR(64) NOT NULL);`,

		`CREATE TABLE IF NOT EXISTS
	cards (
	id SERIAL PRIMARY KEY, 
    bank VARCHAR(255),
	login_owner VARCHAR(255) NOT NULL,
	number INT NOT NULL,
	date_end INT NOT NULL,
	secret_code INT NOT NULL,
	owner VARCHAR(50) NOT NULL);`,

		`CREATE TABLE IF NOT EXISTS
	binary_data (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255), 
	login_owner VARCHAR(255) NOT NULL,
	data bytea NOT NULL);`,
	}

	for _, query := range queries {
		_, err := db.ExecContext(childCtx, query)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping testing connection to database
func (m *ManagerDB) Ping() bool {
	err := m.Db.Ping()
	return err == nil
}

// CheckUser check user from database
func (m *ManagerDB) CheckUser(ctx context.Context, user *storage.User) (bool, error) {

	var check storage.User

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	check.Login = user.Login

	err := m.readUser(childCtx, &check)
	if err != nil {
		return false, err
	}

	return check.Password == user.Password, nil

}

// Add adds new data to database
func (m *ManagerDB) Add(ctx context.Context, src any, login string) error {
	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch data := src.(type) {
	case *storage.User:
		return m.addUser(childCtx, data)
	case storage.Password:
		data.LoginOwner = login
		return m.addPassword(childCtx, &data)
	case storage.BinaryData:
		data.LoginOwner = login
		return m.addBinData(childCtx, &data)
	case storage.Card:
		data.LoginOwner = login
		return m.addCard(childCtx, &data)
	default:

	}

	return errors.New("unknown adding type")
}

// addPassword adds new password
func (m *ManagerDB) addPassword(ctx context.Context, password *storage.Password) error {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `INSERT INTO passwords (service, login_owner, login, password)
							VALUES  (:service, :login_owner, :login, :password);`

	_, err := m.Db.NamedExecContext(childCtx, query, password)
	if err != nil {
		return err
	}

	return nil
}

// addCard adds new card
func (m *ManagerDB) addCard(childCtx context.Context, card *storage.Card) error {

	query := `INSERT INTO cards (bank, login, number, dataEnd, secret_code, owner)
							VALUES  (:bank, :login_owner, :number, :data_end, :secret_code, :owner);`

	_, err := m.Db.NamedExecContext(childCtx, query, card)
	if err != nil {
		return err
	}

	return nil
}

// addBinData adds new binary data
func (m *ManagerDB) addBinData(childCtx context.Context, data *storage.BinaryData) error {

	query := `INSERT INTO binary_data (title, login_owner, password)
							VALUES  (:title, :login_owner, :data);`

	_, err := m.Db.NamedExecContext(childCtx, query, data)
	if err != nil {
		return err
	}

	return nil
}

// addUser adds new user to the database
func (m *ManagerDB) addUser(childCtx context.Context, user *storage.User) error {

	query := `INSERT INTO users (username, password, email, created_at, updated_at) 
							VALUES  (:username, :password, :email, :created_at, :updated_at)`
	_, err := m.Db.NamedExecContext(childCtx, query, user)
	if err != nil {
		return err
	}

	return nil
}

// Read reads data from database
func (m *ManagerDB) Read(ctx context.Context, src any, login string) ([]byte, error) {
	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch data := src.(type) {
	case *storage.User:
		err := m.readUser(childCtx, data)
		if err != nil {
			return []byte(""), err
		}

		if data.Status {
			return []byte(""), err
		}

		res, err := json.Marshal(&data)
		if err != nil {
			return []byte(""), err
		}
		return res, nil
	case storage.Password:
		data.LoginOwner = login
		err := m.readPassword(childCtx, &data)
		if err != nil {
			return []byte(""), err
		}
		res, err := json.Marshal(&data)
		if err != nil {
			return []byte(""), err
		}
		return res, nil
	case storage.BinaryData:
		data.LoginOwner = login
		err := m.readBinary(childCtx, &data)
		if err != nil {
			return []byte(""), err
		}
		res, err := json.Marshal(&data)
		if err != nil {
			return []byte(""), err
		}
		return res, nil
	case storage.Card:
		data.LoginOwner = login
		err := m.readCard(childCtx, &data)
		if err != nil {
			return []byte(""), err
		}
		res, err := json.Marshal(&data)
		if err != nil {
			return []byte(""), err
		}
		return res, nil
	}

	return nil, errors.New("unknown type: " + fmt.Sprintf("%T", src))
}

// readPassword read password
func (m *ManagerDB) readPassword(childCtx context.Context, password *storage.Password) error {

	query := `SELECT * FROM passwords  WHERE service = :service AND login_owner = :login_owner;`

	rows, err := m.Db.NamedQueryContext(childCtx, query, password)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.StructScan(password)

	if err != nil {
		return err
	}

	return nil
}

// readCard read card data
func (m *ManagerDB) readCard(childCtx context.Context, card *storage.Card) error {

	query := `SELECT * FROM cards  WHERE login_owner = :login_owner AND bank = :bank;`

	rows, err := m.Db.NamedQueryContext(childCtx, query, card)
	if err != nil {
		return err
	}

	err = rows.StructScan(card)
	if err != nil {
		return err
	}

	return nil
}

// readBinary read binary data
func (m *ManagerDB) readBinary(childCtx context.Context, binary *storage.BinaryData) error {

	query := `SELECT * FROM binary_data WHERE title = :title AND login_owner = :login_owner;`

	rows, err := m.Db.NamedQueryContext(childCtx, query, binary)
	if err != nil {
		return err
	}

	err = rows.StructScan(binary)
	if err != nil {
		return err
	}

	return nil
}

// readUser read user data
func (m *ManagerDB) readUser(childCtx context.Context, user *storage.User) error {

	query := `SELECT * FROM users WHERE username = :username;`

	rows, err := m.Db.NamedQueryContext(childCtx, query, user)
	if err != nil {
		return err
	}

	rows.Next()
	_ = rows.StructScan(user)

	return nil

}

// Update updates data from database
func (m *ManagerDB) Update(ctx context.Context, src any, login string) error {
	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch data := src.(type) {
	case storage.User:
		return m.updateUser(childCtx, &data)
	case storage.Password:
		data.LoginOwner = login
		return m.updatePassword(childCtx, &data)
	case storage.BinaryData:
		data.LoginOwner = login
		return m.updateBinData(childCtx, &data)
	case storage.Card:
		data.LoginOwner = login
		return m.updateCard(childCtx, &data)
	default:
	}

	return errors.New("unknown updating type")
}

// updatePassword update user password
func (m *ManagerDB) updatePassword(childCtx context.Context, password *storage.Password) error {

	query := `UPDATE passwords SET service = :service, login = :login, password = :password 
                 WHERE login_owner = :login_owner AND service = :service;`

	_, err := m.Db.NamedExecContext(childCtx, query, password)
	if err != nil {
		return err
	}

	return nil
}

// updateCard update user card
func (m *ManagerDB) updateCard(childCtx context.Context, card *storage.Card) error {

	query := `UPDATE cards SET bank = :bank, number = :number, data_end = :data_end,
                 secret_code = :secret_code, owner = :owner
                 WHERE login_owner = :login_owner AND bank = :bank;`

	_, err := m.Db.NamedExecContext(childCtx, query, card)
	if err != nil {
		return err
	}

	return nil
}

// updateBinDAta update user binary data
func (m *ManagerDB) updateBinData(childCtx context.Context, binary *storage.BinaryData) error {

	query := `UPDATE binary_data SET data = :data WHERE login_owner = :login_owner AND title = :title;`

	_, err := m.Db.NamedExecContext(childCtx, query, binary)
	if err != nil {
		return err
	}

	return nil
}

// updateUser update user
func (m *ManagerDB) updateUser(childCtx context.Context, user *storage.User) error {

	query := `UPDATE users SET username = :username, status = :status
                 WHERE username = :username;`

	_, err := m.Db.NamedExecContext(childCtx, query, user)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes data from database
func (m *ManagerDB) Delete(ctx context.Context, src any, login string) error {
	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch data := src.(type) {
	case storage.User:
		return m.deleteUser(childCtx, &data)
	case storage.Password:
		data.LoginOwner = login
		return m.deletePassword(childCtx, &data)
	case storage.BinaryData:
		data.LoginOwner = login
		return m.deleteBinData(childCtx, &data)
	case storage.Card:
		data.LoginOwner = login
		return m.deleteCard(childCtx, &data)
	default:
	}

	return errors.New("unknown updating type")
}

// deletePassword delete pair login:password:service from db
func (m *ManagerDB) deletePassword(childCtx context.Context, password *storage.Password) error {

	query := `DELETE FROM passwords WHERE login_owner = :login_owner AND service = :service;`

	_, err := m.Db.NamedExecContext(childCtx, query, password)
	if err != nil {
		return err
	}

	return nil
}

// deleteCard delete user card
func (m *ManagerDB) deleteCard(childCtx context.Context, card *storage.Card) error {

	query := `DELETE FROM cards WHERE login_owner = :login_owner AND bank = :bank;`

	_, err := m.Db.NamedExecContext(childCtx, query, card)
	if err != nil {
		return err
	}
	return nil
}

// deleteBinData delete user card
func (m *ManagerDB) deleteBinData(childCtx context.Context, data *storage.BinaryData) error {

	query := `DELETE FROM binary_data WHERE login_owner = :login_owner AND title = :title;`

	_, err := m.Db.NamedExecContext(childCtx, query, data)
	if err != nil {
		return err
	}
	return nil
}

// deleteUser delete user profile
func (m *ManagerDB) deleteUser(childCtx context.Context, user *storage.User) error {

	query := `DELETE FROM users WHERE login = :login;`

	_, err := m.Db.NamedExecContext(childCtx, query, user)
	if err != nil {
		return err
	}
	return nil
}
