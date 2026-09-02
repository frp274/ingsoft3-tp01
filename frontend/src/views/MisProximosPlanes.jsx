import React, { useState, useEffect } from 'react';
import './MisProximosPlanes.css';

export default function MisProximosPlanes({ onExplorar }) {
  const [misPlanes, setMisPlanes] = useState([]);
  const [etiquetas, setEtiquetas] = useState({});
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setCargando(true);
        setError(null);

        // Cargar etiquetas
        const resEtiquetas = await fetch('/api/etiquetas');
        if (resEtiquetas.ok) {
          const dataEtiquetas = await resEtiquetas.json();
          const mapEtiquetas = {};
          dataEtiquetas.forEach((e) => {
            mapEtiquetas[e.id] = e.nombre;
          });
          setEtiquetas(mapEtiquetas);
        }

        // Cargar mis planes (JOIN Invitacion + Plan)
        const resMisPlanes = await fetch('/api/mis-planes');
        if (!resMisPlanes.ok) {
          throw new Error(`Error ${resMisPlanes.status}: No se pudieron cargar tus planes`);
        }
        const dataMisPlanes = await resMisPlanes.json();
        setMisPlanes(dataMisPlanes || []);
      } catch (err) {
        setError(err.message);
      } finally {
        setCargando(false);
      }
    };

    fetchData();
  }, []);

  const formatearFecha = (fechaStr) => {
    if (!fechaStr) return '';
    const fecha = new Date(fechaStr);
    return fecha.toLocaleDateString('es-AR', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getBadgeEstado = (estado) => {
    switch (estado?.toLowerCase()) {
      case 'aceptada':
        return <span className="badge badge-estado aceptada">✅ Confirmado</span>;
      case 'rechazada':
        return <span className="badge badge-estado rechazada">❌ Rechazada</span>;
      default:
        return <span className="badge badge-estado pendiente">⏳ Solicitud Pendiente</span>;
    }
  };

  return (
    <div className="mis-planes-container">
      <header className="mis-planes-header">
        <div className="mis-planes-header-content">
          <h1>📅 Mis Próximos Planes</h1>
          <p>Revisá los planes a los que te has postulado y el estado de tus invitaciones.</p>
        </div>
      </header>

      {cargando && (
        <div className="status-message loading">
          <div className="spinner"></div>
          <p>Cargando tus planes e invitaciones...</p>
        </div>
      )}

      {error && (
        <div className="status-message error">
          <p>⚠️ {error}</p>
        </div>
      )}

      {!cargando && !error && misPlanes.length === 0 && (
        <div className="status-message empty-state">
          <div className="empty-icon">🎟️</div>
          <h2>Todavía no te uniste a ningún plan</h2>
          <p>Explorá el feed de novedades para encontrar eventos y juntadas que te interesen.</p>
          <button className="btn-explorar" onClick={onExplorar}>
            Explorar Planes en Novedades
          </button>
        </div>
      )}

      {!cargando && !error && misPlanes.length > 0 && (
        <div className="mis-planes-grid">
          {misPlanes.map((plan) => (
            <article key={plan.idInvitacion || plan.id} className="plan-card-inscrito">
              <div className="plan-card-header">
                {getBadgeEstado(plan.estado)}
                <span className="badge badge-etiqueta">
                  🏷️ {etiquetas[plan.idEtiqueta] || `Etiqueta #${plan.idEtiqueta}`}
                </span>
                <span className="badge badge-visibilidad">
                  🌐 {plan.visibilidad.toUpperCase()}
                </span>
              </div>

              <h2 className="plan-titulo">{plan.nombre}</h2>

              {plan.descripcion && (
                <p className="plan-descripcion">{plan.descripcion}</p>
              )}

              <div className="plan-detalles">
                <div className="detalle-item">
                  <span className="detalle-icon">📅</span>
                  <div className="detalle-texto">
                    <strong>Inicio:</strong> {formatearFecha(plan.fechaInicio)}
                  </div>
                </div>

                <div className="detalle-item">
                  <span className="detalle-icon">🏁</span>
                  <div className="detalle-texto">
                    <strong>Fin:</strong> {formatearFecha(plan.fechaFin)}
                  </div>
                </div>

                <div className="detalle-item">
                  <span className="detalle-icon">📍</span>
                  <div className="detalle-texto">
                    <strong>Ubicación:</strong> {plan.ubicacion}
                  </div>
                </div>

                <div className="detalle-item">
                  <span className="detalle-icon">👥</span>
                  <div className="detalle-texto">
                    <strong>Capacidad:</strong> {plan.capacidad} personas
                  </div>
                </div>
              </div>

              <div className="plan-card-footer">
                <div className="estado-aviso">
                  {plan.estado === 'aceptada' ? (
                    <span className="aviso-ok">🎉 ¡Estás anotado! Te esperamos en el evento.</span>
                  ) : plan.estado === 'rechazada' ? (
                    <span className="aviso-rechazado">Tu solicitud no pudo ser aceptada por el organizador.</span>
                  ) : (
                    <span className="aviso-pendiente">El organizador revisará tu solicitud pronto.</span>
                  )}
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </div>
  );
}