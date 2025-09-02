-- name: GetUserByEmailOrId :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email
FROM
    users u
WHERE
    u.email = $1
    OR u.id = $2;
