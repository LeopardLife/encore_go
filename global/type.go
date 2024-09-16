package global

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
}
