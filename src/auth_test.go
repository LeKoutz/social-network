package forum

import "testing"

func TestSaltGenerator(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "general salt",
			want: "nothing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := SaltGenerator(5)
			if err != nil {
				t.Error(err.Error())
			}
			t.Logf("%#v", a)
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
