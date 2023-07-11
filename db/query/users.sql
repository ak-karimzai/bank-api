-- name: CreateUser :one
INSERT INTO users (username, hashed_pwd, full_name, email)
VALUES ($1, $2, $3, $4) Returning *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1;

-- name: UpdateUser :one
UPDATE 
  users 
SET 
  hashed_pwd = coalesce(sqlc.narg('hashed_pwd'), hashed_pwd),
  pwd_changed_at = coalesce(sqlc.narg('pwd_changed_at'), pwd_changed_at), 
  full_name = coalesce(sqlc.narg('full_name'), full_name), 
  email = coalesce(sqlc.narg('email'), email)
WHERE
  username = sqlc.arg(username)
Returning *;