package permission

import (
	"context"
	"log"

	"encore.app/database"
	"encore.app/global"
)

// PermissionRequest is the request object for adding a permission.
type PermissionRequest struct {
	Name string `json:"name"`
}

// Permission is the model for a permission.
type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//encore:api public method=GET path=/permissions/GetAllRoles
func GetAllPermissions(ctx context.Context) (*global.ApiResponse[[]Permission], error) {
	const query = "SELECT id, name FROM permissions"
	rows, err := database.Database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		if err := rows.Scan(&permission.ID, &permission.Name); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return &global.ApiResponse[[]Permission]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    permissions,
	}, nil
}

//encore:api public method=POST path=/permissions/add
func AddPermission(ctx context.Context, req *PermissionRequest) (*global.ApiResponse[Permission], error) {
	if CheckPermissionExists(ctx, req.Name) {
		return &global.ApiResponse[Permission]{
			Code:    400,
			Status:  400,
			Message: "Permission already exists",
		}, nil
	}

	const query = "INSERT INTO permissions (name) VALUES ($1) RETURNING id"
	row := database.Database.QueryRow(ctx, query, req.Name)
	var id string
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &global.ApiResponse[Permission]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    Permission{ID: id, Name: req.Name},
	}, nil
}

// CheckPermissionExists checks if a permission already exists.
func CheckPermissionExists(ctx context.Context, name string) bool {
	const query = "SELECT id FROM permissions WHERE name = $1"
	row := database.Database.QueryRow(ctx, query, name)
	var id string
	if err := row.Scan(&id); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// GetPermissionID returns the ID of a permission.
func GetPermissionID(ctx context.Context, name string) (string, error) {
	const query = "SELECT id FROM permissions WHERE name = $1"
	row := database.Database.QueryRow(ctx, query, name)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

// GetPermissionName returns the name of a permission.
func GetPermissionName(ctx context.Context, id string) (string, error) {
	const query = "SELECT name FROM permissions WHERE id = $1"
	row := database.Database.QueryRow(ctx, query, id)
	var name string
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}
