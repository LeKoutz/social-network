package utils

import (
	"strings"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "simple password",
			password: "testASDF12!@",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "pässwörd123",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(hash) == 0 {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plaintext password")
				}
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	hash1, _ := HashPassword("samepassword")
	hash2, _ := HashPassword("samepassword")
	if hash1 == hash2 {
		t.Error("HashPassword() should produce different hashes for same input")
	}
}

func TestStringToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{
			name:    "valid number",
			input:   "12345",
			want:    12345,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "large number",
			input:   "9999999999",
			want:    9999999999,
			wantErr: false,
		},
		{
			name:    "negative number",
			input:   "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "letters",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "float",
			input:   "12.34",
			want:    0,
			wantErr: true,
		},
		{
			name:    "with spaces",
			input:   " 123 ",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("StringToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringToTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid timestamp",
			input:   "1700000000",
			wantErr: false,
		},
		{
			name:    "invalid string",
			input:   "not-a-timestamp",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Year() < 2023 {
					t.Errorf("ConvertStringToTime() returned unexpected year %v", got.Year())
				}
			}
		})
	}
}

func TestConvertInt64ToTime(t *testing.T) {
	ts := int64(1700000000)
	got := ConvertInt64ToTime(ts)
	want := time.Unix(ts, 0)
	if !got.Equal(want) {
		t.Errorf("ConvertInt64ToTime() = %v, want %v", got, want)
	}
}

func TestConvertTimeToString(t *testing.T) {
	ts := int64(1700000000)
	tm := time.Unix(ts, 0)
	got := ConvertTimeToString(tm)
	if len(got) == 0 {
		t.Error("ConvertTimeToString() returned empty string")
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	ts := GetCurrentTimestamp()
	if len(ts) == 0 {
		t.Error("GetCurrentTimestamp() returned empty string")
	}
	_, err := StringToInt64(ts)
	if err != nil {
		t.Errorf("GetCurrentTimestamp() returned non-numeric string: %s", ts)
	}
}

func TestGetFunctionName(t *testing.T) {
	err := GetFunctionName()
	if err == nil {
		t.Error("GetFunctionName() returned nil")
	}
	if len(err.Error()) == 0 {
		t.Error("GetFunctionName() returned empty error message")
	}
}
