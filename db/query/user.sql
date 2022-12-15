-- name: CreateUser :one
INSERT INTO users (
    username, email, pass_hash, bio
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users SET
  username = COALESCE(NULLIF(@username, 'NULL'), username),
  email = COALESCE(NULLIF(@email, 'NULL'), email),
  bio = COALESCE(NULLIF(@bio, 'NULL'), bio)
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;