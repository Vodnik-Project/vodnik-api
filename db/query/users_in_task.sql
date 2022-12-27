-- name: AddUserToTask :one
INSERT INTO usersintask (
    user_id, task_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTasksByUserID :many
SELECT tasks.task_id, tasks.title, tasks.info, tasks.tag, 
        tasks.created_by, tasks.created_at, tasks.beggining, 
        tasks.deadline, tasks.color
FROM tasks
INNER JOIN usersintask 
ON tasks.task_id=usersintask.task_id 
WHERE usersintask.user_id=$1 AND tasks.project_id=$2;

-- name: GetUsersByTaskID :many
SELECT users.user_id, users.username, users.bio,
       usersintask.added_at
FROM users
INNER JOIN usersintask
ON users.user_id=usersintask.user_id
WHERE usersintask.task_id = $1;

-- name: DeleteUserFromTask :exec
DELETE FROM usersintask
WHERE user_id = $1 AND task_id = $2;