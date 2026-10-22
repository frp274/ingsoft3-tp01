/**
 * Valida los datos requeridos para la creación de un nuevo plan.
 * @param {Object} formData 
 * @returns {{ valido: boolean, error: string | null }}
 */
export function validarPlan(formData) {
  if (!formData || typeof formData !== 'object') {
    return { valido: false, error: 'Los datos del plan son obligatorios.' };
  }

  const nombre = (formData.nombre || '').trim();
  if (!nombre) {
    return { valido: false, error: 'El nombre del plan es obligatorio.' };
  }

  if (!formData.idEtiqueta || Number(formData.idEtiqueta) <= 0) {
    return { valido: false, error: 'Debes seleccionar una etiqueta/categoría válida.' };
  }

  if (!formData.fechaInicio || !formData.fechaFin) {
    return { valido: false, error: 'Las fechas de inicio y fin son obligatorias.' };
  }

  const inicio = new Date(formData.fechaInicio);
  const fin = new Date(formData.fechaFin);

  if (isNaN(inicio.getTime()) || isNaN(fin.getTime())) {
    return { valido: false, error: 'El formato de las fechas es inválido.' };
  }

  if (fin <= inicio) {
    return { valido: false, error: 'La fecha y hora de fin debe ser posterior a la fecha de inicio.' };
  }

  const ubicacion = (formData.ubicacion || '').trim();
  if (!ubicacion) {
    return { valido: false, error: 'La ubicación es obligatoria.' };
  }

  const capacidadNum = parseInt(formData.capacidad, 10);
  if (isNaN(capacidadNum) || capacidadNum <= 0) {
    return { valido: false, error: 'La capacidad debe ser un número mayor a 0.' };
  }

  return { valido: true, error: null };
}

/**
 * Valida si un valor de capacidad es un entero positivo válido.
 * @param {any} capacidad 
 * @returns {boolean}
 */
export function validarCapacidad(capacidad) {
  if (capacidad === null || capacidad === undefined || capacidad === '') {
    return false;
  }
  const num = Number(capacidad);
  return Number.isInteger(num) && num > 0;
}
