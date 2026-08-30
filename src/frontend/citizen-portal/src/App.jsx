import React, { useState } from 'react';
import { Shield, Fingerprint, Lock, Globe } from 'lucide-react';

function App() {
  const [isFrench, setIsFrench] = useState(true);
  const [consentGranted, setConsentGranted] = useState(true);

  const t = {
    title: isFrench ? "Portail Citoyen SNISID" : "Pòtay Sitwayen SNISID",
    lang: isFrench ? "Passer en Créole" : "Chanje an Fransè",
    walletTitle: isFrench ? "Portefeuille d'Identité" : "Bous Idantite",
    status: isFrench ? "Actif & Vérifié" : "Aktif & Verifye",
    name: isFrench ? "Nom Complet" : "Non Konplè",
    uid: isFrench ? "Identifiant Unique (NIN)" : "Idantifyan Inik (NIN)",
    biometrics: isFrench ? "Empreinte Biométrique" : "Anprent Byometrik",
    consentTitle: isFrench ? "Gouvernance des Données" : "Gouvènans Done",
    consentDesc: isFrench 
      ? "Contrôlez cryptographiquement quelles agences gouvernementales peuvent accéder à vos données." 
      : "Kontwole kriptografikman ki ajans gouvènman ki ka gen aksè a done ou yo.",
    revokeBtn: isFrench ? "Révoquer l'Accès (DGI)" : "Revoke Aksè (DGI)",
    grantBtn: isFrench ? "Autoriser l'Accès (DGI)" : "Otorize Aksè (DGI)"
  };

  return (
    <div className="app-container">
      <div className="glass-panel">
        <header>
          <h1>{t.title}</h1>
          <button 
            className="toggle-btn" 
            style={{ width: 'auto', marginTop: 0, background: 'rgba(255,255,255,0.1)' }}
            onClick={() => setIsFrench(!isFrench)}
          >
            <Globe size={18} style={{ display: 'inline', marginRight: '8px', verticalAlign: 'middle' }}/>
            {t.lang}
          </button>
        </header>

        <div className="dashboard-grid">
          {/* Identity Wallet Card */}
          <div className="card">
            <h2><Fingerprint size={24} color="#3b82f6"/> {t.walletTitle}</h2>
            <div style={{ marginBottom: '2rem' }}>
              <span className="status-badge">{t.status}</span>
            </div>
            
            <div className="data-row">
              <span className="label">{t.name}</span>
              <span className="value">Jean-Pierre Baptiste</span>
            </div>
            <div className="data-row">
              <span className="label">{t.uid}</span>
              <span className="value" style={{ fontFamily: 'monospace' }}>HT-8492-4911-30X</span>
            </div>
            <div className="data-row">
              <span className="label">{t.biometrics}</span>
              <span className="value">SHA256: e3b0c442...</span>
            </div>
          </div>

          {/* Consent Governance Card */}
          <div className="card">
            <h2><Shield size={24} color="#10b981"/> {t.consentTitle}</h2>
            <p style={{ color: 'var(--text-muted)', lineHeight: '1.6', marginBottom: '1.5rem' }}>
              {t.consentDesc}
            </p>
            
            <div style={{ padding: '1.5rem', background: 'rgba(0,0,0,0.2)', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.02)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ fontWeight: '600' }}><Lock size={16} style={{ display:'inline', marginRight:'8px' }}/> DGI (Impôts)</span>
                <span style={{ color: consentGranted ? 'var(--success)' : 'var(--danger)', fontSize: '0.9rem', fontWeight: 'bold' }}>
                  {consentGranted ? 'AUTORISÉ' : 'RÉVOQUÉ'}
                </span>
              </div>
              <button 
                className={`toggle-btn ${!consentGranted ? 'granted' : ''}`}
                onClick={() => setConsentGranted(!consentGranted)}
              >
                {consentGranted ? t.revokeBtn : t.grantBtn}
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}

export default App;
