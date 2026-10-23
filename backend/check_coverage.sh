#!/bin/sh
set -e

THRESHOLD=80.0
OUT_DIR="/out"

echo "=========================================================="
echo "🧪 Ejecutando suite de tests unitarios del Backend (Go)..."
echo "=========================================================="

mkdir -p "$OUT_DIR" coverage 2>/dev/null || true

# 1. Ejecutar go test con coverprofile sobre ./...
go test -v -coverprofile=coverage.raw.out ./...

# 2. Filtrar paquetes de entrada/infraestructura pura (cmd/api y repository que requieren DB activa)
# para medir la cobertura sobre la lógica de negocio y controladores (handlers y services)
grep -v -E 'cmd/api|internal/repository' coverage.raw.out > coverage.out

# 3. Generar reporte HTML
go tool cover -html=coverage.out -o coverage.html || true

# 4. Calcular porcentaje total de cobertura de sentencias/líneas
TOTAL_COV=$(go tool cover -func=coverage.out | grep total: | awk '{print $3}' | tr -d '%')

if [ -z "$TOTAL_COV" ]; then
  echo "❌ Error: No se pudo calcular el porcentaje de cobertura."
  exit 1
fi

echo ""
echo "📊 Cobertura total de líneas/sentencias: ${TOTAL_COV}%"
echo "🎯 Umbral mínimo requerido: ${THRESHOLD}%"
echo ""

# Copiar artefactos a /out para exportar al runner de CI si existe
if [ -d "$OUT_DIR" ]; then
  cp coverage.out "$OUT_DIR/coverage.out" 2>/dev/null || true
  cp coverage.html "$OUT_DIR/coverage.html" 2>/dev/null || true
  
  cat <<EOF > "$OUT_DIR/coverage-summary.md"
### 🧪 Resumen de Cobertura de Tests - Backend (Go)
| Métrica | Cobertura Obtenida | Umbral Requerido | Estado |
| :--- | :--- | :--- | :--- |
| **Sentencias / Líneas** | \`${TOTAL_COV}%\` | \`${THRESHOLD}%\` | $(awk -v cov="$TOTAL_COV" -v th="$THRESHOLD" 'BEGIN { if (cov >= th) print "✅ Aprobado"; else print "❌ Reprobado" }') |
EOF
fi

# 5. Evaluar umbral: salir con exit 1 si es inferior a 80%
PASSED=$(awk -v cov="$TOTAL_COV" -v th="$THRESHOLD" 'BEGIN { if (cov >= th) print "1"; else print "0" }')

if [ "$PASSED" -eq "1" ]; then
  echo "✅ APROBADO: La cobertura (${TOTAL_COV}%) supera el umbral del ${THRESHOLD}%."
  exit 0
else
  echo "❌ FRENO DE CALIDAD ACTIVADO: Cobertura (${TOTAL_COV}%) es inferior al ${THRESHOLD}% requerido."
  exit 1
fi
