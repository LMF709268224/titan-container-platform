package dao

import (
	"context"
	"fmt"

	"titan-container-platform/core"
)

// GetUserResponse 获取用户响应
func GetUserResponse(ctx context.Context, account string) (*core.ResponseUser, error) {
	response := core.ResponseUser{}

	if err := mDB.QueryRowxContext(ctx, fmt.Sprintf(
		`SELECT id,user_name,user_email,created_at FROM %s WHERE account = ?`, userInfoTable), account,
	).StructScan(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUser creates a new user in the database.
func CreateUser(ctx context.Context, user *core.User) error {
	query := fmt.Sprintf(`INSERT INTO %s (account, user_name, user_email, avatar) VALUES (:account, :user_name, :user_email, :avatar) `, userInfoTable)
	_, err := mDB.NamedExec(query, user)

	return err
}

// GetUserByAccount retrieves a user by their account from the database.
func GetUserByAccount(ctx context.Context, account string) (*core.User, error) {
	var out core.User
	if err := mDB.QueryRowxContext(ctx, fmt.Sprintf(
		`SELECT * FROM %s WHERE account = ?`, userInfoTable), account,
	).StructScan(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
