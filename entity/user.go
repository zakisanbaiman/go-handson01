package entity

import (
	"fmt"
	"strconv"
	"time"
)

type UserID int64

// MarshalBinary implements encoding.BinaryMarshaler
func (u UserID) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(u), 10)), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (u *UserID) UnmarshalBinary(data []byte) error {
	id, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal UserID: %w", err)
	}
	*u = UserID(id)
	return nil
}

type User struct {
	ID         UserID    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Password   string    `json:"password" db:"password"`
	Role       string    `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}
