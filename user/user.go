package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"
)

type User struct {
	ID       string
	email    string
	password string
}

type CreateUserParams struct {
	email    string
	password string
}

// CreateUser creates a new user.
//
//encore:api public method=POST path=/user
func CreateUser(ctx context.Context, p *CreateUserParams) (*User, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.email, p.password); err != nil {
		return nil, err
	}
	return &User{ID: id, email: p.email, password: p.password}, nil
}

type GetUserParams struct {
	email    string
	password string
}

// GetUser retrieves the user for the email.
//
//encore:api public method=GET path=/user/:email
func GetUser(ctx context.Context, email string) (*User, error) {
	u := &User{email: email}
	err := db.QueryRow(ctx, `
		SELECT password FROM user
		WHERE email = $1
	`, email).Scan(&u.password)
	return u, err
}

type ListResponse struct {
	Users []*User
}

// List retrieves all users.
//
//encore:api public method=GET path=/user
func List(ctx context.Context) (*ListResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT email, password FROM user
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.email, &u.password); err != nil {
			return nil, err
		}
	}
	return &ListResponse{Users: users}, nil
}

func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func insert(ctx context.Context, id, email, password string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO user (id, email, password)
		VALUES ($1, $2, $3)
	`, id, email, password)
	return err
}

var db = sqldb.NewDatabase("user", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
