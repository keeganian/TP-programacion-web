package main

import (
	"log"
	"net/http"

	"TP-programacion-web/db"
	dbSqlc "TP-programacion-web/db/sqlc"
	"TP-programacion-web/handlers"
)

func main() {
	// 1. Conexión a la base de datos usando db.ConnectDB()
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer conn.Close()

	// 2. Inicializar consultas e instanciar el ServerHandler
	queries := dbSqlc.New(conn)
	serverHandler := handlers.NewServerHandler(queries)

	// 3. Archivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Registrar rutas HTTP redirigidas a los métodos del Handler
	http.HandleFunc("/", serverHandler.HandleIndex)

	// Rutas de Sesión (Login / Logout)
	http.HandleFunc("/login", serverHandler.HandleLogin)
	http.HandleFunc("/logout", serverHandler.HandleLogout)

	// Rutas de Destinos
	http.HandleFunc("/destinos", serverHandler.HandleCreateDestino)
	http.HandleFunc("/destinos/delete", serverHandler.HandleDeleteDestino)

	// Rutas de Usuarios
	http.HandleFunc("/usuarios", serverHandler.HandleCreateUsuario)
	http.HandleFunc("/usuarios/delete", serverHandler.HandleDeleteUsuario)

	// Rutas de Categorías
	http.HandleFunc("/categorias", serverHandler.HandleCreateCategoria)
	http.HandleFunc("/categorias/delete", serverHandler.HandleDeleteCategoria)

	log.Println("Servidor corriendo en el puerto 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
