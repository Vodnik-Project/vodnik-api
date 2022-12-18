-- name: CreateUser :one
INSERT INTO users (
    username, email, pass_hash, bio
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET
  username = COALESCE(NULLIF(@username, ''), username),
  email = COALESCE(NULLIF(@email, ''), email),
  pass_hash = COALESCE(NULLIF(@pass_hash, ''), pass_hash),
  bio = COALESCE(NULLIF(@bio, ''), bio)
WHERE user_id = @user_id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;