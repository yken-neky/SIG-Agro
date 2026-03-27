-- Notification queries

-- name: CreateNotification :one
INSERT INTO notifications (user_id, notification_type, channel, title, message, metadata)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, notification_type, channel, title, message, read, created_at;

-- name: GetNotification :one
SELECT id, user_id, notification_type, channel, title, message, read, created_at
FROM notifications WHERE id = $1;

-- name: ListNotifications :many
SELECT id, user_id, notification_type, channel, title, message, read, created_at
FROM notifications WHERE user_id = $1 AND ($2::boolean = FALSE OR read = FALSE)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: MarkAsRead :exec
UPDATE notifications SET read = TRUE WHERE id = $1;
