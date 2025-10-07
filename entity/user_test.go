package entity

import (
	"testing"
	"time"
)

func TestUser_HashPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "verylongpasswordthatshouldworkfine123456789",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := &User{
				Password: tt.password,
			}

			err := user.HashPassword()

			if tt.wantErr {
				if err == nil {
					t.Errorf("HashPassword() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("HashPassword() unexpected error: %v", err)
				return
			}

			// パスワードがハッシュ化されていることを確認
			if user.Password == tt.password {
				t.Errorf("Password was not hashed")
			}

			// ハッシュ化されたパスワードが空でないことを確認
			if user.Password == "" {
				t.Errorf("Hashed password is empty")
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		originalPass string
		comparePass  string
		wantErr      bool
		description  string
	}{
		{
			name:         "correct password",
			originalPass: "password123",
			comparePass:  "password123",
			wantErr:      false,
			description:  "should succeed with correct password",
		},
		{
			name:         "incorrect password",
			originalPass: "password123",
			comparePass:  "wrongpassword",
			wantErr:      true,
			description:  "should fail with incorrect password",
		},
		{
			name:         "empty password",
			originalPass: "password123",
			comparePass:  "",
			wantErr:      true,
			description:  "should fail with empty password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := &User{
				Password: tt.originalPass,
			}

			// パスワードをハッシュ化
			if err := user.HashPassword(); err != nil {
				t.Fatalf("HashPassword() failed: %v", err)
			}

			// パスワード比較
			err := user.ComparePassword(tt.comparePass)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ComparePassword() expected error but got none: %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("ComparePassword() unexpected error: %v - %s", err, tt.description)
			}
		})
	}
}

func TestUserID_MarshalBinary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  UserID
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "zero user ID",
			userID:  0,
			wantErr: false,
		},
		{
			name:    "negative user ID",
			userID:  -1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.userID.MarshalBinary()

			if tt.wantErr {
				if err == nil {
					t.Errorf("MarshalBinary() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("MarshalBinary() unexpected error: %v", err)
				return
			}

			// 結果が期待される形式であることを確認
			expected := []byte("123")
			if tt.userID == 123 {
				if string(got) != string(expected) {
					t.Errorf("MarshalBinary() got %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestUserID_UnmarshalBinary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		want    UserID
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    []byte("123"),
			want:    123,
			wantErr: false,
		},
		{
			name:    "zero data",
			data:    []byte("0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid data",
			data:    []byte("invalid"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var userID UserID
			err := userID.UnmarshalBinary(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalBinary() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalBinary() unexpected error: %v", err)
				return
			}

			if userID != tt.want {
				t.Errorf("UnmarshalBinary() got %v, want %v", userID, tt.want)
			}
		})
	}
}

func TestUser_StructFields(t *testing.T) {
	t.Parallel()

	now := time.Now()
	user := &User{
		ID:         1,
		Name:       "testuser",
		Password:   "hashedpassword",
		Role:       "user",
		CreatedAt:  now,
		ModifiedAt: now,
	}

	// 各フィールドが正しく設定されていることを確認
	if user.ID != 1 {
		t.Errorf("ID = %v, want 1", user.ID)
	}

	if user.Name != "testuser" {
		t.Errorf("Name = %v, want testuser", user.Name)
	}

	if user.Password != "hashedpassword" {
		t.Errorf("Password = %v, want hashedpassword", user.Password)
	}

	if user.Role != "user" {
		t.Errorf("Role = %v, want user", user.Role)
	}

	if !user.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt, now)
	}

	if !user.ModifiedAt.Equal(now) {
		t.Errorf("ModifiedAt = %v, want %v", user.ModifiedAt, now)
	}
}
