-- name: CreateProject :one
INSERT INTO projects (
    title, info, owner_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetProjectsByUserId :many
SELECT * FROM projects
WHERE owner_id = $1;

-- name: GetProjectData :one
SELECT * FROM projects
WHERE project_id = $1;

-- name: UpdateProject :one
UPDATE projects SET
  title = COALESCE(NULLIF(@title, 'NULL'), title),
  info = COALESCE(NULLIF(@info, 'NULL'), info)
WHERE project_id = @project_id
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;