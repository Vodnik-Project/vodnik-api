-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id, admin
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProjectsByUserID :many
SELECT projects.project_id, projects.title, projects.info, projects.owner_id, projects.created_at, 
       usersinproject.admin 
FROM projects
INNER JOIN usersinproject 
ON projects.project_id=usersinproject.project_id 
WHERE usersinproject.user_id=$1;

-- name: GetUsersByProjectID :many
SELECT users.user_id, users.username, users.bio,
       usersinproject.project_id, usersinproject.added_at, usersinproject.admin 
FROM users
INNER JOIN usersinproject
ON users.user_id=usersinproject.user_id
WHERE usersinproject.project_id = $1;

-- name: IsUserInProject :one
SELECT * FROM usersinproject
WHERE user_id = $1 AND project_id = $2;

-- name: DeleteUserFromProject :exec
DELETE FROM usersinproject
WHERE user_id = $1 AND project_id = $2;