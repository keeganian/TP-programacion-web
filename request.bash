#!/bin/bash

# Terminar el script inmediatamente si un comando falla
set -e

URL="http://localhost:8080"

echo "=== INICIANDO VALIDACIÓN AUTOMÁTICA DE API (CRUD) ==="

# Función auxiliar para validar respuestas HTTP
check_status() {
    local expected="$1"
    local actual="$2"
    local step="$3"

    if [ "$actual" -eq "$expected" ]; then
        echo "✅ OK: $step (Status HTTP $actual)"
    else
        echo "❌ ERROR: $step. Esperado $expected, pero se obtuvo $actual."
        exit 1
    fi
}

echo "=== 1. CREAR UN USUARIO (POST) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/usuarios" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Ian Keegan",
    "email": "ian@example.com"
  }')
check_status 201 "$status" "POST /usuarios"

echo "=== 2. CREAR UNA CATEGORÍA (POST) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/categorias" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Vacaciones"
  }')
check_status 201 "$status" "POST /categorias"

echo "=== 3. CREAR UN DESTINO (POST) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/destinos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre_lugar": "Bariloche",
    "fecha_viaje": "2026-12-01T00:00:00Z",
    "valoracion": 10,
    "descripcion": "Viaje de egresados o vacaciones",
    "categoria_id": 1,
    "usuario_id": 1
  }')
check_status 201 "$status" "POST /destinos"

echo "=== 4. LISTAR TODOS LOS DESTINOS (GET) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$URL/destinos")
check_status 200 "$status" "GET /destinos"

echo "=== 5. OBTENER DESTINO EXISTENTE POR ID (GET) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$URL/destinos/1")
check_status 200 "$status" "GET /destinos/1"

echo "=== 6. OBTENER DESTINO INEXISTENTE (GET - Debe fallar 404) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$URL/destinos/9999")
check_status 404 "$status" "GET /destinos/9999 (Inexistente)"

echo "=== 7. ACTUALIZAR DESTINO EXISTENTE (PUT) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$URL/destinos/1" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre_lugar": "Bariloche Actualizado",
    "fecha_viaje": "2026-12-15T00:00:00Z",
    "valoracion": 9,
    "descripcion": "Modificado correctamente",
    "categoria_id": 1,
    "usuario_id": 1
  }')
check_status 200 "$status" "PUT /destinos/1"

echo "=== 8. ELIMINAR DESTINO (DELETE) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$URL/destinos/1")
# Acepta 200 o 204 según la implementación
if [ "$status" -eq 200 ] || [ "$status" -eq 204 ]; then
    echo "✅ OK: DELETE /destinos/1 (Status HTTP $status)"
else
    echo "❌ ERROR: DELETE /destinos/1. Esperado 200 o 204, obtenido $status"
    exit 1
fi

echo "=== 9. VERIFICAR ELIMINACIÓN DE DESTINO (GET - Debe dar 404) ==="
status=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$URL/destinos/1")
check_status 404 "$status" "GET /destinos/1 tras ser eliminado"

echo "🎉 === TODOS LOS TESTEOS PASARON EXITOSAMENTE ==="
