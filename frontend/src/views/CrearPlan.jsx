import React, { useState, useEffect } from 'react';
import './CrearPlan.css';

export default function CrearPlan({ onPlanCreado }) {
  const [formData, setFormData] = useState({
    nombre: '',
    idEtiqueta: '',
    fechaInicio: '',
    fechaFin: '',
    ubicacion: '',
    capacidad: '',
    visibilidad: 'publico',
    descripcion: '',
  });

  const [etiquetas, setEtiquetas] = useState([]);
  const [cargandoEtiquetas, setCargandoEtiquetas] = useState(true);
  const [enviando, setEnviando] = useState(false);
  const [error, setError] = useState(null);
  const [exito, setExito] = useState(false);

  useEffect(() => {
    const fetchEtiquetas = async () => {
      try {
        const res = await fetch('/api/etiquetas');
        if (!res.ok) throw new Error('No se pudieron cargar las etiquetas');
        const data = await res.json();
        setEtiquetas(data || []);
        if (data && data.length > 0) {
          setFormData((prev) => ({ ...prev, idEtiqueta: data[0].id }));
        }
      } catch (err) {
        setError('Error al obtener categorías de etiquetas: ' + err.message);
      } finally {
        setCargandoEtiquetas(false);
      }
    };

    fetchEtiquetas();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);

    // Validaciones en Frontend
    if (!formData.nombre.trim()) {
      setError('El nombre del plan es obligatorio.');
      return;
    }
    if (!formData.idEtiqueta) {
      setError('Debes seleccionar una etiqueta/categoría.');
      return;
    }
    if (!formData.fechaInicio || !formData.fechaFin) {
      setError('Las fechas de inicio y fin son obligatorias.');
      return;
    }

    const inicio = new Date(formData.fechaInicio);
    const fin = new Date(formData.fechaFin);
    if (fin <= inicio) {
      setError('La fecha y hora de fin debe ser posterior a la fecha de inicio.');
      return;
    }

    if (!formData.ubicacion.trim()) {
      setError('La ubicación es obligatoria.');
      return;
    }

    const capacidadNum = parseInt(formData.capacidad, 10);
    if (isNaN(capacidadNum) || capacidadNum <= 0) {
      setError('La capacidad debe ser un número mayor a 0.');
      return;
    }

    try {
      setEnviando(true);

      const payload = {
        nombre: formData.nombre.trim(),
        idEtiqueta: parseInt(formData.idEtiqueta, 10),
        fechaInicio: formData.fechaInicio,
        fechaFin: formData.fechaFin,
        visibilidad: formData.visibilidad,
        ubicacion: formData.ubicacion.trim(),
        capacidad: capacidadNum,
        descripcion: formData.descripcion.trim() ? formData.descripcion.trim() : null,
      };

      const res = await fetch('/api/planes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (res.status === 201) {
        setExito(true);
        setTimeout(() => {
          if (onPlanCreado) {
            onPlanCreado();
          }
        }, 800);
      } else {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error || `Error ${res.status} al crear el plan`);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  };

  return (
    <div className="crear-plan-container">
      <div className="crear-plan-card">
        <header className="crear-plan-header">
          <h2>✨ Crear Nuevo Plan Público</h2>
          <p>Completá los datos para publicar tu evento y sumar a la comunidad.</p>
        </header>

        {error && (
          <div className="alert-box alert-error">
            <span>⚠️ {error}</span>
          </div>
        )}

        {exito && (
          <div className="alert-box alert-success">
            <span>🎉 ¡Plan creado exitosamente! Redirigiendo a Novedades...</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="crear-plan-form">
          <div className="form-group">
            <label htmlFor="nombre">
              Nombre del Plan <span className="req">*</span>
            </label>
            <input
              type="text"
              id="nombre"
              name="nombre"
              value={formData.nombre}
              onChange={handleChange}
              placeholder="Ej: Juntada de Cine Clásico al Aire Libre"
              required
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="idEtiqueta">
                Etiqueta / Categoría <span className="req">*</span>
              </label>
              <select
                id="idEtiqueta"
                name="idEtiqueta"
                value={formData.idEtiqueta}
                onChange={handleChange}
                disabled={cargandoEtiquetas}
                required
              >
                {cargandoEtiquetas ? (
                  <option value="">Cargando etiquetas...</option>
                ) : (
                  etiquetas.map((e) => (
                    <option key={e.id} value={e.id}>
                      {e.nombre}
                    </option>
                  ))
                )}
              </select>
            </div>

            <div className="form-group">
              <label htmlFor="visibilidad">
                Visibilidad <span className="req">*</span>
              </label>
              <select
                id="visibilidad"
                name="visibilidad"
                value={formData.visibilidad}
                onChange={handleChange}
                required
              >
                <option value="publico">Público (visible en Novedades)</option>
                <option value="privado">Privado (solo por invitación)</option>
              </select>
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="fechaInicio">
                Fecha y Hora de Inicio <span className="req">*</span>
              </label>
              <input
                type="datetime-local"
                id="fechaInicio"
                name="fechaInicio"
                value={formData.fechaInicio}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="fechaFin">
                Fecha y Hora de Fin <span className="req">*</span>
              </label>
              <input
                type="datetime-local"
                id="fechaFin"
                name="fechaFin"
                value={formData.fechaFin}
                onChange={handleChange}
                required
              />
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="ubicacion">
                Ubicación <span className="req">*</span>
              </label>
              <input
                type="text"
                id="ubicacion"
                name="ubicacion"
                value={formData.ubicacion}
                onChange={handleChange}
                placeholder="Ej: Parque Sarmiento, Córdoba"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="capacidad">
                Capacidad (personas) <span className="req">*</span>
              </label>
              <input
                type="number"
                id="capacidad"
                name="capacidad"
                min="1"
                value={formData.capacidad}
                onChange={handleChange}
                placeholder="Ej: 20"
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="descripcion">
              Descripción <span className="opcional">(opcional)</span>
            </label>
            <textarea
              id="descripcion"
              name="descripcion"
              rows="4"
              value={formData.descripcion}
              onChange={handleChange}
              placeholder="Contanos más detalles del plan: qué llevar, punto de encuentro, etc."
            ></textarea>
          </div>

          <div className="form-actions">
            <button
              type="submit"
              className="btn-submit"
              disabled={enviando || exito}
            >
              {enviando ? 'Publicando plan...' : 'Publicar Plan'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}