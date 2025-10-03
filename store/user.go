package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/zakisanbaiman/go-handson01/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, user *entity.User) error {
	user.CreatedAt = r.Clocker.Now()
	user.ModifiedAt = r.Clocker.Now()

	sql := `INSERT INTO users (name, password, role, created_at, modified_at) VALUES (?, ?, ?, ?, ?);`

	result, err := db.ExecContext(ctx, sql, user.Name, user.Password, user.Role, user.CreatedAt, user.ModifiedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeSQLDuplicateEntry {
			return fmt.Errorf("failed to register user: %w", ErrAlreadyExists)
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = entity.UserID(id)

	return nil
}

func (r *Repository) GetUser(
	ctx context.Context,
	db Queryer,
	name string,
) (*entity.User, error) {
	user := &entity.User{}
	sql := `SELECT id, name, password, role, created_at, modified_at FROM users WHERE name = ?;`
	if err := db.GetContext(ctx, user, sql, name); err != nil {
		return nil, err
	}
	return user, nil
}
