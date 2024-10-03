-- queries.sql
-- name: CreateTodo :execresult
INSERT INTO todos (body, completed) VALUES (?, ?);

-- name: ListTodos :many
SELECT * FROM todos ORDER BY created_at DESC;

-- name: UpdateTodo :execresult
UPDATE todos SET completed = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;

-- name: GetTodoByID :one
SELECT id, body, completed, created_at FROM todos WHERE id = ?;