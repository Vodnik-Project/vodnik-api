-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetProjectsByUserID :many
SELECT * FROM usersinproject
WHERE user_id = $1;

-- name: GetUsersByProjectID :many
SELECT * FROM usersinproject
WHERE project_id = $1;

-- name: DeleteUserFromProject :exec
DELETE FROM usersinproject
WHERE user_id = $1 AND project_id = $2;