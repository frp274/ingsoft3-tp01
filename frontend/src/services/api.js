/**
 * Realiza la llamada a la API para solicitar una invitación a un plan.
 * @param {number} idPlan 
 * @param {Function} [fetchFn=fetch] Función fetch inyectable para testing
 * @returns {Promise<Object>}
 */
export async function solicitarInvitacionAPI(idPlan, fetchFn = fetch) {
  if (!idPlan || typeof idPlan !== 'number' || idPlan <= 0) {
    throw new Error('El id del plan debe ser un número mayor a 0');
  }

  const res = await fetchFn('/api/invitaciones', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ idPlan }),
  });

  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || `Error ${res.status}: no se pudo solicitar invitación`);
  }

  return await res.json();
}

/**
 * Consulta la lista de planes públicos desde el backend.
 * @param {Function} [fetchFn=fetch]
 * @returns {Promise<Array>}
 */
export async function obtenerPlanesAPI(fetchFn = fetch) {
  const res = await fetchFn('/api/planes');
  if (!res.ok) {
    throw new Error(`Error ${res.status}: No se pudieron cargar los planes`);
  }
  return await res.json();
}

/**
 * Consulta los planes del usuario activo desde el backend.
 * @param {Function} [fetchFn=fetch]
 * @returns {Promise<Array>}
 */
export async function obtenerMisPlanesAPI(fetchFn = fetch) {
  const res = await fetchFn('/api/mis-planes');
  if (!res.ok) {
    throw new Error(`Error ${res.status}: No se pudieron cargar tus planes`);
  }
  return await res.json();
}
