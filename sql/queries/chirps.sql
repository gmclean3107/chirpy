-- name: CreateChirp :one
INSERT INTO chirps (user_id, body, updated_at)
VALUES (
    $1,
    $2,
    now()
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;