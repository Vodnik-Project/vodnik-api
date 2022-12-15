-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetProjectsOfUser :many
SELECT * FROM usersinproject
WHERE user_id = $1;

-- name: GetUsersOfProject :many
SELECT * FROM usersinproject
WHERE project_id = $1;

-- name: DeleteUserFromProject :exec
DELETE FROM usersinproject
WHERE user_id = $1 AND project_id = $2;