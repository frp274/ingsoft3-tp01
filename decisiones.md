# Decisiones - TP1

## 1. Por qué Git no pudo resolver el conflicto solo — y qué habría tenido que pasar para que nunca apareciera.

Por que las dos ramas partieron de la misma (Main) y no "sabian" que se estaban pisando al mergear, y git necesita saber con que cambio final queda en main.
Para que nunca apareciera, deberia evitar que 2 ramas trabajen en los mismos archivos, y de ser asi, esperar que una mergee a main antes de la otra cree una nueva rama con cambios.

## 2. Qué problemas encontraste y cómo los solucionaste. Los tropiezos bien contados valen más que un camino perfecto: son los que demuestran que entendiste.

El bloqueo del push directo, solucion: solicitar y confirmar un pull request.
Un conflicto entre 2 ramas que mergean a main, solucion: resolver conflicto eligiendo alguno de los dos cambios.

## 3. Declaración de uso de IA: qué partes hiciste con ayuda de inteligencia artificial y cómo verificaste lo que te devolvió ( Uso de IA del enunciado).

La unica parte donde usea ia, predeterminada de github, fue donde te escribe el titulo y descripcion del commit, el resto nada, por que esta bien explicado.

---

# Decisiones - TP2

## 1. Qué app elegiste y por qué (contra los criterios de la guía)

Se eligió desarrollar un **"Organizador de Planes y Eventos"**. Esta aplicación cumple estrictamente con los criterios de la guía al requerir una arquitectura transaccional pura (CRUD) sobre un modelo de datos 100% relacional en MySQL.
Se descartó el uso de bases NoSQL (como MongoDB) para el MVP con el fin de priorizar la integridad referencial y las propiedades ACID al manejar concurrencia en la capacidad de los planes. El modelo abarca entidades centrales como `Usuario`, `Plan`, `Etiqueta` y una tabla intermedia `Invitacion` (M:N), utilizando claves primarias y foráneas de tipo `int` explícitamente definidas para garantizar la coherencia del esquema. El stack elegido (Go para el backend y React+Vite para el frontend) permite generar binarios compilados y estáticos muy eficientes para contenerizar.

## 2. Decisiones de contenerización

* **Estructura Multi-stage:** Se adoptó un enfoque multi-etapa en ambos Dockerfiles para optimizar radicalmente el tamaño. En el backend, la etapa `builder` utiliza la imagen pesada de `golang:1.22` para compilar, pero la imagen final se basa en una distribución mínima (o `scratch`) que solo contiene el binario. En el frontend, se usa Node para el build y la imagen final es puramente un servidor `nginx:alpine` sirviendo los estáticos.
* **Imágenes Base:** Se priorizaron imágenes `alpine` en las etapas de producción para reducir la superficie de ataque y el consumo de almacenamiento en el registro (ghcr.io).
* **Persistencia:** Se definió un *named volume* en el `docker-compose.yml` montado en `/var/lib/mysql`. Los contenedores de la aplicación (Go y Nginx) son 100% efímeros y no guardan estado; toda la persistencia de datos (usuarios, planes, relaciones) sobrevive a la destrucción de los contenedores exclusivamente gracias a este volumen de la base de datos.

## 3. Problemas encontrados y resolución

1. **Nginx fallando al inicio (Host not found):** Nginx fallaba y se cerraba al levantar el ecosistema si el contenedor del backend tardaba en estar listo, ya que intentaba resolver el host estático en el `proxy_pass`. Se solucionó definiendo el resolver interno de Docker (`127.0.0.11`) y asignando el upstream a una variable, obligando a Nginx a resolver la IP en tiempo de ejecución.
2. **Error al pushear al registry:** Al intentar publicar las imágenes en GitHub Packages, Docker arrojaba un error de permisos o de imagen inexistente (`backend:ci`). Se solucionó re-etiquetando las imágenes localmente con el formato `ghcr.io/<usuario>/<imagen>:<tag>` y generando un Personal Access Token (PAT) con permisos explícitos de `write:packages` para el `docker login`.

---

# Decisiones - TP3

## 1. Duración del Sprint

Se estableció una duración de **1 semana**. Esta métrica permite alinear el ciclo de desarrollo continuo (historias de usuario y tareas técnicas) con el ritmo de entregas académicas semanales exigidas en la materia. Sprints más largos generarían un desfase entre la funcionalidad terminada y las fechas de corrección.

## 2. Límite de Trabajo en Progreso (WIP Limit)

