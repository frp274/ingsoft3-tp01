import { describe, it, expect } from 'vitest';
import { formatearFecha, formatearEstadoInvitacion, calcularDuracionHoras } from './formatters';

describe('Suite de Funciones Formateadoras y Calculadoras', () => {
  // Test 1: Caso borde de formatearFecha ante valores nulos, vacíos o fechas inválidas
  it('formatearFecha debe manejar valores nulos, vacíos e inválidos de forma segura', () => {
    // Arrange & Act & Assert
    expect(formatearFecha(null)).toBe('');
    expect(formatearFecha(undefined)).toBe('');
    expect(formatearFecha('')).toBe('');
    expect(formatearFecha('fecha-invalida-xyz')).toBe('Fecha inválida');
  });

  // Test 2: Formateo correcto de fecha válida
  it('formatearFecha debe formatear una fecha válida en español', () => {
    // Arrange
    const fechaISO = '2026-11-20T15:30:00Z';

    // Act
    const resultado = formatearFecha(fechaISO);

    // Assert
    expect(resultado).toBeTruthy();
    expect(typeof resultado).toBe('string');
    expect(resultado).toContain('2026');
  });

  // Test 3: Test parametrizado con it.each para formatearEstadoInvitacion
  it.each([
    { estado: 'aceptada', etiqueta: '✅ Confirmado', clase: 'aceptada' },
    { estado: 'Aceptada', etiqueta: '✅ Confirmado', clase: 'aceptada' },
    { estado: 'rechazada', etiqueta: '❌ Rechazada', clase: 'rechazada' },
    { estado: 'pendiente', etiqueta: '⏳ Solicitud Pendiente', clase: 'pendiente' },
    { estado: 'PENDIENTE', etiqueta: '⏳ Solicitud Pendiente', clase: 'pendiente' },
    { estado: 'otra_cosa', etiqueta: 'Desconocido', clase: 'desconocido' },
    { estado: null, etiqueta: 'Desconocido', clase: 'desconocido' },
    { estado: undefined, etiqueta: 'Desconocido', clase: 'desconocido' },
  ])('formatearEstadoInvitacion("$estado") debe devolver "$etiqueta"', ({ estado, etiqueta, clase }) => {
    // Act
    const res = formatearEstadoInvitacion(estado);

    // Assert
    expect(res.etiqueta).toBe(etiqueta);
    expect(res.clase).toBe(clase);
  });

  // Test 4: Cálculo de duración en horas
  it('calcularDuracionHoras debe calcular la diferencia correcta y manejar fechas invertidas', () => {
    // Arrange
    const inicio = '2026-11-20T10:00:00Z';
    const finNormal = '2026-11-20T12:30:00Z';
    const finInvertida = '2026-11-20T08:00:00Z';

    // Act
    const duracionNormal = calcularDuracionHoras(inicio, finNormal);
    const duracionInvertida = calcularDuracionHoras(inicio, finInvertida);
    const duracionNula = calcularDuracionHoras(null, finNormal);

    // Assert
    expect(duracionNormal).toBe(2.5);
    expect(duracionInvertida).toBe(0);
    expect(duracionNula).toBe(0);
  });
});
