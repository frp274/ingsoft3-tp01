/**
 * Formatea una fecha ISO o string a formato legible en español.
 * Si la fecha es nula, indefinida o inválida, retorna un valor por defecto o cadena vacía.
 * @param {string|Date|null|undefined} fechaStr 
 * @returns {string}
 */
export function formatearFecha(fechaStr) {
  if (!fechaStr) return '';
  const fecha = new Date(fechaStr);
  if (isNaN(fecha.getTime())) {
    return 'Fecha inválida';
  }
  return fecha.toLocaleDateString('es-AR', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * Normaliza y devuelve los atributos visuales para mostrar el estado de una invitación.
 * @param {string|null|undefined} estado 
 * @returns {{ etiqueta: string, clase: string }}
 */
export function formatearEstadoInvitacion(estado) {
  const normalizado = (estado || '').toLowerCase().trim();
  switch (normalizado) {
    case 'aceptada':
      return { etiqueta: '✅ Confirmado', clase: 'aceptada' };
    case 'rechazada':
      return { etiqueta: '❌ Rechazada', clase: 'rechazada' };
    case 'pendiente':
      return { etiqueta: '⏳ Solicitud Pendiente', clase: 'pendiente' };
    default:
      return { etiqueta: 'Desconocido', clase: 'desconocido' };
  }
}

/**
 * Calcula la duración en horas entre dos fechas.
 * @param {string|Date} fechaInicio 
 * @param {string|Date} fechaFin 
 * @returns {number} Duración en horas redondeada a 1 decimal
 */
export function calcularDuracionHoras(fechaInicio, fechaFin) {
  if (!fechaInicio || !fechaFin) return 0;
  const inicio = new Date(fechaInicio);
  const fin = new Date(fechaFin);
  if (isNaN(inicio.getTime()) || isNaN(fin.getTime()) || fin <= inicio) {
    return 0;
  }
  const diffMs = fin.getTime() - inicio.getTime();
  return Math.round((diffMs / (1000 * 60 * 60)) * 10) / 10;
}
