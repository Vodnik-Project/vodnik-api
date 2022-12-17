-- name: AddUserToTask :one
INSERT INTO usersintask (
    user_id, task_id, admin
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTasksByUserID :many
SELECT * FROM usersintask
WHERE user_id = $1;

-- name: GetUsersByTaskID :many
SELECT * FROM usersintask
WHERE task_id = $1;

-- name: DeleteUserFromTask :exec
DELETE FROM usersintask
WHERE user_id = $1 AND task_id = $2;