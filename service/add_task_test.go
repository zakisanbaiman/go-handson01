package service

import (
	"context"
	"errors"
	"testing"

	"github.com/zakisanbaiman/go-handson01/auth"
	"github.com/zakisanbaiman/go-handson01/entity"
	"github.com/zakisanbaiman/go-handson01/store"
)

func TestAddTask_AddTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		title       string
		userID      entity.UserID
		userIDFound bool
		mockError   error
		wantError   bool
		wantTask    *entity.Task
	}{
		{
			name:        "successful task creation",
			title:       "Test Task",
			userID:      1,
			userIDFound: true,
			mockError:   nil,
			wantError:   false,
			wantTask: &entity.Task{
				UserID: 1,
				Title:  "Test Task",
				Status: entity.TaskStatusTodo,
			},
		},
		{
			name:        "user ID not found in context",
			title:       "Test Task",
			userIDFound: false,
			wantError:   true,
		},
		{
			name:        "repository error",
			title:       "Test Task",
			userID:      1,
			userIDFound: true,
			mockError:   errors.New("database error"),
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// モックの設定
			mockRepo := &TaskAdderMock{
				AddTaskFunc: func(ctx context.Context, db store.Execer, t *entity.Task) error {
					if tt.mockError != nil {
						return tt.mockError
					}
					// 成功時はIDを設定
					t.ID = 1
					return nil
				},
			}

			// サービスインスタンスの作成
			addTaskService := &AddTask{
				Repo: mockRepo,
			}

			// コンテキストの設定
			ctx := context.Background()
			if tt.userIDFound {
				ctx = auth.SetUserID(ctx, tt.userID)
			}

			// テスト実行
			gotTask, err := addTaskService.AddTask(ctx, tt.title)

			// 結果の検証
			if tt.wantError {
				if err == nil {
					t.Errorf("AddTask() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("AddTask() unexpected error: %v", err)
				return
			}

			if gotTask == nil {
				t.Errorf("AddTask() got nil task")
				return
			}

			if gotTask.UserID != tt.wantTask.UserID {
				t.Errorf("AddTask() got userID = %v, want %v", gotTask.UserID, tt.wantTask.UserID)
			}

			if gotTask.Title != tt.wantTask.Title {
				t.Errorf("AddTask() got title = %v, want %v", gotTask.Title, tt.wantTask.Title)
			}

			if gotTask.Status != tt.wantTask.Status {
				t.Errorf("AddTask() got status = %v, want %v", gotTask.Status, tt.wantTask.Status)
			}

			// モックの呼び出し回数を検証
			if tt.userIDFound && tt.mockError == nil {
				if len(mockRepo.AddTaskCalls()) != 1 {
					t.Errorf("AddTask() was called %d times, want 1", len(mockRepo.AddTaskCalls()))
				}
			}
		})
	}
}
