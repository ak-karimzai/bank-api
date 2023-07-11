// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, hashed_pwd, full_name, email)
VALUES ($1, $2, $3, $4) Returning username, hashed_pwd, full_name, email, pwd_changed_at, created_at
`

type CreateUserParams struct {
	Username  string `json:"username"`
	HashedPwd string `json:"hashed_pwd"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPwd,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPwd,
		&i.FullName,
		&i.Email,
		&i.PwdChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_pwd, full_name, email, pwd_changed_at, created_at FROM users WHERE username = $1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPwd,
		&i.FullName,
		&i.Email,
		&i.PwdChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE 
  users 
SET 
  hashed_pwd = coalesce($1, hashed_pwd),
  pwd_changed_at = coalesce($2, pwd_changed_at), 
  full_name = coalesce($3, full_name), 
  email = coalesce($4, email)
WHERE
  username = $5
Returning username, hashed_pwd, full_name, email, pwd_changed_at, created_at
`

type UpdateUserParams struct {
	HashedPwd    pgtype.Text        `json:"hashed_pwd"`
	PwdChangedAt pgtype.Timestamptz `json:"pwd_changed_at"`
	FullName     pgtype.Text        `json:"full_name"`
	Email        pgtype.Text        `json:"email"`
	Username     string             `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.HashedPwd,
		arg.PwdChangedAt,
		arg.FullName,
		arg.Email,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPwd,
		&i.FullName,
		&i.Email,
		&i.PwdChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
