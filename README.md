# TP5: Renderizado del Lado del Servidor con Templ

Aplicacion web desarrollada en Go que implementa un sistema completo de gestion de viajes, usuarios y categorias, conectandose a una base de datos PostgreSQL mediante consultas tipadas con sqlc, renderizado de componentes mediante Templ, automatizada con Docker y controlada mediante un Makefile.

## Tecnologias Utilizadas

* Lenguaje: Go (Golang) (`net/http`)

* Motor de Vistas: Templ (componentes UI fuertemente tipados)

* Base de Datos: PostgreSQL

* Generador SQL: `sqlc`

* Contenedorizacion: Docker

* Automatizacion: GNU Make (`Makefile`)

## Prerrequisitos

Asegurese de tener instalado en su sistema operativo:

* Go (version 1.18 o superior)

* Docker (con el servicio activo)

* Make

## Configuracion del Entorno (.env)

El proyecto utiliza un archivo de configuracion de entorno (`.env`). No es necesario crearlo manualmente; el Makefile se encarga de autogenerarlo a partir de la plantilla provista (`.env.example`) la primera vez que se ejecute el comando de construccion o ejecucion.

## Instrucciones de Ejecucion y Pruebas

Para clonar la rama especifica del proyecto, ejecute en la terminal el siguiente comando:

```
git clone --branch tp5 --single-branch https://github.com/keeganian/TP-programacion-web.git

```

Una vez dentro de la carpeta del proyecto, para levantar el entorno completo de base de datos en Docker, generar el codigo de sqlc y templ, compilar la aplicacion e iniciar el servidor web de manera totalmente automatizada, utilice el comando:

```
make run-app

```

Con este comando, el servidor quedara corriendo y accesible en `http://localhost:8080`. Al presionar `Ctrl+C` en la terminal (o al finalizar la ejecucion), el sistema realizara de manera automatizada la limpieza y eliminacion de los contenedores y volumenes de Docker.

Una vez corriendo la app y dentro de `http://localhost:8080` para poder acceder a todas las funciones es necesario crear un usuario, y luego iniciar sesion.