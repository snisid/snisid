import React, { useState } from 'react';

function App() {
  const [status, setStatus] = useState('Pending');

  const signDocument = async () => {
    try {
      const res = await fetch('/api/sign', { method: 'POST' });
      const data = await res.json();
      setStatus(data.status);
    } catch (e) {
      console.error(e);
      setStatus('Error signing');
    }
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial' }}>
      <h1>Parapheur Électronique - Visa Circuit</h1>
      <div style={{ border: '1px solid #ccc', padding: '20px', marginTop: '20px' }}>
        <h3>Document: Arrêté Présidentiel</h3>
        <p>ID: DOC-12345</p>
        <p>Status: <strong>{status}</strong></p>
        <button 
          onClick={signDocument}
          style={{ padding: '10px 20px', backgroundColor: '#00205B', color: 'white', border: 'none', cursor: 'pointer' }}
        >
          Sign with Smartcard (PKI)
        </button>
      </div>
    </div>
  );
}

export default App;
