package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"time"

	"encore.app/database"
	"encore.app/global"
	"encore.app/users"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type MyAuthParams struct {
	// SessionCookie is set to the value of the "session" cookie.
	// If the cookie is not set it's nil.
	SessionCookie *http.Cookie `cookie:"session"`

	// ClientID is the unique id of the client, sourced from the URL query string.
	ClientID string `query:"client_id"`

	// Authorization is the raw value of the "Authorization" header
	// without any parsing.
	Authorization string `header:"Authorization"`
}

type MyAuthResponse struct {
	Username string
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Xây dựng URL mới với path tùy chỉnh
	newURL := url.URL{
		Scheme: "https",
		Host:   r.Host,
		Path:   "/new-path",
	}

	// Trả về URL tùy chỉnh cho client
	fmt.Fprintf(w, "New URL: %s", newURL.String())
}

//encore:authhandler
func AuthHandler(ctx context.Context, p *MyAuthParams) (auth.UID, *global.DataAuth, error) {
	var id string
	if p.SessionCookie != nil {
		// Check the session cookie
	}

	if p.ClientID != "" {
		// Check the client ID
	}

	print("Authorization: ")

	if p.Authorization != "" {
		// Check the Authorization header
		read, err := VerifyJWT(p.Authorization)
		if err != nil {
			return "", &global.DataAuth{}, err
		}

		err = database.Database.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", read).Scan(&id)
		if err != nil {
			return "", &global.DataAuth{}, &errs.Error{
				Code:    400,
				Message: "user not found",
			}
		}

		return auth.UID(id), &global.DataAuth{Username: read}, nil
	}

	return "", &global.DataAuth{}, nil
}

type MyLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MyLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// login is a private function that checks the username and password
// and returns a user ID if the credentials are correct.
//
//encore:api public method=POST path=/auth/login
func login(ctx context.Context, p *MyLoginParams) (*global.ApiResponse[MyLoginResponse], error) {
	var hashedPassword string
	const query = "SELECT password FROM users WHERE username = $1"
	row := database.Database.QueryRow(ctx, query, p.Username)
	if err := row.Scan(&hashedPassword); err != nil {
		return nil, &errs.Error{
			Code:    400,
			Message: "user not found",
		}
	}

	if err := CheckPassword(hashedPassword, p.Password); err != nil {
		return nil, &errs.Error{
			Code:    400,
			Message: "invalid password",
		}
	}

	token, err := GenerateJWT(p.Username)
	if err != nil {
		return nil, err
	}

	refreshtoken, err := GenerateRefreshToken(p.Username)
	if err != nil {
		return nil, err
	}

	return &global.ApiResponse[MyLoginResponse]{
		Code: 200,
		Data: MyLoginResponse{
			Token:        token,
			RefreshToken: refreshtoken,
		},
	}, nil
}

// signUp is a private function that checks if the user already exists
// and returns an error if the user already exists.
//
//encore:api public method=POST path=/auth/signup
func signUp(ctx context.Context, p *MyLoginParams) error {
	if users.CheckUserExists(ctx, p.Username) {
		return &errs.Error{
			Code:    400,
			Message: "user already",
		}
	}

	hashedPassword, err := HashPassword(p.Password)
	if err != nil {
		return err
	}

	const query = "INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())"
	_, err = database.Database.Exec(ctx, query, p.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

var jwtKey = []byte("your-secret-key")

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return token.Claims.(jwt.MapClaims)["username"].(string), nil
}

func GenerateRefreshToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}
