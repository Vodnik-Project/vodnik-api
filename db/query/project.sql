-- name: CreateProject :one
INSERT INTO projects (
    title, info, owner_id
) VALUES (
    @title, @info, @owner_id
)
RETURNING *;

-- name: GetProjectsByUserId :many
SELECT * FROM projects
WHERE owner_id = @owner_id;

-- name: GetProjectData :one
SELECT * FROM projects
WHERE id = @id;

-- name: UpdateProject :one
UPDATE projects SET
  title = COALESCE(NULLIF(@title, 'NULL'), title),
  info = COALESCE(NULLIF(@info, 'NULL'), info)
WHERE id = @id
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = @id;