package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const connStr = "postgres://nosotros:1234@localhost:5433/tpespecial?sslmode=disable"

func setupTestDB(t *testing.T) (*Queries, *sql.DB) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("Error de conexión (ping): %v", err)
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return New(conn), conn
}

func TestDestinosCRUD(t *testing.T) {
	ctx := context.Background()
	q, _ := setupTestDB(t)

	// 1. CreateDestino
	nuevo, err := q.CreateDestino(ctx, CreateDestinoParams{
		NombreLugar: "Bariloche",
		FechaViaje:  time.Now(),
		Valoracion:  sql.NullInt32{Int32: 5, Valid: true},
		Descripcion: sql.NullString{String: "Viaje a las montañas", Valid: true},
	})
	if err != nil {
		t.Fatalf("Error en CreateDestino: %v", err)
	}
	if nuevo.ID == 0 {
		t.Errorf("ID inválido, esperado > 0")
	}

	// 2. GetDestino
	obtenido, err := q.GetDestino(ctx, nuevo.ID)
	if err != nil {
		t.Fatalf("Error en GetDestino: %v", err)
	}
	if obtenido.NombreLugar != nuevo.NombreLugar {
		t.Errorf("Esperado %s, obtenido %s", nuevo.NombreLugar, obtenido.NombreLugar)
	}

	// 3. ListDestinos
	lista, err := q.ListDestinos(ctx)
	if err != nil {
		t.Fatalf("Error en ListDestinos: %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Se esperaba al menos un elemento en la lista")
	}

	// 4. UpdateDestino
	err = q.UpdateDestino(ctx, UpdateDestinoParams{
		ID:          nuevo.ID,
		NombreLugar: "Bariloche Modificado",
		FechaViaje:  time.Now(),
		Valoracion:  sql.NullInt32{Int32: 4, Valid: true},
		Descripcion: sql.NullString{String: "Nueva descripción", Valid: true},
	})
	if err != nil {
		t.Fatalf("Error en UpdateDestino: %v", err)
	}

	// 5. DeleteDestino
	err = q.DeleteDestino(ctx, nuevo.ID)
	if err != nil {
		t.Fatalf("Error en DeleteDestino: %v", err)
	}

	// Verificar eliminación
	_, err = q.GetDestino(ctx, nuevo.ID)
	if err == nil {
		t.Errorf("El registro sigue existiendo tras DeleteDestino")
	}
}
