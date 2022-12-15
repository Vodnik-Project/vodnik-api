-- name: AddUserToTask :one
INSERT INTO usersintask (
    user_id, task_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTasksOfUser :many
SELECT * FROM usersintask
WHERE user_id = $1;

-- name: GetUsersOfTask :many
SELECT * FROM usersintask
WHERE task_id = $1;

-- name: DeleteUserFromTask :exec
DELETE FROM usersintask
WHERE user_id = $1 AND task_id = $2;