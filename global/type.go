package global

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
}

type DataAuth struct {
	Username    string       `json:"username"`
	Permissions []Permission `json:"permissions,omitempty"`
	Roles       []Role       `json:"roles,omitempty"`
	Profile     Profile      `json:"profile,omitempty"`
}

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Method      string `json:"method"`
	Endpoint    string `json:"endpoint"`
	Description string `json:"description,omitempty"`
	ModuleID    string `json:"module_id,omitempty"`
}

type Module struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id,omitempty"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Profile struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	UserID      string `json:"user_id"`
}
