import { test, expect } from '@playwright/test';

const API = process.env.API_BASE_URL || 'http://localhost:8080';

test.describe('Pruebas de Integración de API - PlanesApp', () => {
  test('Prueba A: Crear un plan, verificar existencia con GET y borrar con DELETE', async ({ request }) => {
    const tituloUnico = `Plan E2E ${Date.now()}`;
    const nuevoPlan = {
      nombre: tituloUnico,
      fechaInicio: '2026-11-01T10:00:00Z',
      fechaFin: '2026-11-01T18:00:00Z',
      visibilidad: 'publico',
      descripcion: 'Plan generado automáticamente',
      idEtiqueta: 1,
      ubicacion: 'Auditorio',
      capacidad: 50,
    };

    // 1. Crear el plan y verificar código 201
    const createRes = await request.post(`${API}/api/planes`, { data: nuevoPlan });
    expect(createRes.status()).toBe(201);
    const planCreado = await createRes.json();
    const idReal = planCreado.id || planCreado.ID;

    let borrado = false;
    try {
      // 2. Hacer GET para verificar que aparece
      const getRes = await request.get(`${API}/api/planes`);
      expect(getRes.ok()).toBeTruthy();
      const planes = await getRes.json();
      const planEncontrado = planes.find((p) => p.nombre === tituloUnico);
      expect(planEncontrado).toBeDefined();

      // 3. Hacer DELETE y verificar que el backend responde OK
      const deleteRes = await request.delete(`${API}/api/planes/${idReal}`);
      expect(deleteRes.status()).toBe(200);
      borrado = true;

      // 4. Confirmar persistencia del borrado: GET por ID devuelve 404 y no aparece en el listado
      const getDeletedRes = await request.get(`${API}/api/planes/${idReal}`);
      expect(getDeletedRes.status()).toBe(404);

      const getAfterRes = await request.get(`${API}/api/planes`);
      expect(getAfterRes.ok()).toBeTruthy();
      const planesDespues = await getAfterRes.json();
      expect(planesDespues.find((p) => p.nombre === tituloUnico)).toBeUndefined();
    } finally {
      if (!borrado && idReal) {
        await request.delete(`${API}/api/planes/${idReal}`).catch(() => {});
      }
    }
  });

  // Prueba B: Mandar un plan inválido (ej. título vacío), esperar que la API devuelva HTTP 400,
  // y hacer un GET para asegurar que la cantidad de planes no cambió.
  test('Prueba B: Mandar un plan inválido, esperar HTTP 400 y verificar que la cantidad de planes no cambió', async ({ request }) => {
    // 1. Obtener cantidad inicial de planes
    const listInicialRes = await request.get(`${API}/api/planes`);
    expect(listInicialRes.ok()).toBeTruthy();
    const planesIniciales = await listInicialRes.json();
    const totalInicial = planesIniciales.length;

    // 2. Mandar plan inválido (nombre / título vacío)
    const planInvalido = {
      nombre: '',
      fechaInicio: '2026-11-01T10:00:00Z',
      fechaFin: '2026-11-01T18:00:00Z',
      visibilidad: 'publico',
      idEtiqueta: 1,
      ubicacion: 'Lugar Inválido',
      capacidad: 25,
    };
    const invalidRes = await request.post(`${API}/api/planes`, {
      data: planInvalido,
    });
    expect(invalidRes.status()).toBe(400);

    // 3. Hacer GET para asegurar que la cantidad de planes no cambió
    const listFinalRes = await request.get(`${API}/api/planes`);
    expect(listFinalRes.ok()).toBeTruthy();
    const planesFinales = await listFinalRes.json();
    expect(planesFinales.length).toBe(totalInicial);
  });

  // Prueba C: Buscar un plan con un ID inexistente y validar que devuelva 404, o listar los planes y validar la estructura.
  test('Prueba C: Buscar un plan con ID inexistente (404) o listar los planes y validar la estructura', async ({ request }) => {
    // 1. Listar los planes y validar la estructura de los datos devueltos
    const listRes = await request.get(`${API}/api/planes`);
    expect(listRes.status()).toBe(200);
    const planes = await listRes.json();
    expect(Array.isArray(planes)).toBeTruthy();

    if (planes.length > 0) {
      const primerPlan = planes[0];
      expect(primerPlan).toHaveProperty('id');
      expect(primerPlan).toHaveProperty('nombre');
      expect(primerPlan).toHaveProperty('ubicacion');
      expect(primerPlan).toHaveProperty('capacidad');
    }

    // 2. Buscar un plan con ID inexistente
    const notFoundRes = await request.get(`${API}/api/planes/99999999`);
    expect(notFoundRes.status()).toBe(404);
  });
});
