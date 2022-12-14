-- name: SetSession :exec
INSERT INTO refresh_token (
    token, user_id, fingerprint, device
) VALUES (
    $1, $2, $3, $4
);

-- name: GetSessionByToken :one
SELECT * FROM refresh_token
WHERE token = $1;

-- name: GetDeviceSession :one
SELECT * FROM refresh_token
WHERE user_id = $1 AND fingerprint = $2;

-- name: DeleteSession :exec
DELETE FROM refresh_token
WHERE token = $1;