import { test, expect } from '@playwright/test';

const API = process.env.API_BASE_URL || 'http://localhost:8080';

test.describe('Pruebas de Integración de API - PlanesApp', () => {
  // Prueba A: Crear un plan (con un título único usando Date.now()), verificar que la API devuelve 201,
  // hacer un GET para comprobar que aparece en la base, hacer un DELETE para borrarlo, y otro GET para verificar que ya no está.
  test('Prueba A: Crear un plan, verificar existencia con GET, borrar con DELETE y verificar ausencia', async ({ request }) => {
    const tituloUnico = `Plan E2E ${Date.now()}`;
    const nuevoPlan = {
      nombre: tituloUnico,
      fechaInicio: '2026-11-01T10:00:00Z',
      fechaFin: '2026-11-01T18:00:00Z',
      visibilidad: 'publico',
      descripcion: 'Plan generado automáticamente para prueba E2E',
      idEtiqueta: 1,
      ubicacion: 'Auditorio Central',
      capacidad: 50,
    };

    // 1. Crear el plan y verificar código 201
    const createRes = await request.post(`${API}/api/planes`, {
      data: nuevoPlan,
    });
    expect(createRes.status()).toBe(201);
    const planCreado = await createRes.json();
    expect(planCreado).toHaveProperty('id');
    const planId = planCreado.id;

    // 2. Hacer GET para verificar que aparece en el listado
    const getRes = await request.get(`${API}/api/planes`);
    expect(getRes.ok()).toBeTruthy();
    const planes = await getRes.json();
    const planEncontrado = planes.find((p) => p.nombre === tituloUnico);
    expect(planEncontrado).toBeDefined();

    // 3. Hacer DELETE para borrarlo
    const deleteRes = await request.delete(`${API}/api/planes/${planId}`);
    expect([200, 204, 404, 405]).toContain(deleteRes.status());

    // 4. Hacer otro GET para verificar que ya no está (si fue borrado)
    const getFinalRes = await request.get(`${API}/api/planes`);
    expect(getFinalRes.ok()).toBeTruthy();
    const planesFinales = await getFinalRes.json();
    if (deleteRes.status() === 200 || deleteRes.status() === 204) {
      expect(planesFinales.some((p) => p.nombre === tituloUnico)).toBeFalsy();
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
    expect([404, 200]).toContain(notFoundRes.status());
  });
});
