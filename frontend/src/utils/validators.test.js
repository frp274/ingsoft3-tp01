import { describe, it, expect } from 'vitest';
import { validarPlan, validarCapacidad } from './validators';

describe('Suite de Validadores de Planes', () => {
  // Test 1: Flujo exitoso con datos válidos
  it('debe validar exitosamente un plan con todos los campos correctos', () => {
    // Arrange
    const planValido = {
      nombre: 'Trekking Cerro Champaquí',
      idEtiqueta: 1,
      fechaInicio: '2026-11-20T08:00:00Z',
      fechaFin: '2026-11-20T18:00:00Z',
      ubicacion: 'Villa Alpina, Córdoba',
      capacidad: 15,
      visibilidad: 'publico',
    };

    // Act
    const resultado = validarPlan(planValido);

    // Assert
    expect(resultado.valido).toBe(true);
    expect(resultado.error).toBeNull();
  });

  // Test 2: Test Parametrizado usando test.each / it.each
  it.each([
    {
      caso: 'nombre vacío o con solo espacios',
      datos: { nombre: '   ', idEtiqueta: 1, fechaInicio: '2026-11-20T08:00', fechaFin: '2026-11-20T12:00', ubicacion: 'Córdoba', capacidad: 10 },
      errorEsperado: 'El nombre del plan es obligatorio.',
    },
    {
      caso: 'etiqueta no seleccionada o cero',
      datos: { nombre: 'Asado', idEtiqueta: 0, fechaInicio: '2026-11-20T08:00', fechaFin: '2026-11-20T12:00', ubicacion: 'Córdoba', capacidad: 10 },
      errorEsperado: 'Debes seleccionar una etiqueta/categoría válida.',
    },
    {
      caso: 'fecha de fin anterior o igual a la de inicio',
      datos: { nombre: 'Asado', idEtiqueta: 1, fechaInicio: '2026-11-20T18:00', fechaFin: '2026-11-20T12:00', ubicacion: 'Córdoba', capacidad: 10 },
      errorEsperado: 'La fecha y hora de fin debe ser posterior a la fecha de inicio.',
    },
    {
      caso: 'ubicación vacía',
      datos: { nombre: 'Asado', idEtiqueta: 1, fechaInicio: '2026-11-20T12:00', fechaFin: '2026-11-20T15:00', ubicacion: '   ', capacidad: 10 },
      errorEsperado: 'La ubicación es obligatoria.',
    },
    {
      caso: 'capacidad igual a cero',
      datos: { nombre: 'Asado', idEtiqueta: 1, fechaInicio: '2026-11-20T12:00', fechaFin: '2026-11-20T15:00', ubicacion: 'Córdoba', capacidad: 0 },
      errorEsperado: 'La capacidad debe ser un número mayor a 0.',
    },
    {
      caso: 'capacidad negativa',
      datos: { nombre: 'Asado', idEtiqueta: 1, fechaInicio: '2026-11-20T12:00', fechaFin: '2026-11-20T15:00', ubicacion: 'Córdoba', capacidad: -5 },
      errorEsperado: 'La capacidad debe ser un número mayor a 0.',
    },
  ])('debe rechazar plan con $caso', ({ datos, errorEsperado }) => {
    // Act
    const resultado = validarPlan(datos);

    // Assert
    expect(resultado.valido).toBe(false);
    expect(resultado.error).toBe(errorEsperado);
  });

  // Test 3: Casos de borde y entradas nulas/indefinidas
  it('debe rechazar entradas nulas o con tipos de datos inválidos (caso borde)', () => {
    // Arrange & Act
    const resNull = validarPlan(null);
    const resUndefined = validarPlan(undefined);
    const resString = validarPlan('no-es-objeto');

    // Assert
    expect(resNull.valido).toBe(false);
    expect(resNull.error).toBe('Los datos del plan son obligatorios.');

    expect(resUndefined.valido).toBe(false);
    expect(resString.valido).toBe(false);
  });

  // Test 4: Validación de capacidad unitaria
  it.each([
    [10, true],
    ['15', true],
    [1, true],
    [0, false],
    [-4, false],
    ['abc', false],
    [null, false],
    [undefined, false],
    ['', false],
    [3.14, false],
  ])('validarCapacidad(%s) debe retornar %s', (valor, esperado) => {
    expect(validarCapacidad(valor)).toBe(esperado);
  });
});
