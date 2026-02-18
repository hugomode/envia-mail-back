# envia-mail-back

API en Go para envio de correos via SMTP con Gin y documentacion Swagger.

## Requisitos

- Go 1.22+
- Docker (opcional)
- Acceso al servidor SMTP configurado

## Variables de entorno

- `SERVER`: host del servidor SMTP
- `SMTP_PORT`: puerto SMTP (para Gmail usar `587`)
- `SMTP_USERNAME`: usuario SMTP
- `SMTP_PASSWORD`: password SMTP (para Gmail usar App Password)
- `FROM`: remitente por defecto cuando `from` no termina en `@uchile.cl`
- `TIMEOUT_SECONDS`: timeout de envio SMTP en segundos
- `HOST` (opcional): host publicado en Swagger (default: `localhost:8080`)

Para Gmail:
- `SERVER=smtp.gmail.com`
- `SMTP_PORT=587`
- `SMTP_USERNAME` y `FROM` deben ser la misma cuenta Gmail
- `SMTP_PASSWORD` debe ser una App Password de Google (no la clave normal)

## Ejecucion local

```bash
go mod tidy
go run ./cmd/main.go
```

Tambien puedes usar Make:

```bash
make run
```

## Calidad y utilidades

```bash
make fmt
make tidy
make test
```

## Docker

Build de imagen de desarrollo:

```bash
make docker-build
```

Shell dentro del contenedor de desarrollo:

```bash
make docker-shell
```

## Endpoints

- `GET /api/v1/healthcheck`
- `POST /api/v1/notificacion`
- `GET /api/v1/swagger/index.html`

## Formato esperado para `POST /api/v1/notificacion`

Enviar `application/json` en el body.
Si incluyes adjuntos, `archivo` debe ir en base64.

```json
{
  "asunto": "Asunto del correo",
  "contenido": "<p>Contenido HTML</p>",  
  "to": ["destino@uchile.cl"],
  "cc": ["copia@uchile.cl"],
  "bcc": ["oculto@uchile.cl"],
  "adjuntos": [
    {
      "archivo": "JVBERi0xLjQKJcfsj6IK...", 
      "nombre_archivo": "reporte.pdf"
    }
  ]
}
```
