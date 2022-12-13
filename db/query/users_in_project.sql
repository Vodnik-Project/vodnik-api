-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetUserProjects :many
SELECT * FROM usersinproject
WHERE user_id = $1;

-- name: GetProjectUsers :many
SELECT * FROM usersinproject
WHERE project_id = $1;