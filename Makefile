APP_NAME := my-app
DB_CONTAINER := tp2-postgres-test
DB_PORT := 5433
DB_USER := nosotros
DB_PASS := 1234
DB_NAME := tpespecial

.PHONY: all run generate build test clean db-up db-down

all: build

# 1. Generar código SQL
generate:
	@sqlc generate

# 2. Compilar aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Tarea para levantar PostgreSQL de test
db-up:
	@echo "Limpiando contenedores anteriores si existen..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true
	@echo "Levantando base de datos temporal..."
	@docker run --name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 -d postgres:latest
	@echo "Esperando que PostgreSQL esté listo..."
	@until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		sleep 1; \
	done
	@echo "Cargando esquema en la base de datos..."
	@docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < db/schema/schema.sql

# Tarea para borrar contenedor y volúmenes
db-down:
	@echo "Borrando contenedor y volúmenes de prueba..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true

# Ejecuta el flujo completo pedido:
# Previo (sqlc, build, limpiar y levantar BD, esperar BD) -> Tests -> Posterior (borrar BD y volumen)
test: generate build
	@$(MAKE) db-up
	@echo "Ejecutando tests de Go..."
	@-go test -v ./db/sqlc; \
	STATUS=$$?; \
	$(MAKE) db-down; \
	exit $$STATUS

clean:
	@rm -rf tmp
