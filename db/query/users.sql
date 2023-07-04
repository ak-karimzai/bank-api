-- name: CreateUser :one
INSERT INTO users (username, hashed_pwd, full_name, email)
VALUES ($1, $2, $3, $4) Returning *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1;
