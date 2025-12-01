-- name: MonitorLogs :many
SELECT * FROM logs ORDER BY created_at DESC;
