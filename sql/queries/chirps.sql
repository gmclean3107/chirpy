-- name: CreateChirp :one
INSERT INTO chirps (user_id, body, updated_at)
VALUES (
    $1,
    $2,
    now()
)
RETURNING *;