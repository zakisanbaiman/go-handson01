package entity

import (
	"testing"
	"time"
)

func TestTaskID_MarshalBinary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		taskID  TaskID
		wantErr bool
	}{
		{
			name:    "valid task ID",
			taskID:  456,
			wantErr: false,
		},
		{
			name:    "zero task ID",
			taskID:  0,
			wantErr: false,
		},
		{
			name:    "negative task ID",
			taskID:  -1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.taskID.MarshalBinary()

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
			expected := []byte("456")
			if tt.taskID == 456 {
				if string(got) != string(expected) {
					t.Errorf("MarshalBinary() got %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestTaskID_UnmarshalBinary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		want    TaskID
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    []byte("456"),
			want:    456,
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

			var taskID TaskID
			err := taskID.UnmarshalBinary(tt.data)

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

			if taskID != tt.want {
				t.Errorf("UnmarshalBinary() got %v, want %v", taskID, tt.want)
			}
		})
	}
}

func TestTaskStatus_Constants(t *testing.T) {
	t.Parallel()

	// タスクステータスの定数が正しく定義されていることを確認
	if TaskStatusTodo != "todo" {
		t.Errorf("TaskStatusTodo = %v, want todo", TaskStatusTodo)
	}

	if TaskStatusDoing != "doing" {
		t.Errorf("TaskStatusDoing = %v, want doing", TaskStatusDoing)
	}

	if TaskStatusDone != "done" {
		t.Errorf("TaskStatusDone = %v, want done", TaskStatusDone)
	}
}

func TestTask_StructFields(t *testing.T) {
	t.Parallel()

	now := time.Now()
	task := &Task{
		ID:         1,
		UserID:     123,
		Title:      "Test Task",
		Status:     TaskStatusTodo,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	// 各フィールドが正しく設定されていることを確認
	if task.ID != 1 {
		t.Errorf("ID = %v, want 1", task.ID)
	}

	if task.UserID != 123 {
		t.Errorf("UserID = %v, want 123", task.UserID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Title = %v, want Test Task", task.Title)
	}

	if task.Status != TaskStatusTodo {
		t.Errorf("Status = %v, want %v", task.Status, TaskStatusTodo)
	}

	if !task.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", task.CreatedAt, now)
	}

	if !task.ModifiedAt.Equal(now) {
		t.Errorf("ModifiedAt = %v, want %v", task.ModifiedAt, now)
	}
}

func TestTasks_Slice(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tasks := Tasks{
		{
			ID:         1,
			UserID:     123,
			Title:      "Task 1",
			Status:     TaskStatusTodo,
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:         2,
			UserID:     123,
			Title:      "Task 2",
			Status:     TaskStatusDone,
			CreatedAt:  now,
			ModifiedAt: now,
		},
	}

	// スライスの長さを確認
	if len(tasks) != 2 {
		t.Errorf("Tasks length = %v, want 2", len(tasks))
	}

	// 各タスクの内容を確認
	if tasks[0].Title != "Task 1" {
		t.Errorf("First task title = %v, want Task 1", tasks[0].Title)
	}

	if tasks[1].Title != "Task 2" {
		t.Errorf("Second task title = %v, want Task 2", tasks[1].Title)
	}

	if tasks[0].Status != TaskStatusTodo {
		t.Errorf("First task status = %v, want %v", tasks[0].Status, TaskStatusTodo)
	}

	if tasks[1].Status != TaskStatusDone {
		t.Errorf("Second task status = %v, want %v", tasks[1].Status, TaskStatusDone)
	}
}
