Trabajo Práctico 2 - Persistencia y Testing
Proyecto backend en Go con persistencia en PostgreSQL utilizando sqlc y un ciclo automatizado de tests de integración con Docker mediante Makefile.

---

Requisitos del sistema
Para clonar y ejecutar este proyecto se debe contar previamente con:

Sistema operativo: Linux o entorno tipo Unix compatible con Bash.
Go: Versión 1.21 o superior instalada en el sistema.
Docker Engine: Servicio instalado y en ejecución .
GNU Make: Herramienta instalada.
sqlc : herramienta instalada.

---

Instrucciones de ejecución
Todo el ciclo de vida de los tests (levantar la base, cargar esquema, probar y limpiar) se ejecuta con un único comando:

Clonar el repositorio y situarse en la branch de entrega
bash:
git clone https://github.com/keeganian/TP-programacion-web.git
cd TP-programacion-web
git checkout tp2

Verificación del resultado
Al ejecutar make test, la salida en la terminal debe mostrar secuencialmente las siguientes etapas sin arrojar códigos de error:

Fase de inicio y carga de esquema:

La terminal espera la disponibilidad del motor y muestra:

/var/run/postgresql:5432 - accepting connections
Cargando esquema en la base de datos...
CREATE TABLE
Fase de ejecución de tests:

Se inicia la prueba unitaria: === RUN   TestDestinosCRUD.

Concluye de forma satisfactoria marcando:

--- PASS: TestDestinosCRUD (xs)
PASS
ok  WebEntregable/db/sqlc  xs
Fase de limpieza posterior:

Se detiene y borra el contenedor temporal junto con sus volúmenes:

Borrando contenedor y volúmenes de prueba...
tp2-postgres-test
Criterio de validación: El resultado es correcto si la prueba finaliza con PASS, el paquete reporta estado ok, y la terminal retorna al prompt sin advertencias de fallo ni contenedores residuales activos.
