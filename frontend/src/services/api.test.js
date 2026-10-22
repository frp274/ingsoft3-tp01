import { describe, it, expect, vi } from 'vitest';
import { solicitarInvitacionAPI, obtenerPlanesAPI, obtenerMisPlanesAPI } from './api';

describe('Suite de Servicios de API con Mocks (vi.fn)', () => {
  // Test 1: Mock de fetch exitoso verificado con toHaveBeenCalledWith
  it('solicitarInvitacionAPI debe llamar a fetch con la ruta y payload correctos', async () => {
    // Arrange (Patrón AAA)
    const mockRespuesta = {
      idInvitacion: 12,
      idPlan: 5,
      idUsuario: 1,
      estado: 'pendiente',
    };

    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: vi.fn().mockResolvedValue(mockRespuesta),
    });

    const idPlan = 5;

    // Act
    const resultado = await solicitarInvitacionAPI(idPlan, mockFetch);

    // Assert
    expect(mockFetch).toHaveBeenCalledTimes(1);
    expect(mockFetch).toHaveBeenCalledWith('/api/invitaciones', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ idPlan: 5 }),
    });
    expect(resultado).toEqual(mockRespuesta);
  });

  // Test 2: Error HTTP desde el backend usando Mock
  it('solicitarInvitacionAPI debe lanzar error descriptivo cuando fetch responde ok: false', async () => {
    // Arrange
    const mockFetchError = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: vi.fn().mockResolvedValue({ error: 'El plan especificado no existe' }),
    });

    // Act & Assert
    await expect(solicitarInvitacionAPI(999, mockFetchError)).rejects.toThrow(
      'El plan especificado no existe'
    );
    expect(mockFetchError).toHaveBeenCalledTimes(1);
  });

  // Test 3: Caso de borde: idPlan no positivo sin disparar fetch
  it.each([0, -1, null, undefined, 'invalido'])(
    'solicitarInvitacionAPI con idPlan=%s debe lanzar error antes de invocar la red (caso borde)',
    async (idInvalido) => {
      // Arrange
      const mockFetch = vi.fn();

      // Act & Assert
      await expect(solicitarInvitacionAPI(idInvalido, mockFetch)).rejects.toThrow(
        'El id del plan debe ser un número mayor a 0'
      );
      expect(mockFetch).not.toHaveBeenCalled();
    }
  );

  // Test 4: obtenerPlanesAPI realiza GET /api/planes
  it('obtenerPlanesAPI debe consultar /api/planes y devolver los datos en JSON', async () => {
    // Arrange
    const planesMock = [
      { id: 1, nombre: 'Trekking', visibilidad: 'publico' },
      { id: 2, nombre: 'Fiesta', visibilidad: 'publico' },
    ];

    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(planesMock),
    });

    // Act
    const resultado = await obtenerPlanesAPI(mockFetch);

    // Assert
    expect(mockFetch).toHaveBeenCalledWith('/api/planes');
    expect(resultado).toHaveLength(2);
    expect(resultado[0].nombre).toBe('Trekking');
  });

  // Test 5: obtenerMisPlanesAPI realiza GET /api/mis-planes
  it('obtenerMisPlanesAPI debe consultar /api/mis-planes exitosamente', async () => {
    // Arrange
    const misPlanesMock = [{ id: 1, nombre: 'Trekking', estado: 'pendiente' }];
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue(misPlanesMock),
    });

    // Act
    const resultado = await obtenerMisPlanesAPI(mockFetch);

    // Assert
    expect(mockFetch).toHaveBeenCalledWith('/api/mis-planes');
    expect(resultado).toEqual(misPlanesMock);
  });
});
