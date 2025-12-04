package forum

import "testing"

func TestAuth(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		want     error
	}{
		{
			name:     "tester",
			email:    "newuser@tester.er",
			password: "testASDF12!@",
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
