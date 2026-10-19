package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	db "TP-programacion-web/db/sqlc"

	_ "github.com/lib/pq"
)

type Server struct {
	queries *db.Queries
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5433"
	}

	connStr := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("No se puede hacer ping a la base de datos: %v", err)
	}

	server := &Server{
		queries: db.New(database),
	}

	// Rutas de Destinos
	http.HandleFunc("/destinos", server.handleDestinosCollection)
	http.HandleFunc("/destinos/", server.handleDestinosItem)

	// Rutas de Usuarios
	http.HandleFunc("/usuarios", server.handleUsuariosCollection)
	http.HandleFunc("/usuarios/", server.handleUsuariosItem)

	// Rutas de Categorías
	http.HandleFunc("/categorias", server.handleCategoriasCollection)
	http.HandleFunc("/categorias/", server.handleCategoriasItem)

	log.Println("Servidor corriendo en el puerto 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

// ==========================================
// MANEJADORES DE DESTINOS
// ==========================================

func (s *Server) handleDestinosCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listDestinos(w, r)
	case http.MethodPost:
		s.createDestino(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDestinosItem(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/destinos/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getDestino(w, r, int32(id))
	case http.MethodPut:
		s.updateDestino(w, r, int32(id))
	case http.MethodDelete:
		s.deleteDestino(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listDestinos(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	destinos, err := s.queries.ListDestinos(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destinos)
}

func (s *Server) createDestino(w http.ResponseWriter, r *http.Request) {
	var input struct {
		NombreLugar string    `json:"nombre_lugar"`
		FechaViaje  time.Time `json:"fecha_viaje"`
		Valoracion  int32     `json:"valoracion"`
		Descripcion string    `json:"descripcion"`
		CategoriaID int32     `json:"categoria_id"`
		UsuarioID   int32     `json:"usuario_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.NombreLugar) == "" || input.Valoracion < 1 || input.Valoracion > 10 || input.UsuarioID == 0 {
		http.Error(w, "Datos de validación incorrectos", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	created, err := s.queries.CreateDestino(ctx, db.CreateDestinoParams{
		NombreLugar: input.NombreLugar,
		FechaViaje:  input.FechaViaje,
		Valoracion:  sql.NullInt32{Int32: input.Valoracion, Valid: true},
		Descripcion: sql.NullString{String: input.Descripcion, Valid: input.Descripcion != ""},
		CategoriaID: sql.NullInt32{Int32: input.CategoriaID, Valid: input.CategoriaID != 0},
		UsuarioID:   sql.NullInt32{Int32: input.UsuarioID, Valid: true},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (s *Server) getDestino(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	destino, err := s.queries.GetDestino(ctx, id)
	if err != nil {
		http.Error(w, "Destino no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destino)
}

func (s *Server) updateDestino(w http.ResponseWriter, r *http.Request, id int32) {
	var input struct {
		NombreLugar string    `json:"nombre_lugar"`
		FechaViaje  time.Time `json:"fecha_viaje"`
		Valoracion  int32     `json:"valoracion"`
		Descripcion string    `json:"descripcion"`
		CategoriaID int32     `json:"categoria_id"`
		UsuarioID   int32     `json:"usuario_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := s.queries.UpdateDestino(ctx, db.UpdateDestinoParams{
		ID:          id,
		NombreLugar: input.NombreLugar,
		FechaViaje:  input.FechaViaje,
		Valoracion:  sql.NullInt32{Int32: input.Valoracion, Valid: true},
		Descripcion: sql.NullString{String: input.Descripcion, Valid: input.Descripcion != ""},
		CategoriaID: sql.NullInt32{Int32: input.CategoriaID, Valid: input.CategoriaID != 0},
		UsuarioID:   sql.NullInt32{Int32: input.UsuarioID, Valid: true},
	})
	if err != nil {
		http.Error(w, "No se pudo actualizar el destino", http.StatusNotFound)
		return
	}

	updated, err := s.queries.GetDestino(ctx, id)
	if err != nil {
		http.Error(w, "Destino actualizado pero no se pudo recuperar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (s *Server) deleteDestino(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	err := s.queries.DeleteDestino(ctx, id)
	if err != nil {
		http.Error(w, "Destino no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// MANEJADORES DE USUARIOS
// ==========================================

func (s *Server) handleUsuariosCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listUsuarios(w, r)
	case http.MethodPost:
		s.createUsuario(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUsuariosItem(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/usuarios/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getUsuario(w, r, int32(id))
	case http.MethodPut:
		s.updateUsuario(w, r, int32(id))
	case http.MethodDelete:
		s.deleteUsuario(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listUsuarios(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	usuarios, err := s.queries.ListUsuarios(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func (s *Server) createUsuario(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" || strings.TrimSpace(input.Email) == "" {
		http.Error(w, "Datos de validación incorrectos", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	created, err := s.queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Nombre: input.Nombre,
		Email:  input.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (s *Server) getUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	usuario, err := s.queries.GetUsuario(ctx, id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func (s *Server) updateUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	var input struct {
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := s.queries.UpdateUsuario(ctx, db.UpdateUsuarioParams{
		ID:     id,
		Nombre: input.Nombre,
		Email:  input.Email,
	})
	if err != nil {
		http.Error(w, "No se pudo actualizar el usuario", http.StatusNotFound)
		return
	}

	updated, err := s.queries.GetUsuario(ctx, id)
	if err != nil {
		http.Error(w, "Usuario actualizado pero no se pudo recuperar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (s *Server) deleteUsuario(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	err := s.queries.DeleteUsuario(ctx, id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// MANEJADORES DE CATEGORÍAS
// ==========================================

func (s *Server) handleCategoriasCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listCategorias(w, r)
	case http.MethodPost:
		s.createCategoria(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCategoriasItem(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/categorias/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getCategoria(w, r, int32(id))
	case http.MethodPut:
		s.updateCategoria(w, r, int32(id))
	case http.MethodDelete:
		s.deleteCategoria(w, r, int32(id))
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listCategorias(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	categorias, err := s.queries.ListCategorias(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func (s *Server) createCategoria(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" {
		http.Error(w, "Datos de validación incorrectos", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	created, err := s.queries.CreateCategoria(ctx, input.Nombre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (s *Server) getCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	categoria, err := s.queries.GetCategoria(ctx, id)
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func (s *Server) updateCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	var input struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := s.queries.UpdateCategoria(ctx, db.UpdateCategoriaParams{
		ID:     id,
		Nombre: input.Nombre,
	})
	if err != nil {
		http.Error(w, "No se pudo actualizar la categoría", http.StatusNotFound)
		return
	}

	updated, err := s.queries.GetCategoria(ctx, id)
	if err != nil {
		http.Error(w, "Categoría actualizada pero no se pudo recuperar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (s *Server) deleteCategoria(w http.ResponseWriter, r *http.Request, id int32) {
	ctx := context.Background()
	err := s.queries.DeleteCategoria(ctx, id)
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
