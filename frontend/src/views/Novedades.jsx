import React, { useState, useEffect } from 'react';
import './Novedades.css';

export default function Novedades() {
  const [planes, setPlanes] = useState([]);
  const [etiquetas, setEtiquetas] = useState({});
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setCargando(true);
        setError(null);

        // Cargar etiquetas para mapear idEtiqueta -> nombre
        const resEtiquetas = await fetch('/api/etiquetas');
        if (resEtiquetas.ok) {
          const dataEtiquetas = await resEtiquetas.json();
          const mapEtiquetas = {};
          dataEtiquetas.forEach((e) => {
            mapEtiquetas[e.id] = e.nombre;
          });
          setEtiquetas(mapEtiquetas);
        }

        // Cargar planes públicos
        const resPlanes = await fetch('/api/planes');
        if (!resPlanes.ok) {
          throw new Error(`Error ${resPlanes.status}: No se pudieron cargar los planes`);
        }
        const dataPlanes = await resPlanes.json();
        setPlanes(dataPlanes || []);
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

  return (
    <div className="novedades-container">
      <header className="novedades-header">
        <div className="novedades-header-content">
          <h1>🔥 Feed de Novedades</h1>
          <p>Descubrí los mejores planes públicos organizados por la comunidad.</p>
        </div>
      </header>

      {cargando && (
        <div className="status-message loading">
          <div className="spinner"></div>
          <p>Cargando planes públicos...</p>
        </div>
      )}

      {error && (
        <div className="status-message error">
          <p>⚠️ {error}</p>
        </div>
      )}

      {!cargando && !error && planes.length === 0 && (
        <div className="status-message empty">
          <p>No hay planes públicos disponibles en este momento.</p>
        </div>
      )}

      {!cargando && !error && planes.length > 0 && (
        <div className="planes-grid">
          {planes.map((plan) => (
            <article key={plan.id} className="plan-card">
              <div className="plan-card-header">
                <span className="badge badge-visibilidad">
                  🌐 {plan.visibilidad.toUpperCase()}
                </span>
                <span className="badge badge-etiqueta">
                  🏷️ {etiquetas[plan.idEtiqueta] || `Etiqueta #${plan.idEtiqueta}`}
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
                <button className="btn-unirse" onClick={() => alert(`¡Te uniste a: ${plan.nombre}!`)}>
                  Unirme al Plan
                </button>
              </div>
            </article>
          ))}
        </div>
      )}
    </div>
  );
}