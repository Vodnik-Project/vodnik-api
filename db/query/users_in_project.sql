-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id, admin
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProjectsByUserID :many
SELECT * FROM usersinproject
WHERE user_id = $1;

-- name: GetUsersByProjectID :many
SELECT * FROM usersinproject
WHERE project_id = $1;

-- name: IsUserInProject :one
SELECT * FROM usersinproject
WHERE user_id = $1 AND project_id = $2;

-- name: DeleteUserFromProject :exec
DELETE FROM usersinproject
WHERE user_id = $1 AND project_id = $2;