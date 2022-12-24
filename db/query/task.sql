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
WHERE project_id = $1 AND
      title =       COALESCE(NULLIF(@title, ''), title) AND
      info =        COALESCE(NULLIF(@info, ''), info) AND
      tag =         COALESCE(NULLIF(@tag, ''), tag) AND
      created_by =  COALESCE(NULLIF(@created_by, uuid '00000000-0000-0000-0000-000000000000'), created_by) AND
      created_at >= COALESCE(NULLIF(@created_at_from, timestamptz '0001-01-01 03:25:44+03:25:44'), created_at) AND
      created_at <= COALESCE(NULLIF(@created_at_until, timestamptz '0001-01-01 03:25:44+03:25:44'), created_at) AND
      beggining >=  COALESCE(NULLIF(@beggining_from, timestamptz '0001-01-01 03:25:44+03:25:44'), beggining) AND
      beggining <=  COALESCE(NULLIF(@beggining_until, timestamptz '0001-01-01 03:25:44+03:25:44'), beggining) AND
      deadline >=   COALESCE(NULLIF(@deadline_from, timestamptz '0001-01-01 03:25:44+03:25:44'), deadline) AND
      deadline <=   COALESCE(NULLIF(@deadline_until, timestamptz '0001-01-01 03:25:44+03:25:44'), deadline)
ORDER BY
  CASE WHEN @sortDirection = 'asc' THEN
    CASE
      WHEN @sortBy = 'created_at' THEN created_at
      WHEN @sortBy = 'beggining' THEN beggining
      WHEN @sortBy = 'deadline' THEN deadline
    END
  END ASC,
  CASE WHEN @sortDirection = 'desc' THEN
    CASE
      WHEN @sortBy = 'created_at' THEN created_at
      WHEN @sortBy = 'beggining' THEN beggining
      WHEN @sortBy = 'deadline' THEN deadline
    END
  END DESC
LIMIT $2
OFFSET $3;

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