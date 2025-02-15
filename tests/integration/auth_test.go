package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func TestAuth(t *testing.T) {
	// Устанавливаем URL сервера
	serverURL := "http://localhost:8080/api/auth"

	tests := []struct {
		name         string
		authRequest  AuthRequest
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful Auth",
			authRequest: AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Auth - Invalid username",
			authRequest: AuthRequest{
				Username: "",
				Password: "password123",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"errors": "Invalid username"`,
		},
		{
			name: "Invalid Auth - Invalid password",
			authRequest: AuthRequest{
				Username: "testuser",
				Password: "",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"errors": "Invalid password"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Преобразуем запрос в JSON
			requestBody, err := json.Marshal(tt.authRequest)
			assert.NoError(t, err)

			// Создаем новый запрос
			req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Отправляем запрос
			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Проверяем код статуса ответа
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Читаем тело ответа
			responseBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Проверяем тело ответа
			assert.Contains(t, string(responseBody), tt.expectedBody)
		})
	}
}
