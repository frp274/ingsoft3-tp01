# 🎟️ PlanesApp - Plataforma de Organización y Descubrimiento de Planes

MVP de una plataforma colaborativa para descubrir eventos, crear planes públicos y gestionar solicitudes de invitación entre usuarios de una comunidad.

---

## 🏗️ Arquitectura y Tecnologías

- **Backend:** Go 1.22 con Clean Architecture (Handlers, Services, Repositories) e inyección de dependencias (`*sql.DB`).
- **Base de Datos:** MySQL 8 con persistencia en volumen nombrado (`db_data`), migraciones DDL automáticas y script Seed inicial.
- **Frontend:** React + Vite, interfaz reactiva con tarjetas interactivas y formularios controlados.
- **Proxy Inverso & Servidor Web:** Nginx (Alpine) con resolución DNS dinámica interna y fallback SPA.
- **Contenedores & Orquestación:** Docker Multi-stage Builds, Docker Compose y distribución vía GitHub Packages Container Registry (GHCR).
- **CI / CD:** Pipeline as Code con GitHub Actions (`ci.yml`) y caché optimizado de capas Docker (`mode=max, scope=backend/frontend`).

---

## 🚀 Guía de Instalación y Despliegue Rápido

Sigue estos pasos para levantar la plataforma completa desde cero:

### 1. Clonar el repositorio
```bash
git clone https://github.com/frp274/ingsoft3-tp01.git
cd ingsoft3-tp01
```

### 2. Configurar variables de entorno
Crea tu archivo `.env` local a partir de la plantilla:
```bash
cp .env.example .env
```
*(En Windows PowerShell puedes usar `Copy-Item .env.example .env`)*

### 3. Levantar los contenedores

#### Opción A: Construcción y ejecución local (Recomendado para desarrollo)
```bash
docker compose up -d --build
```

#### Opción B: Ejecución desde el Container Registry (GHCR)
Si deseas ejecutar directamente las imágenes precompiladas y publicadas en GitHub Container Registry:
```bash
docker compose -f docker-compose.registry.yml up -d
```

---

## 🌐 Puntos de Acceso

Una vez levantados los servicios, accede a través de tu navegador:

| Servicio | URL | Descripción |
| :--- | :--- | :--- |
| **Frontend Web** | [http://localhost:3000](http://localhost:3000) | Aplicación React (Novedades, Mis Próximos Planes, Crear Plan) |
| **Nginx API Proxy** | [http://localhost:3000/api/health](http://localhost:3000/api/health) | Acceso al backend mediante proxy inverso (mismo origen) |
| **Backend API Directo** | [http://localhost:8080/api/health](http://localhost:8080/api/health) | Endpoint directo del servidor Go |
| **Base de Datos** | `localhost:3306` | MySQL (red interna compose `db:3306`) |

---

## 📡 Endpoints de la API REST

- `GET /api/health`: Estado de salud del servidor y conexión con MySQL.
- `GET /api/planes`: Lista todos los planes públicos disponibles.
- `POST /api/planes`: Crea un nuevo plan (valida campos obligatorios y coherencia de fechas).
- `GET /api/mis-planes`: Consulta los planes donde el usuario se ha postulado (JOIN con `invitaciones`).
- `POST /api/invitaciones`: Registra una postulación a un plan (`idPlan`).
- `GET /api/etiquetas`: Lista las categorías disponibles (ej. "Outdoor", "Fiesta", "Cine").
- `GET /api/usuarios`: Consulta los usuarios registrados.

---

## 🧪 Datos Iniciales (Seed)

Al iniciar el backend por primera vez, el sistema inicializa automáticamente:
- **1 Usuario de prueba:** `juan.perez@example.com`
- **3 Etiquetas base:** "Outdoor", "Fiesta", "Cine"
- **5 Planes públicos de ejemplo** con fechas y ubicaciones listas para interactuar.