#!/bin/bash

echo "=== 1. CREAR UN USUARIO (POST) ==="
response_user=$(curl -s -X POST http://localhost:8080/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Ian Keegan",
    "email": "ian@example.com"
  }')
echo "$response_user"
echo -e "\n"

echo "=== 2. CREAR UNA CATEGORÍA (POST) ==="
response_cat=$(curl -s -X POST http://localhost:8080/categorias \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Vacaciones"
  }')
echo "$response_cat"
echo -e "\n"

echo "=== 3. CREAR UN DESTINO ASOCIADO (POST) ==="
response_dest=$(curl -s -X POST http://localhost:8080/destinos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre_lugar": "Bariloche",
    "fecha_viaje": "2026-12-01T00:00:00Z",
    "valoracion": 10,
    "descripcion": "Viaje de egresados o vacaciones",
    "categoria_id": 1,
    "usuario_id": 1
  }')
echo "$response_dest"
echo -e "\n"

echo "=== 4. LISTAR TODOS LOS DESTINOS (GET) ==="
curl -s -X GET http://localhost:8080/destinos
echo -e "\n\n"

echo "=== 5. OBTENER DESTINO POR ID 1 (GET) ==="
curl -s -X GET http://localhost:8080/destinos/1
echo -e "\n\n"

echo "=== 6. ACTUALIZAR DESTINO 1 (PUT) ==="
curl -s -X PUT http://localhost:8080/destinos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nombre_lugar": "Bariloche Actualizado",
    "fecha_viaje": "2026-12-15T00:00:00Z",
    "valoracion": 9,
    "descripcion": "Modificado correctamente",
    "categoria_id": 1,
    "usuario_id": 1
  }'
echo -e "\n\n"

echo "=== 7. ELIMINAR DESTINO 1 (DELETE) ==="
curl -s -o /dev/null -w "Código de estado HTTP: %{http_code}\n" -X DELETE http://localhost:8080/destinos/1
