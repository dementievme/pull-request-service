package application

type ErrorDTO struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}
