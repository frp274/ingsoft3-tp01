# ETAPA 1: BUILD (SDK pesado)
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copiamos solo los archivos de dependencias primero. 
# Esto aprovecha el caché de capas de Docker: si no agregaste librerías nuevas, 
# Docker reutiliza esta capa y no vuelve a descargar todo de internet.
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código fuente y compilamos.
COPY . .
# Compilamos un binario estático llamado 'main'
RUN go build -o main ./cmd/api/main.go

# ETAPA 2: RUNTIME (Imagen mínima)
FROM alpine:latest AS final
WORKDIR /app

# Traemos ÚNICAMENTE el binario compilado desde la etapa 'build'.
# Así, el compilador de Go no viaja a producción, reduciendo el peso de gigas a megas.
COPY --from=build /app/main .

EXPOSE 8080
# Definimos el ejecutable del contenedor.
ENTRYPOINT ["./main"]