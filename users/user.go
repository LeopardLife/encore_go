package users

import (
	"context"
	"log"

	"encore.app/global"
	"encore.dev/storage/sqldb"
)

// UserRequest is the request object for adding a user.
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User is the model for a user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

var userdb = sqldb.NewDatabase("users", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:api public method=GET path=/users/all
func GetAllUsers(ctx context.Context) (*global.ApiResponse[[]User], error) {
	const query = "SELECT id, username FROM users"
	rows, err := userdb.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &global.ApiResponse[[]User]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    users,
	}, nil
}

//encore:api public method=POST path=/users/add
func AddUser(ctx context.Context, req *UserRequest) (*global.ApiResponse[User], error) {
	if CheckUserExists(ctx, req.Username) {
		return &global.ApiResponse[User]{
			Code:    400,
			Status:  400,
			Message: "User already exists",
		}, nil
	}

	const query = "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := userdb.Exec(ctx, query, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	var user User
	err = userdb.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", req.Username).Scan(&user.ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &global.ApiResponse[User]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    User{Username: req.Username, ID: user.ID},
	}, nil
}

// CheckUserExists checks if a user with the given username already exists.
func CheckUserExists(ctx context.Context, username string) (exists bool) {
	const query = "SELECT id FROM users WHERE username = $1"
	var id string
	err := userdb.QueryRow(ctx, query, username).Scan(&id)
	if err == nil {
		return true
	}

	return false
}

//encore:api public method=GET path=/users/get/:id
func GetUser(ctx context.Context, id string) (*global.ApiResponse[User], error) {
	const query = "SELECT id, username FROM users WHERE id = $1"
	var user User
	err := userdb.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return &global.ApiResponse[User]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    user,
	}, nil
}

//encore:api public method=DELETE path=/users/delete/:id
func DeleteUser(ctx context.Context, id string) (*global.ApiResponse[User], error) {
	const query = "DELETE FROM users WHERE id = $1"
	_, err := userdb.Exec(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &global.ApiResponse[User]{
		Code:    200,
		Status:  200,
		Message: "Success",
	}, nil
}

//encore:api public method=PUT path=/users/update/:id
func UpdateUser(ctx context.Context, id string, req *UserRequest) (*global.ApiResponse[User], error) {
	const query = "UPDATE users SET username = $1, password = $2 WHERE id = $3"
	_, err := userdb.Exec(ctx, query, req.Username, req.Password, id)
	if err != nil {
		return nil, err
	}

	return &global.ApiResponse[User]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    User{Username: req.Username, ID: id},
	}, nil
}
