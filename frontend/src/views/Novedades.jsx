import React, { useState, useEffect } from 'react';
import './Novedades.css';

export default function Novedades() {
  const [planes, setPlanes] = useState([]);
  const [etiquetas, setEtiquetas] = useState({});
  const [solicitados, setSolicitados] = useState({});
  const [solicitandoId, setSolicitandoId] = useState(null);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setCargando(true);
        setError(null);

        // 1. Cargar etiquetas
        const resEtiquetas = await fetch('/api/etiquetas');
        if (resEtiquetas.ok) {
          const dataEtiquetas = await resEtiquetas.json();
          const mapEtiquetas = {};
          dataEtiquetas.forEach((e) => {
            mapEtiquetas[e.id] = e.nombre;
          });
          setEtiquetas(mapEtiquetas);
        }

        // 2. Cargar mis planes / invitaciones existentes del usuario activo
        const resMisPlanes = await fetch('/api/mis-planes');
        if (resMisPlanes.ok) {
          const misPlanesData = await resMisPlanes.json();
          const mapSolicitados = {};
          misPlanesData.forEach((p) => {
            mapSolicitados[p.id] = p.estado || 'pendiente';
          });
          setSolicitados(mapSolicitados);
        }

        // 3. Cargar planes públicos
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

  const handleSolicitarInvitacion = async (planId) => {
    try {
      setSolicitandoId(planId);
      const res = await fetch('/api/invitaciones', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ idPlan: planId }),
      });

      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || 'No se pudo solicitar la invitación');
      }

      const invData = await res.json();
      setSolicitados((prev) => ({
        ...prev,
        [planId]: invData.estado || 'pendiente',
      }));
    } catch (err) {
      alert(`⚠️ ${err.message}`);
    } finally {
      setSolicitandoId(null);
    }
  };

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
          {planes.map((plan) => {
            const estaSolicitado = Boolean(solicitados[plan.id]);
            const estadoInv = solicitados[plan.id];

            return (
              <article key={plan.id} className="plan-card">
                <div className="plan-card-header">
                  <span className="badge badge-visibilidad">
                    🌐 {plan.visibilidad.toUpperCase()}
                  </span>
                  <span className="badge badge-etiqueta">
                    🏷️ {etiquetas[plan.idEtiqueta] || `Etiqueta #${plan.idEtiqueta}`}
                  </span>
                  {estaSolicitado && (
                    <span className="badge badge-estado-inv">
                      ⏳ {estadoInv.toUpperCase()}
                    </span>
                  )}
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
                  <button
                    className={`btn-solicitar ${estaSolicitado ? 'solicitado' : ''}`}
                    onClick={() => handleSolicitarInvitacion(plan.id)}
                    disabled={estaSolicitado || solicitandoId === plan.id}
                  >
                    {solicitandoId === plan.id
                      ? 'Enviando solicitud...'
                      : estaSolicitado
                      ? '✓ Solicitado'
                      : 'Solicitar Invitación'}
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      )}
    </div>
  );
}