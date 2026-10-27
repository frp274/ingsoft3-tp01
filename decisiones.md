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

---

# Decisiones - TP5

## 1. Enlaces de verificación y evidencias de este TP
- **Corrida roja en GitHub Actions (freno por umbral de cobertura):**  
  [Fallo en CI por cobertura insuficiente (<80%)](https://github.com/frp274/ingsoft3-tp01/actions/runs/39980746055/job/120507063296?pr=27)
- **PR #27 (Secuencia completa: rojo -> corrección con tests -> verde -> merge):**  
  [Pull Request #27 - Implementación de Quality Gate y resolución](https://github.com/frp274/ingsoft3-tp01/pull/27)
- **PR #28 (Gate bloqueado permanentemente - Evidencia activa para la defensa):**  
  [Pull Request #28 - Freno de calidad activo (en rojo)](https://github.com/frp274/ingsoft3-tp01/pull/28)

---

## 2. Qué lógica se testeó y por qué (análisis de impacto)
Se priorizó testear la lógica nuclear del negocio del Organizador de Eventos:
1. **En el backend (Go):** Las reglas de validación de planes, cálculo y control de capacidad máxima de asistentes, verificación de estados de invitaciones y transiciones de estados del plan.
2. **En el frontend (React):** Los módulos utilitarios puros (`src/lib/`), encargados de procesar fechas, filtrar listas de eventos según visibilidad/estado y clasificar prioridades de planes.

**Por qué esa lógica:** Un bug en estos componentes compromete directamente la integridad del sistema (por ejemplo, permitir sobreventa de cupos para un evento, procesar invitaciones a planes cancelados o admitir usuarios no registrados). Los fallos en la lógica de negocio no se manifiestan como caídas del servidor; son errores silenciosos que corrompen los datos.

---

## 3. Adaptación del Stack tecnológico frente a la guía (.NET -> Go + React/Vitest)
Dado que el stack del proyecto utiliza Go y React en lugar del sample canónico (.NET + Vitest), se mapearon todas las técnicas y herramientas requeridas para cumplir con la matriz de la cátedra:

| Concepto a lograr | Herramienta implementada | Justificación técnica |
|---|---|---|
| **Ubicación de tests (Back)** | Archivos `*_test.go` | Convención estándar nativa de Go junto a sus paquetes de lógica. |
| **Ubicación de tests (Front)** | Archivos `*.test.js` en `src/lib/` | Pruebas de lógica pura al lado de los módulos testeados sin tocar el DOM. |
| **Test parametrizado (Back)** | *Table-Driven Tests* (`t.Run`) | Estructuras de datos iteradas con slices para evaluar múltiples inputs/outputs en un solo método. |
| **Test parametrizado (Front)** | `it.each` de Vitest | Mismo comportamiento de evaluación tabular en entorno JavaScript. |
| **Inyección de dependencias** | Interfaces en Go / Funciones por parámetro en JS | Permite sustituir repositorios y clientes HTTP por dobles de prueba desde el exterior. |
| **Fabricación de Mocks (Back)** | `stretchr/testify/mock` | Mocks basados en interfaces que permiten aserciones sobre llamadas (`AssertExpectations`). |
| **Fabricación de Mocks (Front)** | `vi.fn()` en Vitest | Mocking de llamadas asíncronas y verificación de argumentos con `toHaveBeenCalledWith`. |
| **Medición de cobertura** | `go test -coverprofile` y `@vitest/coverage-v8` | Extracción de porcentajes de líneas y ramas tanto en local como en los runners del CI. |
| **Umbral que rompe el build** | Script bash con `exit 1` (Back) y `thresholds` (Front) | Interrupción de la ejecución y salida con código de error si el porcentaje es menor a 80%. |
| **Filtro de exclusión** | Banderas de paquetes en Go y `include: ['src/lib/**']` en Vite | Aísla el código de arranque y modelos de datos para medir exclusivamente la lógica de negocio. |

---

## 4. Umbral de cobertura elegido y justificación
Se definió un umbral estricto del **80% de cobertura** tanto para líneas (*lines*) como para ramas de decisión (*branches*).
- **Justificación:** Un 80% asegura que ninguna función central quede desprotegida ante refactorizaciones o regresiones, obligando a contemplar caminos alternativos y casos de borde. Se evitó exigir un 95% o 100% para no caer en la *Ley de Goodhart*, donde el equipo termina escribiendo aserciones vacías o testeando código trivial solo para satisfacer una métrica numérica.
- **Ramas medidas:** En el frontend, Vitest evalúa tanto *lines* como *branches*. En Go, si bien `go test -cover` mide sentencias (*statements*), la suite fue diseñada recorriendo explícitamente todas las ramas de control (`if/else` y `switch`).

---

## 5. Qué se dejó afuera de la cuenta de cobertura y por qué
Tanto en el backend como en el frontend se aplicó el principio de **"medir lo que tiene lógica y excluir la fontanería"**:
- **Backend (Go):** Se excluyeron `main.go` y la carpeta de arranque (`cmd/api/`), los modelos/structs de datos puros que solo definen esquemas sin comportamiento, y las inicializaciones de conexiones a MySQL. Si el cableado de la app falla, el servicio no levanta y el problema se detecta de inmediato sin necesidad de tests unitarios.
- **Frontend (React):** Se excluyeron los componentes visuales de React (`.jsx` / `.tsx`), los archivos de configuración (`vite.config.js`) y el punto de entrada `index.html`. La cátedra desaconseja montar entornos DOM pesados en esta etapa; la UI se valida de punta a punta con pruebas end-to-end en el TP7. La cobertura del frontend mide exclusivamente la lógica pura contenida en `src/lib/**`.

---

## 6. Por qué cobertura alta no garantiza calidad (con ejemplo concreto)
El code coverage mide qué líneas fueron **ejecutadas**, pero es ciego respecto a si fueron **verificadas**. 

**Ejemplo:**
```javascript
// Función bajo prueba:
export function calcularCapacidadRestante(total, confirmados) {
    if (confirmados > total) return 0;
    return total - confirmados;
}

// Test engañoso que otorga 100% de cobertura:
it('ejecuta la función de capacidad', () => {
    calcularCapacidadRestante(100, 20); // Se ejecutan todas las líneas... ¡pero no hay ningún expect()!
});
```
---

# Decisiones - TP6
## Despliegue Continuo (CD) y Gestión de Entornos en la Nube
---
### 1. Decisiones de Infraestructura y Despliegue
* **Plataforma de Hosting (Render):** Se optó por utilizar Render mediante Web Services configurados con runtime de Docker y el *Auto-Deploy* desactivado. Esto nos permite un control total sobre cuándo se actualizan los servicios y asegura que el despliegue sea coordinado exclusivamente desde el pipeline de integración continua.
* **Separación de Entornos (QA vs. Producción):** Se estructuraron cuatro servicios independientes (`api-qa`, `api-prod`, `front-qa`, `front-prod`) para garantizar el aislamiento absoluto entre el entorno de pruebas y el entorno real de producción, evitando interferencias en los datos y comportamientos de la aplicación.
* **Configuración Dinámica de Nginx en el Frontend:** Se reemplazó el archivo de configuración estático por una plantilla (`default.conf.template`) utilizando variables de entorno (`BACKEND_URL`) y un resolvedor DNS dinámico (`DNS_RESOLVER=8.8.8.8`). Esto evita hardcodear nombres internos de contenedores y permite que el frontend redirija correctamente las peticiones API hacia el backend correspondiente según el entorno.
---
### 2. Elección de Base de Datos: Aiven vs. Neon
* **Decisión:** Se seleccionó **Aiven** como proveedor de base de datos MySQL en la nube en lugar de **Neon**.
* **Justificación técnica:** 
  1. *Compatibilidad de motor:* Neon es un servicio exclusivo para bases de datos **PostgreSQL**. El backend del proyecto está desarrollado en Go utilizando el driver oficial de MySQL (`go-sql-driver/mysql`) y consultas SQL nativas de dicho motor. Utilizar Neon habría implicado una reescritura profunda del backend, cambio de driver y validación completa de consultas, introduciendo un riesgo innecesario de rotura en una etapa avanzada del proyecto. Aiven nos proveyó una instancia de MySQL gratuita y compatible al 100% sin tocar el código fuente.
  2. *Aislamiento de entornos con recursos limitados:* Ante la restricción del plan gratuito de Aiven (límite de un servicio/servidor de base de datos), se resolvió la separación entre QA y Producción utilizando **bases de datos lógicas independientes** dentro del mismo servidor (`defaultdb` para QA y `proddb` para Producción), cumpliendo rigurosamente con el requisito de no mezclar los datos de ambos entornos.
---
### 3. Automatización del Despliegue Continuo (CD)
* **GitHub Actions y GHCR:** Se configuró el pipeline para compilar las imágenes Docker del backend y frontend y publicarlas automáticamente en el GitHub Container Registry (GHCR) etiquetadas con el *commit sha* correspondiente.
* **Deploy Hooks y Environments:** Se integró GitHub Actions con Render utilizando *Deploy Hooks* seguros almacenados como secretos en los *Environments* de GitHub (`qa` y `production`).
* **Compuerta Humana (Required Reviewers):** Se configuró el entorno de `production` en GitHub exigiendo una aprobación manual obligatoria antes de que el pipeline ejecute el despliegue hacia producción, implementando formalmente una estrategia de *Continuous Delivery*.
* **Congelamiento de Versión:** Se cerró el práctico generando y publicando el tag y la release `v6.0.0` en GitHub, garantizando la trazabilidad exacta del estado de la aplicación.


---

## TP7 - CI/CD y Pruebas Automáticas

Se implementó un pipeline en GitHub Actions con separación de entornos (QA y PROD) usando Render y contenedores inmutables desde GHCR. Se agregó un gate de pruebas de integración y E2E con Playwright. Se modificó la aserción de borrado en pruebas de integración para lidiar con el soft-delete del backend en Go, y se extrajo el ID dinámicamente de la base de datos. Uso de IA: Asistencia para armar jobs de CI/CD, configuración de Playwright y resolución de conflictos de Git.


---

## TP8 - Infraestructura como Código

Se implementó un entorno de `preprod` utilizando Terraform con los providers de Aiven (MySQL) y Render. Se configuró una base de datos lógica aislada en la instancia existente de Aiven para optimizar recursos, y se asoció al Render Owner ID del equipo. Uso de IA: Asistencia en la creación del código HCL, debugging de IDs de la API de Render y resolución de exportación de variables en Windows.
