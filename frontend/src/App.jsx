import React, { useState } from 'react';
import Novedades from './views/Novedades';
import CrearPlan from './views/CrearPlan';
import MisProximosPlanes from './views/MisProximosPlanes';
import './App.css';

function App() {
  const [vistaActual, setVistaActual] = useState('novedades');

  return (
    <div className="app-layout">
      <nav className="navbar">
        <div className="nav-content">
          <div className="nav-brand" onClick={() => setVistaActual('novedades')}>
            <span className="brand-logo">🎉</span>
            <span className="brand-name">PlanesApp</span>
          </div>

          <div className="nav-links">
            <button
              className={`nav-btn ${vistaActual === 'novedades' ? 'active' : ''}`}
              onClick={() => setVistaActual('novedades')}
            >
              Novedades
            </button>
            <button
              className={`nav-btn ${vistaActual === 'mis-planes' ? 'active' : ''}`}
              onClick={() => setVistaActual('mis-planes')}
            >
              Mis Próximos Planes
            </button>
            <button
              className={`nav-btn btn-crear ${vistaActual === 'crear-plan' ? 'active' : ''}`}
              onClick={() => setVistaActual('crear-plan')}
            >
              + Crear Plan
            </button>
          </div>
        </div>
      </nav>

      <main className="main-content">
        {vistaActual === 'novedades' && <Novedades />}
        {vistaActual === 'mis-planes' && (
          <MisProximosPlanes onExplorar={() => setVistaActual('novedades')} />
        )}
        {vistaActual === 'crear-plan' && (
          <CrearPlan onPlanCreado={() => setVistaActual('novedades')} />
        )}
      </main>
    </div>
  );
}

export default App;