# ETAPA 1: BUILD (Node.js)
FROM node:22-alpine AS build
WORKDIR /app

# Igual que en Go, copiamos dependencias primero para optimizar el caché.
COPY package*.json ./
# Usamos 'ci' en lugar de 'install' para respetar estrictamente las versiones del package-lock.json.
RUN npm ci

# Copiamos el resto del código y generamos los estáticos de Vite.
COPY . .
RUN npm run build

# ETAPA TEST: EJECUCIÓN DE PRUEBAS UNITARIAS Y COBERTURA (TP5)
FROM build AS test
WORKDIR /app
ENTRYPOINT ["npm", "run", "test:ci"]

# ETAPA 2: SERVIDOR (Nginx)
FROM nginx:alpine

# Copiamos tu plantilla de configuración para que Nginx procese variables de entorno al iniciar
COPY default.conf.template /etc/nginx/templates/default.conf.template

# Copiamos los archivos HTML/JS/CSS generados en la carpeta /dist de la etapa anterior.
COPY --from=build /app/dist /usr/share/nginx/html

ENV BACKEND_URL=http://backend:8080 DNS_RESOLVER=127.0.0.11

EXPOSE 80
