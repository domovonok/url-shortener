-- name: Create :one
INSERT INTO url (url)
VALUES ($1)
ON CONFLICT (url) DO NOTHING
RETURNING id;

-- name: GetIDByURL :one
SELECT id
FROM url
WHERE url = $1;

-- name: Get :one
SELECT url
FROM url
WHERE id = $1;
