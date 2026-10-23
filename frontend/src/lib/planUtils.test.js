import { describe, it, expect } from 'vitest';
import { calcularPrioridadPlan, calcularDescuentoPlan } from './planUtils';

describe('planUtils - calcularPrioridadPlan (Patrón AAA)', () => {
  it('debe retornar BAJA cuando el plan es nulo o indefinido', () => {
    // Arrange
    const plan = null;

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('BAJA');
  });

  it('debe retornar CRITICA para plan PUBLICO con capacidad mayor a 50', () => {
    // Arrange
    const plan = { visibilidad: 'PUBLICO', capacidadMax: 60 };

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('CRITICA');
  });

  it('debe retornar ALTA para plan PUBLICO con capacidad entre 21 y 50', () => {
    // Arrange
    const plan = { visibilidad: 'PUBLICO', capacidadMax: 30 };

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('ALTA');
  });

  it('debe retornar MEDIA-ALTA para plan PRIVADO con capacidad mayor a 30', () => {
    // Arrange
    const plan = { visibilidad: 'PRIVADO', capacidadMax: 35 };

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('MEDIA-ALTA');
  });

  it('debe retornar MEDIA cuando la capacidad es mayor a 10 pero no cumple condiciones anteriores', () => {
    // Arrange
    const plan = { visibilidad: 'PRIVADO', capacidadMax: 15 };

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('MEDIA');
  });

  it('debe retornar BAJA cuando la capacidad es 10 o menor', () => {
    // Arrange
    const plan = { visibilidad: 'PRIVADO', capacidadMax: 5 };

    // Act
    const resultado = calcularPrioridadPlan(plan);

    // Assert
    expect(resultado).toBe('BAJA');
  });
});

describe('planUtils - calcularDescuentoPlan (Patrón AAA)', () => {
  it('debe retornar 0 si falta el plan o el tipo de usuario', () => {
    // Arrange & Act
    const sinPlan = calcularDescuentoPlan(null, 'VIP');
    const sinTipo = calcularDescuentoPlan({ capacidadMax: 20 }, null);

    // Assert
    expect(sinPlan).toBe(0);
    expect(sinTipo).toBe(0);
  });

  it('debe otorgar 35% a usuarios VIP si la capacidad supera 50, o 25% si no', () => {
    // Arrange
    const planGrande = { capacidadMax: 60 };
    const planChico = { capacidadMax: 40 };

    // Act
    const descGrande = calcularDescuentoPlan(planGrande, 'VIP');
    const descChico = calcularDescuentoPlan(planChico, 'VIP');

    // Assert
    expect(descGrande).toBe(35);
    expect(descChico).toBe(25);
  });

  it('debe otorgar 20% a usuarios PREMIUM si la capacidad supera 20, o 15% si no', () => {
    // Arrange
    const planMedio = { capacidadMax: 25 };
    const planChico = { capacidadMax: 15 };

    // Act
    const descMedio = calcularDescuentoPlan(planMedio, 'PREMIUM');
    const descChico = calcularDescuentoPlan(planChico, 'PREMIUM');

    // Assert
    expect(descMedio).toBe(20);
    expect(descChico).toBe(15);
  });

  it('debe otorgar 15% a ESTUDIANTE para plan PUBLICO y 10% para PRIVADO', () => {
    // Arrange
    const planPublico = { visibilidad: 'PUBLICO' };
    const planPrivado = { visibilidad: 'PRIVADO' };

    // Act
    const descPub = calcularDescuentoPlan(planPublico, 'ESTUDIANTE');
    const descPriv = calcularDescuentoPlan(planPrivado, 'ESTUDIANTE');

    // Assert
    expect(descPub).toBe(15);
    expect(descPriv).toBe(10);
  });

  it('debe otorgar 5% a usuarios REGULAR', () => {
    // Arrange
    const plan = { capacidadMax: 10 };

    // Act
    const desc = calcularDescuentoPlan(plan, 'REGULAR');

    // Assert
    expect(desc).toBe(5);
  });

  it('debe retornar 0 para tipos de usuario desconocidos (default)', () => {
    // Arrange
    const plan = { capacidadMax: 10 };

    // Act
    const desc = calcularDescuentoPlan(plan, 'DESCONOCIDO');

    // Assert
    expect(desc).toBe(0);
  });
});
