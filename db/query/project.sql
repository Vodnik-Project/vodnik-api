-- name: CreateProject :one
INSERT INTO projects (
    title, info, owner_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetProjectData :one
SELECT * FROM projects
WHERE project_id = $1;

-- name: UpdateProject :one
UPDATE projects SET
  title = COALESCE(NULLIF(@title, ''), title),
  info = COALESCE(NULLIF(@info, ''), info)
WHERE project_id = @project_id
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;