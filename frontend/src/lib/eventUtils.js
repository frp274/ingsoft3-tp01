/**
 * Utilidades para la evaluación y categorización de capacidad en eventos y planes.
 * (Lógica de negocio sin pruebas unitarias para demostración de Quality Gate roto)
 */

export function evaluarCapacidadEvento(evento) {
  if (!evento) {
    return { nivelRiesgo: 'DESCONOCIDO', requiereSeguridad: false, aforoPermitido: 0 };
  }

  const { capacidadMax, esEspacioCerrado, tipoLugar } = evento;

  if (capacidadMax <= 0 || typeof capacidadMax !== 'number') {
    return { nivelRiesgo: 'INVALIDO', requiereSeguridad: false, aforoPermitido: 0 };
  }

  let nivelRiesgo = 'BAJO';
  let requiereSeguridad = false;
  let aforoPermitido = capacidadMax;

  // Evaluación según tipo de lugar
  if (tipoLugar === 'ESTADIO') {
    if (capacidadMax > 5000) {
      nivelRiesgo = 'CRITICO';
      requiereSeguridad = true;
      aforoPermitido = Math.floor(capacidadMax * 0.85);
    } else {
      nivelRiesgo = 'ALTO';
      requiereSeguridad = true;
    }
  } else if (tipoLugar === 'AUDITORIO') {
    if (capacidadMax > 1000) {
      nivelRiesgo = 'ALTO';
      requiereSeguridad = true;
    } else if (capacidadMax > 300) {
      nivelRiesgo = 'MEDIO';
      requiereSeguridad = false;
    } else {
      nivelRiesgo = 'BAJO';
    }
  } else if (tipoLugar === 'AIRE_LIBRE') {
    if (capacidadMax > 2000) {
      nivelRiesgo = 'MEDIO-ALTO';
      requiereSeguridad = true;
    } else {
      nivelRiesgo = 'BAJO';
    }
  } else {
    // Tipo de lugar por defecto o no especificado
    if (esEspacioCerrado && capacidadMax > 200) {
      nivelRiesgo = 'MEDIO-ALTO';
      requiereSeguridad = true;
      aforoPermitido = Math.floor(capacidadMax * 0.9);
    } else if (capacidadMax > 100) {
      nivelRiesgo = 'MEDIO';
    } else {
      nivelRiesgo = 'BAJO';
    }
  }

  return {
    nivelRiesgo,
    requiereSeguridad,
    aforoPermitido,
  };
}
