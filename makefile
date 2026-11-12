APP_NAME := my-app

# Carga automática del archivo .env si existe
-include .env
export

.PHONY: all generate build db-up db-down run-app clean

# Asegura que el archivo .env existe a partir de .env.example
.env:
	@if [ ! -f .env ]; then \
		echo "Creando archivo .env desde .env.example..."; \
		cp .env.example .env; \
	fi

# 1. Tareas previas: Generar código SQL y plantillas Templ
generate: .env
	@echo "Generando código con sqlc y templ..."
	@sqlc generate
	@go run github.com/a-h/templ/cmd/templ@latest generate

# 2. Compilar aplicación Go
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# 3. Tareas previas: Borrar contenedor previo, levantar Postgres, esperar y aplicar schema.sql
db-up: .env
	@echo "Limpiando contenedores y volúmenes anteriores..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true
	@echo "Levantando contenedor de base de datos..."
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

# 4. Tareas posteriores: Borrar contenedor y volúmenes
db-down:
	@echo "Borrando contenedor y volúmenes..."
	@docker rm -f -v $(DB_CONTAINER) 2>/dev/null || true

# 5. Tarea principal exigida para probar la entrega
run-app: db-up build
	@echo "Iniciando la aplicación en http://localhost:8080..."
	@./tmp/$(APP_NAME) || $(MAKE) db-down
	@$(MAKE) db-down

clean: db-down
	@rm -rf tmp
