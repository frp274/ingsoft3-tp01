import { test, expect } from '@playwright/test';

test.describe('Pruebas E2E de Interfaz - Gestión de Planes', () => {
  // Prueba A (Alta y borrado):
  // Navegar a '/', llenar el formulario de creación con un título único (Date.now()),
  // hacer clic en el botón de agregar, verificar que el texto del plan aparece visible en la lista,
  // hacer clic en el botón de borrar ese plan específico, y verificar que desaparece.
  test('Prueba A (Alta y borrado): Crear un plan, verificar su aparición en la lista y borrarlo', async ({ page }) => {
    await page.goto('/');

    // 1. Abrir vista de creación de plan desde la navegación
    const btnCrearNav = page.getByRole('button', { name: /\+ Crear Plan/i });
    await expect(btnCrearNav).toBeVisible();
    await btnCrearNav.click();

    // 2. Completar formulario con título único
    const tituloUnico = `Plan E2E ${Date.now()}`;
    await page.locator('#nombre').fill(tituloUnico);
    await page.locator('#fechaInicio').fill('2026-11-15T10:00');
    await page.locator('#fechaFin').fill('2026-11-15T18:00');
    await page.locator('#ubicacion').fill('Parque Sarmiento, Córdoba');
    await page.locator('#capacidad').fill('25');
    await page.locator('#descripcion').fill('Plan de prueba E2E automatizado con Playwright');

    // 3. Hacer clic en el botón de agregar / publicar
    const btnPublicar = page.getByRole('button', { name: /Publicar Plan|Agregar/i });
    await expect(btnPublicar).toBeEnabled();
    await btnPublicar.click();

    // 4. Verificar que el texto del plan aparece visible en la lista (feed de novedades)
    const planCard = page.locator('.plan-card').filter({ hasText: tituloUnico });
    await expect(planCard).toBeVisible();
    await expect(planCard.locator('.plan-titulo')).toHaveText(tituloUnico);

    // 5. Hacer clic en el botón de borrar de ese plan específico
    const btnBorrar = planCard.getByRole('button', { name: /Borrar|Eliminar/i });
    await expect(btnBorrar).toBeVisible();
    await btnBorrar.click();

    // 6. Verificar que desaparece de la vista
    await expect(planCard).not.toBeVisible();
  });

  // Prueba B (Validación):
  // Navegar a '/', intentar crear un plan inválido, hacer clic en agregar,
  // y verificar que aparece un mensaje de error visual (usando getByRole 'alert' o similar) y que no se agregó a la lista.
  test('Prueba B (Validación): Formulario inválido muestra mensaje de error y no agrega el plan', async ({ page }) => {
    await page.goto('/');

    // 1. Ir a la vista del formulario
    const btnCrearNav = page.getByRole('button', { name: /\+ Crear Plan/i });
    await expect(btnCrearNav).toBeVisible();
    await btnCrearNav.click();

    // 2. Completar campos con fechas inválidas (fin anterior a inicio) para disparar validación de negocio
    const tituloInvalido = `Plan Inválido ${Date.now()}`;
    await page.locator('#nombre').fill(tituloInvalido);
    await page.locator('#fechaInicio').fill('2026-11-20T18:00');
    await page.locator('#fechaFin').fill('2026-11-20T10:00'); // Fin antes de inicio
    await page.locator('#ubicacion').fill('Auditorio Central');
    await page.locator('#capacidad').fill('15');

    // 3. Hacer clic en agregar
    const btnPublicar = page.getByRole('button', { name: /Publicar Plan|Agregar/i });
    await btnPublicar.click();

    // 4. Verificar mensaje de error visual accesible mediante getByRole 'alert'
    const errorAlert = page.getByRole('alert').or(page.locator('.alert-box.alert-error'));
    await expect(errorAlert).toBeVisible();
    await expect(errorAlert).toContainText(/posterior|error|obligatorio|inválido/i);

    // 5. Volver al feed de Novedades y verificar que el plan inválido NO fue agregado
    const btnNovedades = page.getByRole('button', { name: /Novedades/i });
    await btnNovedades.click();

    const planNoCreado = page.locator('.plan-card').filter({ hasText: tituloInvalido });
    await expect(planNoCreado).toHaveCount(0);
  });

  // Prueba C (Persistencia visual):
  // Cargar la página, contar cuántos planes hay en la lista, recargar la página con `page.reload()`,
  // y verificar que la lista mantiene la misma cantidad de elementos.
  test('Prueba C (Persistencia visual): La cantidad de planes se mantiene al recargar la página', async ({ page }) => {
    await page.goto('/');

    // 1. Asegurar que estamos en el listado y esperar que termine la carga inicial
    await expect(page.locator('.status-message.loading')).not.toBeVisible();

    // 2. Contar cantidad de tarjetas de planes visibles
    const planesLocator = page.locator('.plan-card');
    const cantidadInicial = await planesLocator.count();

    // 3. Recargar la página usando page.reload()
    await page.reload();

    // 4. Esperar nuevamente el fin de la carga
    await expect(page.locator('.status-message.loading')).not.toBeVisible();

    // 5. Verificar que la lista mantiene exactamente la misma cantidad de elementos
    await expect(planesLocator).toHaveCount(cantidadInicial);
  });
});
