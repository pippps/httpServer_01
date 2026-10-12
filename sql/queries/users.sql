-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid() ,
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT
    *
FROM users
WHERE email = $1;

-- name: ModUser :one
UPDATE users
SET
    email = $2,
    updated_at = NOW(),
    hashed_password = $3
WHERE id = $1
RETURNING
    id,
    created_at,
    updated_at,
    email,
    is_chirpy_red;

-- name: UpgradeChirpyRed :exec
UPDATE users
SET
    is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1;
