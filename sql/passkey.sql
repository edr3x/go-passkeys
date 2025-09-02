-- name: GetUserPasskeyCredentials :many
SELECT
    *
FROM
    credentials c
WHERE
    c.user_id = $1;

-- name: AddCredential :exec
INSERT INTO
    credentials (
        id,
        user_id,
        public_key,
        sign_count,
        transports,
        attestation_type,
        aaguid,
        attestation,
        flags,
        clone_warning,
        attachment
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateCredential :exec
UPDATE
    credentials
SET
    sign_count = $2,
    updated_at = NOW()
WHERE
    id = $1;