Se configuró un WIP Limit de **2** en la columna "In Progress". Al ser un proyecto individual, un límite de 2 evita el cambio de contexto excesivo (context-switching) pero otorga la flexibilidad necesaria para tomar una segunda tarea si la principal queda bloqueada temporalmente (por ejemplo, por una duda técnica, la caída de un servicio, o esperando el build de un pipeline), sin acumular trabajo a medio terminar.

## 3. Diagnóstico de la historia mal escrita

La historia *"Como desarrollador quiero crear la tabla usuarios para guardar los datos"* viola el principio de centrarse en el valor para el negocio; describe un detalle de implementación (una tabla SQL) desde la perspectiva técnica.
**Reescritura:** *"Como usuario, quiero poder registrarme en la plataforma para acceder a funciones personalizadas e invitar a mis amigos a planes"*.

## 4. Problemas encontrados y resolución

1. **Automatización del tablero:** Al cerrar issues manualmente o a través de Pull Requests, las tarjetas quedaban atascadas en la columna "In Progress". Se resolvió accediendo a la configuración de los Workflows internos del GitHub Project y activando explícitamente el disparador de transición "Item closed -> Status: Done".


2. **Desconexión en la Trazabilidad:** Los primeros PRs no cerraban las tareas al hacer el merge. Se solucionó adoptando la convención estricta de incluir la palabra clave `Closes #N` en la descripción del PR, lo que habilitó el enlace nativo de GitHub hacia el Issue tracker.

## 5. Declaración de uso de IA

Para este TP, se utilizó un Agente LLM (Gemini) como "Thought Partner" para estructurar los requerimientos y asegurar el cumplimiento del criterio INVEST en el backlog. Por otro lado, la herramienta Antigravity (VS Code) asistió en la generación iterativa del código base para cumplir dichas tareas. La verificación se realizó de forma empírica: compilando el código localmente, testeando los endpoints con clientes HTTP y verificando el cumplimiento de los criterios de aceptación detallados en cada historia.

---

# Decisiones - TP4

## 1. Estructura del pipeline

Se estructuró el archivo `ci.yml` con dos jobs separados (`build-backend` y `build-frontend`) configurados para ejecutarse **en paralelo**. Al no haber dependencias de compilación entre el código de Go y los componentes de React, ejecutar estos procesos en máquinas virtuales (runners) distintas ahorra tiempo de feedback al desarrollador. El pipeline actúa como un "Gate" estricto activado en los eventos de `push` y `pull_request` hacia la rama `main`.

## 2. Estrategia de Caché

El pipeline implementa la API de caché nativa de GitHub Actions mediante la configuración `type=gha` en el módulo de Buildx.

* **Qué se reutiliza:** Se cachean las capas intermedias de Docker (como la descarga de los módulos de Go `go.mod` y las librerías de Node `node_modules`). Para evitar conflictos de almacenamiento, se aislaron mediante scopes (`scope=backend` y `scope=frontend`).
* **Pérdida de caché:** Si el caché se purga o caduca, el pipeline no falla; simplemente el runner descargará nuevamente todas las dependencias desde cero e incurrirá en un tiempo de ejecución significativamente mayor.

## 3. Por qué el pipeline construye con el Dockerfile

Se delega la construcción directamente a la directiva `docker/build-push-action` en lugar de instalar Go/Node en el runner por un principio de paridad de entornos. Si el CI instala su propia versión del SDK para compilar, no hay garantías de que el artefacto funcione igual en el contenedor de producción. Al usar el Dockerfile, el pipeline verifica y sella exactamente el mismo binario inmutable y la misma infraestructura que se desplegará.

## 4. Problemas encontrados y resolución

1. **Bloqueo absoluto de PRs (Gate strict):** Al activar la regla "Require branches to be up to date before merging", los PRs quedaban bloqueados en rojo indicando que la rama estaba desactualizada respecto a main. La solución fue realizar un rebase/merge previo de `main` hacia la rama de *feature*, volver a pushear, y permitir que los checks corran nuevamente con la historia integrada.
2. **Nomenclatura dinámica del runner:** Durante los Pull Requests, la variable nativa para extraer el nombre de la rama generaba conflictos (devolvía números de merge en lugar del nombre real). Se solucionó parametrizando el entorno de validación con un fallback explícito: `RAMA: ${{ github.head_ref || github.ref_name }}`.

## 5. Declaración de uso de IA

Se utilizó IA conversacional para diseñar la sintaxis YAML del pipeline y asegurar el uso correcto de las Actions oficiales actualizadas (ej. `actions/checkout@v6`, `docker/setup-buildx-action@v4`). La verificación se llevó a cabo empíricamente mediante la técnica de "broken build": introduciendo un error de sintaxis intencional en una rama, abriendo un Pull Request, y validando visualmente que el Gate bloqueara el merge de forma efectiva.
