-- name: SetSession :exec
INSERT INTO refresh_token (
    token, username, fingerprint
) VALUES (
    $1, $2, $3
);

-- name: GetSessionByToken :one
SELECT * FROM refresh_token
WHERE token = $1;

-- name: GetSessionByUsername :one
SELECT * FROM refresh_token
WHERE username = $1;

-- name: DeleteSession :exec
DELETE FROM refresh_token
WHERE token = $1;