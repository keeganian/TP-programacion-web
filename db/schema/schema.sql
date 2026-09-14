CREATE TABLE destinos (
    id SERIAL PRIMARY KEY,
    nombre_lugar VARCHAR(255) NOT NULL,
    fecha_viaje DATE NOT NULL,
    valoracion INT CHECK (valoracion >= 1 AND valoracion <= 10),
    descripcion TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);