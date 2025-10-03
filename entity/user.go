package entity

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
