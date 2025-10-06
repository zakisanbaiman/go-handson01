package service

import (
	"context"
	"errors"
	"testing"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

// RegisterUserRepositoryMock は RegisterUserRepository のモック
type RegisterUserRepositoryMock struct {
	RegisterUserFunc func(ctx context.Context, db store.Execer, user *entity.User) error
}

func (m *RegisterUserRepositoryMock) RegisterUser(ctx context.Context, db store.Execer, user *entity.User) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, db, user)
	}
	return nil
}

func TestRegisterUser_RegisterUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userName  string
		password  string
		role      string
		mockError error
		wantError bool
	}{
		{
			name:      "successful user registration",
			userName:  "testuser",
			password:  "password123",
			role:      "user",
			mockError: nil,
			wantError: false,
		},
		{
			name:      "repository error",
			userName:  "testuser",
			password:  "password123",
			role:      "user",
			mockError: errors.New("database error"),
			wantError: true,
		},
		{
			name:      "duplicate user error",
			userName:  "existinguser",
			password:  "password123",
			role:      "user",
			mockError: store.ErrAlreadyExists,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// モックの設定
			mockRepo := &RegisterUserRepositoryMock{
				RegisterUserFunc: func(ctx context.Context, db store.Execer, user *entity.User) error {
					return tt.mockError
				},
			}

			// サービスインスタンスの作成
			registerUserService := &RegisterUser{
				Repo: mockRepo,
			}

			// テスト実行
			ctx := context.Background()
			gotUser, err := registerUserService.RegisterUser(ctx, tt.userName, tt.password, tt.role)

			// 結果の検証
			if tt.wantError {
				if err == nil {
					t.Errorf("RegisterUser() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("RegisterUser() unexpected error: %v", err)
				return
			}

			if gotUser == nil {
				t.Errorf("RegisterUser() got nil user")
				return
			}

			if gotUser.Name != tt.userName {
				t.Errorf("RegisterUser() got name = %v, want %v", gotUser.Name, tt.userName)
			}

			if gotUser.Role != tt.role {
				t.Errorf("RegisterUser() got role = %v, want %v", gotUser.Role, tt.role)
			}

			// パスワードがハッシュ化されていることを確認
			if gotUser.Password == tt.password {
				t.Errorf("RegisterUser() password was not hashed")
			}

			if gotUser.Password == "" {
				t.Errorf("RegisterUser() password is empty")
			}
		})
	}
}

func TestRegisterUser_RegisterUser_PasswordHashing(t *testing.T) {
	t.Parallel()

	// パスワードハッシュ化のテスト
	mockRepo := &RegisterUserRepositoryMock{
		RegisterUserFunc: func(ctx context.Context, db store.Execer, user *entity.User) error {
			return nil
		},
	}

	registerUserService := &RegisterUser{
		Repo: mockRepo,
	}

	ctx := context.Background()
	user, err := registerUserService.RegisterUser(ctx, "testuser", "password123", "user")

	if err != nil {
		t.Fatalf("RegisterUser() unexpected error: %v", err)
	}

	// パスワードがハッシュ化されていることを確認
	if user.Password == "password123" {
		t.Errorf("Password was not hashed")
	}

	// ハッシュ化されたパスワードで認証できることを確認
	if err := user.ComparePassword("password123"); err != nil {
		t.Errorf("Hashed password comparison failed: %v", err)
	}

	// 間違ったパスワードで認証が失敗することを確認
	if err := user.ComparePassword("wrongpassword"); err == nil {
		t.Errorf("Wrong password should fail authentication")
	}
}
