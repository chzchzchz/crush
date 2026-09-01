-- name: CreateSession :one
INSERT INTO sessions (
    id,
    parent_session_id,
    title,
    message_count,
    prompt_tokens,
    completion_tokens,
    cached_tokens,
    reasoning_tokens,
    total_requests,
    total_prompt_tokens,
    total_completion_tokens,
    total_cached_tokens,
    total_reasoning_tokens,
    cost,
    summary_message_id,
    updated_at,
    created_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    null,
    strftime('%s', 'now'),
    strftime('%s', 'now')
) RETURNING *;

-- name: GetSessionByID :one
SELECT *
FROM sessions
WHERE id = ? LIMIT 1;

-- name: GetLastSession :one
SELECT *
FROM sessions
ORDER BY updated_at DESC
LIMIT 1;

-- name: ListSessions :many
SELECT *
FROM sessions
WHERE parent_session_id is NULL
ORDER BY updated_at DESC;

-- name: UpdateSession :one
UPDATE sessions
SET
    title = ?,
    prompt_tokens = ?,
    completion_tokens = ?,
    cached_tokens = ?,
    reasoning_tokens = ?,
    total_requests = ?,
    total_prompt_tokens = ?,
    total_completion_tokens = ?,
    total_cached_tokens = ?,
    total_reasoning_tokens = ?,
    summary_message_id = ?,
    cost = ?,
    todos = ?,
    updated_at = strftime('%s', 'now')
WHERE id = ?
RETURNING *;

-- name: UpdateSessionTitleAndUsage :exec
UPDATE sessions
SET
    title = ?,
    prompt_tokens = prompt_tokens + ?,
    completion_tokens = completion_tokens + ?,
    cached_tokens = cached_tokens + ?,
    reasoning_tokens = reasoning_tokens + ?,
    total_requests = total_requests + ?,
    total_prompt_tokens = total_prompt_tokens + ?,
    total_completion_tokens = total_completion_tokens + ?,
    total_cached_tokens = total_cached_tokens + ?,
    total_reasoning_tokens = total_reasoning_tokens + ?,
    cost = cost + ?,
    updated_at = strftime('%s', 'now')
WHERE id = ?;


-- name: RenameSession :exec
UPDATE sessions
SET
    title = ?
WHERE id = ?;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;
