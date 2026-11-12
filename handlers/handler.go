package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "TP-programacion-web/db/sqlc"
	"TP-programacion-web/views"
)

type ServerHandler struct {
	queries *db.Queries
}

func NewServerHandler(q *db.Queries) *ServerHandler {
	return &ServerHandler{queries: q}
}

// ==========================================
// AUXILIAR DE AUTENTICACIÓN
// ==========================================

// getLoggedUser lee la cookie de sesión "user_id" y busca el usuario en la BD
func (h *ServerHandler) getLoggedUser(r *http.Request) *db.Usuario {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return nil
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return nil
	}

	usuario, err := h.queries.GetUsuario(r.Context(), int32(userID))
	if err != nil {
		return nil
	}

	return &usuario
}

// ==========================================
// CONTROLADORES DE SESIÓN (LOGIN / LOGOUT)
// ==========================================

func (h *ServerHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	usuario, err := h.queries.GetUsuarioByEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "/?error=user_not_found", http.StatusSeeOther)
		return
	}

	// Guardar Cookie de sesión por 24 horas
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(usuario.ID)),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ServerHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ==========================================
// VISTA PRINCIPAL (GET /)
// ==========================================

func (h *ServerHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	loggedUser := h.getLoggedUser(r)
	filtroLugar := r.URL.Query().Get("lugar")

	var destinos []db.Destino
	var err error

	if filtroLugar != "" {
		destinos, err = h.queries.SearchDestinosByLugar(ctx, sql.NullString{String: filtroLugar, Valid: filtroLugar != ""})
	} else {
		destinos, err = h.queries.ListDestinos(ctx)
	}

	if err != nil {
		log.Printf("Error al consultar destinos: %v", err)
		destinos = []db.Destino{}
	}

	var misResenas []db.Destino
	if loggedUser != nil {
		for _, d := range destinos {
			if d.UsuarioID.Valid && d.UsuarioID.Int32 == loggedUser.ID {
				misResenas = append(misResenas, d)
			}
		}
	}

	usuarios, _ := h.queries.ListUsuarios(ctx)
	categorias, _ := h.queries.ListCategorias(ctx)

	err = views.Layout("Destinados a no laburar", destinos, misResenas, usuarios, categorias, loggedUser, filtroLugar).Render(ctx, w)
	if err != nil {
		log.Printf("Error al renderizar plantilla: %v", err)
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
		return
	}
}

// ==========================================
// DESTINOS
// ==========================================

func (h *ServerHandler) HandleCreateDestino(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	loggedUser := h.getLoggedUser(r)

	nombreLugar := r.FormValue("nombre_lugar")
	fechaStr := r.FormValue("fecha_viaje")
	valoracionStr := r.FormValue("valoracion")
	descripcion := r.FormValue("descripcion")
	categoriaIDStr := r.FormValue("categoria_id")

	// Determinar ID de usuario: priorizar el logueado, sino usar el del form si existe
	var usuarioID int
	if loggedUser != nil {
		usuarioID = int(loggedUser.ID)
	} else {
		usuarioID, _ = strconv.Atoi(r.FormValue("usuario_id"))
	}

	fechaViaje, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		fechaViaje = time.Now()
	}

	valoracion, _ := strconv.Atoi(valoracionStr)
	categoriaID, _ := strconv.Atoi(categoriaIDStr)

	_, err = h.queries.CreateDestino(r.Context(), db.CreateDestinoParams{
		NombreLugar: nombreLugar,
		FechaViaje:  fechaViaje,
		Valoracion:  sql.NullInt32{Int32: int32(valoracion), Valid: valoracion > 0},
		Descripcion: sql.NullString{String: descripcion, Valid: descripcion != ""},
		CategoriaID: sql.NullInt32{Int32: int32(categoriaID), Valid: categoriaID > 0},
		UsuarioID:   sql.NullInt32{Int32: int32(usuarioID), Valid: usuarioID > 0},
	})
	if err != nil {
		log.Printf("Error al crear destino: %v", err)
		http.Error(w, "Error al crear el destino", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ServerHandler) HandleDeleteDestino(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteDestino(r.Context(), int32(id)); err != nil {
		log.Printf("Error al eliminar destino: %v", err)
		http.Error(w, "Error al eliminar destino", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ==========================================
// USUARIOS
// ==========================================

func (h *ServerHandler) HandleCreateUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	email := r.FormValue("email")

	_, err := h.queries.CreateUsuario(r.Context(), db.CreateUsuarioParams{
		Nombre: nombre,
		Email:  email,
	})
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ServerHandler) HandleDeleteUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteUsuario(r.Context(), int32(id)); err != nil {
		log.Printf("Error al eliminar usuario: %v", err)
		http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ==========================================
// CATEGORÍAS
// ==========================================

func (h *ServerHandler) HandleCreateCategoria(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")

	_, err := h.queries.CreateCategoria(r.Context(), nombre)
	if err != nil {
		log.Printf("Error al crear categoria: %v", err)
		http.Error(w, "Error al crear la categoría", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ServerHandler) HandleDeleteCategoria(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteCategoria(r.Context(), int32(id)); err != nil {
		log.Printf("Error al eliminar categoria: %v", err)
		http.Error(w, "Error al eliminar categoría", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
