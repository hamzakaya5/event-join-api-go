package models

type LoginResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Token   string      `json:"token"`
	Id      int         `json:"id"`
	Level   int         `json:"level"`
	Data    interface{} `json:"data,omitempty"`
}

type RegisterResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Level   int         `json:"level"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Join struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
