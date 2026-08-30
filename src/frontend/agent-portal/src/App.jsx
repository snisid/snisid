import React, { useState } from 'react';
import { ShieldAlert, Users, FileCheck, Settings, Globe, AlertTriangle } from 'lucide-react';

function App() {
  const [isFrench, setIsFrench] = useState(true);

  const t = {
    logo: isFrench ? "Portail Opérationnel" : "Pòtay Operasyonèl",
    navEnroll: isFrench ? "Enrôlements" : "Anwòlman",
    navFraud: isFrench ? "Alertes Fraude" : "Alèt Fwod",
    navAudit: isFrench ? "Journaux d'Audit" : "Jounal Odit",
    title: isFrench ? "Gestion des Cas de Fraude" : "Jesyon Ka Fwod",
    stat1: isFrench ? "Alertes Actives" : "Alèt Aktif",
    stat2: isFrench ? "Anomalies Biométriques" : "Anomali Byometrik",
    stat3: isFrench ? "Tentatives API Bloquées" : "Tantativ API Bloke",
    thCase: isFrench ? "ID Cas" : "ID Ka",
    thType: isFrench ? "Type d'Alerte" : "Kalite Alèt",
    thSev: isFrench ? "Sévérité" : "Severite",
    thAct: isFrench ? "Action" : "Aksyon",
    btnInv: isFrench ? "Investiguer" : "Envestige"
  };

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <div className="sidebar-logo">
          <ShieldAlert size={28} color="#2563eb" />
          SNISID Ops
        </div>
        <nav className="nav-menu">
          <div className="nav-item"><Users size={20} /> {t.navEnroll}</div>
          <div className="nav-item active"><AlertTriangle size={20} /> {t.navFraud}</div>
          <div className="nav-item"><FileCheck size={20} /> {t.navAudit}</div>
          <div className="nav-item"><Settings size={20} /> Paramètres</div>
        </nav>
      </aside>

      <main className="main-content">
        <div className="topbar">
          <h1 className="page-title">{t.title}</h1>
          <button className="action-btn" onClick={() => setIsFrench(!isFrench)}>
            <Globe size={16} style={{ display:'inline', marginRight:'6px', verticalAlign:'text-bottom' }}/>
            {isFrench ? "Créole" : "Français"}
          </button>
        </div>

        <div className="dashboard-grid">
          <div className="stat-card">
            <div className="stat-title">{t.stat1}</div>
            <div className="stat-value danger">12</div>
          </div>
          <div className="stat-card">
            <div className="stat-title">{t.stat2}</div>
            <div className="stat-value warning">3</div>
          </div>
          <div className="stat-card">
            <div className="stat-title">{t.stat3} (OPA)</div>
            <div className="stat-value">4,892</div>
          </div>
        </div>

        <div className="table-container">
          <div className="table-header">
            Files d'attente d'Investigation
          </div>
          <table>
            <thead>
              <tr>
                <th>{t.thCase}</th>
                <th>{t.thType}</th>
                <th>{t.thSev}</th>
                <th>{t.thAct}</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td style={{ fontFamily: 'monospace' }}>CAS-9011</td>
                <td>Biométrie en Double (AFIS Match)</td>
                <td><span className="badge high">CRITIQUE</span></td>
                <td><button className="action-btn">{t.btnInv}</button></td>
              </tr>
              <tr>
                <td style={{ fontFamily: 'monospace' }}>CAS-9012</td>
                <td>Rate Limit API Dépassé (IP Suspecte)</td>
                <td><span className="badge medium">ÉLEVÉ</span></td>
                <td><button className="action-btn">{t.btnInv}</button></td>
              </tr>
              <tr>
                <td style={{ fontFamily: 'monospace' }}>CAS-9013</td>
                <td>Échec de Vérification d'Agence (DGI)</td>
                <td><span className="badge medium">ÉLEVÉ</span></td>
                <td><button className="action-btn">{t.btnInv}</button></td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
}

export default App;
