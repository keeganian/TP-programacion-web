-- ==========================================
-- USUARIOS
-- ==========================================

-- name: CreateUsuario :one
INSERT INTO usuarios (nombre, email)
VALUES ($1, $2)
RETURNING *;

-- name: GetUsuario :one
SELECT * FROM usuarios WHERE id = $1;

-- name: ListUsuarios :many
SELECT * FROM usuarios ORDER BY id;

-- name: UpdateUsuario :exec
UPDATE usuarios 
SET nombre = $2, email = $3 
WHERE id = $1;

-- name: DeleteUsuario :exec
DELETE FROM usuarios WHERE id = $1;


-- ==========================================
-- CATEGORÍAS
-- ==========================================

-- name: CreateCategoria :one
INSERT INTO categorias (nombre)
VALUES ($1)
RETURNING *;

-- name: GetCategoria :one
SELECT * FROM categorias WHERE id = $1;

-- name: ListCategorias :many
SELECT * FROM categorias ORDER BY id;

-- name: UpdateCategoria :exec
UPDATE categorias 
SET nombre = $2 
WHERE id = $1;

-- name: DeleteCategoria :exec
DELETE FROM categorias WHERE id = $1;


-- ==========================================
-- DESTINOS
-- ==========================================

-- name: CreateDestino :one
INSERT INTO destinos (nombre_lugar, fecha_viaje, valoracion, descripcion, categoria_id, usuario_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetDestino :one
SELECT * FROM destinos WHERE id = $1;

-- name: ListDestinos :many
SELECT * FROM destinos ORDER BY fecha_viaje DESC;

-- name: UpdateDestino :exec
UPDATE destinos 
SET nombre_lugar = $2, fecha_viaje = $3, valoracion = $4, descripcion = $5, categoria_id = $6, usuario_id = $7 
WHERE id = $1;

-- name: DeleteDestino :exec
DELETE FROM destinos WHERE id = $1;
