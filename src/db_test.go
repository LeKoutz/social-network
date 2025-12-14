package forum

import "testing"

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test",
			want: "mpla",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitDB("./data/db.db")
			if err != nil {
				t.Logf("%s", err.Error())
			}

		})
	}
}
