# FranchiseAPI

Es una aplicación de servidor que proporciona endpoints para gestionar información sobre franquicias, compañías y detalles de dominio. La API utiliza el framework [Gin](https://gin-gonic.com/) para las rutas y [GORM](https://gorm.io/) como ORM para interactuar con la base de datos.

## Configuración del Proyecto

### Requisitos Previos

Asegúrate de tener Go instalado en tu sistema. Puedes descargar Go desde [https://golang.org/dl/](https://golang.org/dl/).

### Instalación de Dependencias

Ejecuta el siguiente comando para instalar las dependencias del proyecto:

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u github.com/likexian/whois
```

### Configuración de la Base de Datos

La aplicación utiliza GORM como ORM y actualmente está configurada para conectarse a una base de datos PostgreSQL. Asegúrate de tener un servidor PostgreSQL en ejecución y actualiza la configuración de la base de datos en `config/config.go` según tus necesidades.

## Ejecución del Proyecto

Ejecuta el siguiente comando para iniciar el servidor:

```bash
go run main.go
```

El servidor se iniciará en el puerto 3000 por defecto. Puedes cambiar el puerto en el archivo `main.go` si es necesario.

## Rutas de la API

1. **POST /franchise**: Crea una nueva franquicia.
2. **PUT /franchise/:id**: Actualiza una franquicia existente por ID.
3. **GET /franchise/:id**: Obtiene información sobre una franquicia por ID.
4. **GET /franchises**: Obtiene todas las franquicias.
5. **GET /companies/:id**: Obtiene todas las franquicias asociadas a una compañía por ID.
6. **GET /location/:country**: Obtiene todas las franquicias en un país específico.
7. **GET /ssl-info**: Obtiene información SSL para un host específico.
8. **GET /domain-info**: Obtiene información de dominio para un dominio específico.

## Estructura de Datos

La aplicación gestiona datos sobre franquicias, compañías, ubicaciones, propietarios, información de dominio y detalles de información. La estructura de datos está definida en el archivo `main.go`.

## Servicios Externos

La API utiliza el servicio de [SSL Labs](https://www.ssllabs.com/) para obtener información SSL y la biblioteca [likexian/whois](https://github.com/likexian/whois) para obtener detalles de dominio.
