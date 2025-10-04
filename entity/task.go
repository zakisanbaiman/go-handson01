package entity

import (
	"fmt"
	"strconv"
	"time"
)

type TaskID int64

// MarshalBinary implements encoding.BinaryMarshaler
func (t TaskID) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(t), 10)), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (t *TaskID) UnmarshalBinary(data []byte) error {
	id, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TaskID: %w", err)
	}
	*t = TaskID(id)
	return nil
}

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

type Task struct {
	ID         TaskID     `json:"id" db:"id"`
	UserID     UserID     `json:"user_id" db:"user_id"`
	Title      string     `json:"title" db:"title"`
	Status     TaskStatus `json:"status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt time.Time  `json:"modified_at" db:"modified_at"`
}

type Tasks []*Task
