-- TODO: auto time for beggining without input doesn't work.

-- name: CreateTask :one
INSERT INTO tasks (
    project_id, title, info, tag, created_by, beggining, deadline, color
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetTasksInProject :many
SELECT * FROM tasks
WHERE project_id = $1;

-- name: GetTaskData :one
SELECT * FROM tasks
WHERE task_id = $1;

-- TODO: get tasks by filtering: title, tag, created_by, beggining, deadline

-- name: UpdateTask :one
UPDATE tasks SET
  title =     COALESCE(NULLIF(@title, 'NULL'), title),
  info =      COALESCE(NULLIF(@info, 'NULL'), info),
  tag =       COALESCE(NULLIF(@tag, 'NULL'), tag),
  beggining = COALESCE(@beggining, beggining),
  deadline =  COALESCE(@deadline, deadline),
  color =     COALESCE(NULLIF(@color, 'NULL'), color)
WHERE task_id = @id
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1;