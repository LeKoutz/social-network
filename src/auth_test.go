package forum

import "testing"

func TestSaltGenerator(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{
			name:   "general salt",
			length: 5,
			want:   5,
		},
		{
			name:   "limit",
			length: 26,
			want:   26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt := SaltGenerator(tt.length)
			got := len(salt)
			if got != tt.want {
				t.Errorf("Expected: %d, got %d", tt.want, got)
			}
		})
	}
}

func TestAuth(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		want     error
	}{
		{
			name:     "placeholder",
			email:    "test@test.test",
			password: "test",
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Auth(tt.email, tt.password)
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}
