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

# ETAPA 2: SERVIDOR (Nginx)
FROM nginx:alpine

# Copiamos tu archivo nginx.conf (que deberá tener el proxy_pass hacia http://backend:8080/api/)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiamos los archivos HTML/JS/CSS generados en la carpeta /dist de la etapa anterior.
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80