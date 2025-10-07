package service

import (
	"context"
	"errors"
	"testing"

	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

func TestLogin_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		userName       string
		password       string
		mockUser       *entity.User
		mockUserError  error
		mockToken      []byte
		mockTokenError error
		wantToken      string
		wantError      bool
	}{
		{
			name:     "successful login",
			userName: "testuser",
			password: "password123",
			mockUser: func() *entity.User {
				user := &entity.User{
					ID:       1,
					Name:     "testuser",
					Password: "password123", // テスト用に平文で設定
					Role:     "user",
				}
				if err := user.HashPassword(); err != nil {
					panic(err) // テスト用なのでpanicで十分
				}
				return user
			}(),
			mockToken: []byte("mock-jwt-token"),
			wantToken: "mock-jwt-token",
			wantError: false,
		},
		{
			name:          "user not found",
			userName:      "nonexistent",
			password:      "password123",
			mockUserError: errors.New("user not found"),
			wantError:     true,
		},
		{
			name:     "wrong password",
			userName: "testuser",
			password: "wrongpassword",
			mockUser: func() *entity.User {
				user := &entity.User{
					ID:       1,
					Name:     "testuser",
					Password: "password123", // テスト用に平文で設定
					Role:     "user",
				}
				if err := user.HashPassword(); err != nil {
					panic(err) // テスト用なのでpanicで十分
				}
				return user
			}(),
			wantError: true,
		},
		{
			name:     "token generation error",
			userName: "testuser",
			password: "password123",
			mockUser: func() *entity.User {
				user := &entity.User{
					ID:       1,
					Name:     "testuser",
					Password: "password123", // テスト用に平文で設定
					Role:     "user",
				}
				if err := user.HashPassword(); err != nil {
					panic(err) // テスト用なのでpanicで十分
				}
				return user
			}(),
			mockTokenError: errors.New("token generation failed"),
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// モックの設定
			mockUserGetter := &UserGetterMock{
				GetUserFunc: func(ctx context.Context, db store.Queryer, userName string) (*entity.User, error) {
					if tt.mockUserError != nil {
						return nil, tt.mockUserError
					}
					return tt.mockUser, nil
				},
			}

			mockTokenGenerator := &TokenGeneratorMock{
				GenerateTokenFunc: func(ctx context.Context, user entity.User) ([]byte, error) {
					if tt.mockTokenError != nil {
						return nil, tt.mockTokenError
					}
					return tt.mockToken, nil
				},
			}

			// サービスインスタンスの作成
			loginService := &Login{
				Repo:           mockUserGetter,
				TokenGenerator: mockTokenGenerator,
			}

			// テスト実行
			ctx := context.Background()
			gotToken, err := loginService.Login(ctx, tt.userName, tt.password)

			// 結果の検証
			if tt.wantError {
				if err == nil {
					t.Errorf("Login() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
				return
			}

			if gotToken != tt.wantToken {
				t.Errorf("Login() got token = %v, want %v", gotToken, tt.wantToken)
			}

			// モックの呼び出し回数を検証
			if len(mockUserGetter.GetUserCalls()) != 1 {
				t.Errorf("GetUser() was called %d times, want 1", len(mockUserGetter.GetUserCalls()))
			}

			if tt.mockTokenError == nil && len(mockTokenGenerator.GenerateTokenCalls()) != 1 {
				t.Errorf("GenerateToken() was called %d times, want 1", len(mockTokenGenerator.GenerateTokenCalls()))
			}
		})
	}
}
