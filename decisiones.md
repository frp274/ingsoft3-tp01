# Decisiones - TP1

## 1. Por qué Git no pudo resolver el conflicto solo — y qué habría tenido que pasar para que nunca apareciera.
Por que las dos ramas partieron de la misma (Main) y no "sabian" que se estaban pisando al mergear, y git necesita saber con que cambio final queda en main.
Para que nunca apareciera, deberia evitar que 2 ramas trabajen en los mismos archivos, y de ser asi, esperar que una mergee a main antes de la otra cree una nueva rama con cambios

## 2. Qué problemas encontraste y cómo los solucionaste. Los tropiezos bien contados valen más que un camino perfecto: son los que demuestran que entendiste.
El bloqueo del push directo, solucion: solicitar y confirmar un pull request
Un conflicto entre 2 ramas que mergean a main, solucion: resolver conflicto eligiendo alguno de los dos cambios

## 3. Declaración de uso de IA: qué partes hiciste con ayuda de inteligencia artificial y cómo verificaste lo que te devolvió (§ Uso de IA del enunciado).
La unica parte donde usea ia, predeterminada de github, fue donde te escribe el titulo y descripcion del commit, el resto nada, por que esta bien explicado
