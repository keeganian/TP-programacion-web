APP_NAME := my-app

# Carga automática del archivo .env generado
-include .env
export

.PHONY: all run generate build test-api clean db-up db-down

all: build

# Regla para asegurar que el archivo .env existe a partir de .env.example
.env:
	@if [ ! -f .env ]; then \
		echo "Creando archivo .env desde .env.example..."; \
		cp .env.example .env; \
	fi

# 1. Generar código SQL
generate: .env
	@sqlc generate

# 2. Compilar aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Tarea para levantar PostgreSQL usando las credenciales del .env
db-up: .env
	@echo "Limpiando contenedores anteriores si existen..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true
	@echo "Levantando base de datos temporal..."
	@docker run --name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 -d postgres:latest
	@echo "Esperando que PostgreSQL esté listo..."
	@until docker exec $(DB_CONTAINER) pg_isready -h localhost -U $(DB_USER) -d $(DB_NAME); do \
		sleep 1; \
	done
	@echo "Cargando esquema en la base de datos..."
	@docker exec -i $(DB_CONTAINER) psql -h localhost -U $(DB_USER) -d $(DB_NAME) < db/schema/schema.sql

# Tarea para borrar contenedor y volúmenes
db-down:
	@echo "Borrando contenedor y volúmenes de prueba..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true

# Ejecuta la base de datos, el servidor en segundo plano, corre las validaciones y apaga todo
test-api: db-up build
	@echo "Iniciando servidor para tests..."
	@./tmp/$(APP_NAME) & SERVER_PID=$$!; \
	sleep 2; \
	chmod +x request.bash; \
	./request.bash; \
	TEST_RESULT=$$?; \
	kill $$SERVER_PID 2>/dev/null || true; \
	$(MAKE) db-down; \
	exit $$TEST_RESULT

clean:
	@rm -rf tmp
