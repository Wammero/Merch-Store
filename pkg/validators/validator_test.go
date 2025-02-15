package validators

import "testing"

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"Empty password", "", false},
		{"Non-empty password", "securepassword", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.password); got != tt.want {
				t.Errorf("IsValidPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"Empty username", "", false},
		{"Non-empty username", "john_doe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUsername(tt.username); got != tt.want {
				t.Errorf("IsValidUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}
