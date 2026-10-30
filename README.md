# TP3: Servidor Web en Go con PostgreSQL y sqlc

API RESTful desarrollada en Go que implementa un sistema de gestion de viajes, conectandose a una base de datos PostgreSQL mediante consultas tipadas con sqlc, automatizada con Docker y controlada mediante un Makefile.

## Tecnologias Utilizadas

* Lenguaje: Go (Golang) (net/http)

* Base de Datos: PostgreSQL

* Generador SQL: sqlc

* Contenedorizacion: Docker

* Automatizacion: GNU Make (Makefile)

## Prerrequisitos

Asegurese de tener instalado en su sistema operativo:

* Go (version 1.18 o superior)

* Docker (con el servicio activo)

* Make

## Configuracion del Entorno (.env)

El proyecto utiliza un archivo de configuracion de entorno (.env). No es necesario crearlo manualmente; el Makefile se encarga de autogenerarlo a partir de la plantilla provista (.env.example) la primera vez que se ejecute el comando de construccion o pruebas.

## Instrucciones de Ejecucion y Pruebas

Para clonar el branch, hay que ejecutar en la terminal el comando: 

```
git clone --branch tp3 --single-branch https://github.com/keeganian/TP-programacion-web.git
```
Para levantar el entorno completo, compilar la aplicacion, iniciar el servidor de manera automatizada y gestionar la limpieza de contenedores, utilice el siguiente comando principal:

```
sudo make test-api
```


Con este comando, el servidor queda corriendo en la terminal y ejecuta las pruebas CRUD automáticamente, utilizando ./request.bash, indicando si pasa el testeo o no. Finalizada la ejecución, cierra todos los docker.

