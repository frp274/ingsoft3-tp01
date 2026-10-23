/**
 * Utilidades para cálculo de prioridades y descuentos en planes.
 * (Lógica de negocio para planes)
 */

export function calcularPrioridadPlan(plan) {
  if (!plan) {
    return 'BAJA';
  }

  // Prioridad según visibilidad, capacidad y categoría
  if (plan.visibilidad === 'PUBLICO' && plan.capacidadMax > 50) {
    return 'CRITICA';
  } else if (plan.visibilidad === 'PUBLICO' && plan.capacidadMax > 20) {
    return 'ALTA';
  } else if (plan.visibilidad === 'PRIVADO' && plan.capacidadMax > 30) {
    return 'MEDIA-ALTA';
  } else if (plan.capacidadMax > 10) {
    return 'MEDIA';
  } else {
    return 'BAJA';
  }
}

export function calcularDescuentoPlan(plan, tipoUsuario) {
  if (!plan || !tipoUsuario) {
    return 0;
  }

  let descuento = 0;

  switch (tipoUsuario) {
    case 'VIP':
      if (plan.capacidadMax > 50) {
        descuento = 35;
      } else {
        descuento = 25;
      }
      break;
    case 'PREMIUM':
      if (plan.capacidadMax > 20) {
        descuento = 20;
      } else {
        descuento = 15;
      }
      break;
    case 'ESTUDIANTE':
      if (plan.visibilidad === 'PUBLICO') {
        descuento = 15;
      } else {
        descuento = 10;
      }
      break;
    case 'REGULAR':
      descuento = 5;
      break;
    default:
      descuento = 0;
      break;
  }

  return descuento;
}
