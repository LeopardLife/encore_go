package role

import (
	"context"

	"encore.app/database"
	"encore.app/global"
)

// RoleRequest is the request object for adding a role.
type RoleRequest struct {
	Name string `json:"name"`
}

// Role is the model for a role.
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//encore:api public method=GET path=/roles/all
func GetAllRoles(ctx context.Context) (*global.ApiResponse[[]Role], error) {
	const query = "SELECT id, name FROM roles"
	rows, err := database.Database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return &global.ApiResponse[[]Role]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    roles,
	}, nil
}

//encore:api public method=POST path=/roles/add
func AddRole(ctx context.Context, req *RoleRequest) (*global.ApiResponse[Role], error) {
	if CheckRoleExists(ctx, req.Name) {
		return &global.ApiResponse[Role]{
			Code:    400,
			Status:  400,
			Message: "Role already exists",
		}, nil
	}

	const query = "INSERT INTO roles (name) VALUES ($1) RETURNING id"
	row := database.Database.QueryRow(ctx, query, req.Name)
	var role Role
	if err := row.Scan(&role.ID); err != nil {
		return nil, err
	}

	return &global.ApiResponse[Role]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    role,
	}, nil
}

func CheckRoleExists(ctx context.Context, name string) bool {
	const query = "SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)"
	row := database.Database.QueryRow(ctx, query, name)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false
	}
	return exists
}

//encore:api public method=POST path=/roles/delete
func DeleteRole(ctx context.Context, req *RoleRequest) (*global.ApiResponse[Role], error) {
	const query = "DELETE FROM roles WHERE name = $1 RETURNING id"
	row := database.Database.QueryRow(ctx, query, req.Name)
	var role Role
	if err := row.Scan(&role.ID); err != nil {
		return nil, err
	}

	return &global.ApiResponse[Role]{
		Code:    200,
		Status:  200,
		Message: "Success",
		Data:    role,
	}, nil
}
