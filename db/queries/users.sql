-- name: CreateDestino :one
INSERT INTO destinos (nombre_lugar, fecha_viaje, valoracion, descripcion)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDestino :one
SELECT * FROM destinos WHERE id = $1;

-- name: ListDestinos :many
SELECT * FROM destinos ORDER BY fecha_viaje DESC;

-- name: UpdateDestino :exec
UPDATE destinos 
SET nombre_lugar = $2, fecha_viaje = $3, valoracion = $4, descripcion = $5 
WHERE id = $1;

-- name: DeleteDestino :exec
DELETE FROM destinos WHERE id = $1;