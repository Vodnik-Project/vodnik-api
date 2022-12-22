-- TODO: auto time for beggining without input doesn't work.

-- name: CreateTask :one
INSERT INTO tasks (
    project_id, title, info, tag, created_by, beggining, deadline, color
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetTasksByProjectID :many
SELECT * FROM tasks
WHERE project_id = $1;

-- name: GetTaskData :one
SELECT * FROM tasks
WHERE task_id = $1;

-- TODO: get tasks by filtering: title, tag, created_by, beggining, deadline

-- name: UpdateTask :one
UPDATE tasks SET
  title =     COALESCE(NULLIF(@title, ''), title),
  info =      COALESCE(NULLIF(@info, ''), info),
  tag =       COALESCE(NULLIF(@tag, ''), tag),
  beggining = COALESCE(NULLIF(@beggining, timestamptz '0001-01-01 03:25:44+03:25:44'), beggining),
  deadline =  COALESCE(NULLIF(@deadline, timestamptz '0001-01-01 03:25:44+03:25:44'), deadline),
  color =     COALESCE(NULLIF(@color, ''), color)
WHERE task_id = @task_id
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1;