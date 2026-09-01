# Decisiones - TP1

## 1. Por qué Git no pudo resolver el conflicto solo — y qué habría tenido que pasar para que nunca apareciera.
Por que las dos ramas partieron de la misma (Main) y no "sabian" que se estaban pisando al mergear, y git necesita saber con que cambio final queda en main.
Para que nunca apareciera, deberia evitar que 2 ramas trabajen en los mismos archivos, y de ser asi, esperar que una mergee a main antes de la otra cree una nueva rama con cambios

## 2. Qué problemas encontraste y cómo los solucionaste. Los tropiezos bien contados valen más que un camino perfecto: son los que demuestran que entendiste.
El bloqueo del push directo, solucion: solicitar y confirmar un pull request
Un conflicto entre 2 ramas que mergean a main, solucion: resolver conflicto eligiendo alguno de los dos cambios

## 3. Declaración de uso de IA: qué partes hiciste con ayuda de inteligencia artificial y cómo verificaste lo que te devolvió (§ Uso de IA del enunciado).
La unica parte donde usea ia, predeterminada de github, fue donde te escribe el titulo y descripcion del commit, el resto nada, por que esta bien explicado

# Decisiones - TP2
## Decisiones de disenio de BD de la APP 
### Arquitectura de Persistencia y Modelo de Datos (MVP)
- **Qué:** Se adopta un modelo 100% relacional en MySQL, excluyendo MongoDB y el motor de recomendaciones en la fase inicial.
- **Cómo:** Se estructuran 4 tablas core (`Usuario`, `Plan`, `Etiqueta`, `Invitacion`). `Invitacion` actúa como tabla intermedia (relación M:N) entre usuarios y planes. El organizador se maneja con una FK directa `idOrganizador` desde `Plan` hacia `Usuario`.
- **Por qué:** Simplifica el desarrollo del MVP, asegura integridad transaccional (ACID) en el manejo de la capacidad y asistencias, y reduce la complejidad operativa de mantener infraestructura políglota tempranamente.

# Decisiones - TP3

### TP3 — Planificación y trazabilidad
- **Duración del Sprint:** Se establece en 1 semana. Esto permite alinear el esfuerzo de desarrollo (historias y tareas) con el ritmo de entregas semanales de los trabajos prácticos de la materia[cite: 3].
- **Límite de trabajo en progreso (WIP):** Se configura en 2. Al ser un proyecto individual (1 persona + 1), este límite evita la sobrecarga de contexto, permitiendo tener una tarea en revisión o bloqueada mientras se avanza con una segunda, sin acumular trabajo sin terminar[cite: 3].
- **Diagnóstico de historia mal escrita:** La historia *"Como desarrollador quiero crear la tabla usuarios para guardar los datos"* es en realidad una tarea técnica disfrazada[cite: 3]. No entrega un incremento de valor observable por el usuario final. Se reescribiría como: *"Como usuario, quiero poder registrarme en la plataforma para acceder a funciones personalizadas"*.
- **Problemas encontrados y resoluciones:** La automatización de GitHub Projects no movía las tareas a "Done". Se solucionó activando manualmente el workflow interno del tablero "Item closed -> Status: Done"[cite: 3].
- **Uso de IA:** Antigravity (VS Code) asistió en la generación del código base iterativo, mientras que Gemini validó la estructura de requerimientos (Épicas/Historias) asegurando el cumplimiento del criterio INVEST.

# Decisiones - TP4

### CI/CD y Estrategia de QA (Ajuste TP4)
- **Qué:** Implementación de Pipeline as Code con GitHub Actions para verificar la compilación (build de imágenes Docker) en cada PR, posponiendo la ejecución de tests automatizados para el TP5.
- **Cómo:** Se configuran dos jobs en paralelo (`build-backend` y `build-frontend`) usando `docker/build-push-action`, implementando caché de capas de Docker en el almacén de GitHub (`type=gha`), y configurando el workflow como *Required Status Check* en la rama main.
- **Por qué:** Cumple con el requerimiento del TP4 de garantizar que el código compila en un entorno limpio antes de integrarse. Separar la compilación (TP4) de los tests (TP5) permite iterar rápido sobre la infraestructura base sin bloquearse escribiendo pruebas tempranas.
