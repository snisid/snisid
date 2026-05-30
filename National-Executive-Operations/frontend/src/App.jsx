import React, { useState, useEffect } from 'react';

const App = () => {
  const [documents, setDocuments] = useState([]);
  const [pinCode, setPinCode] = useState('');
  const [selectedDoc, setSelectedDoc] = useState(null);
  const [message, setMessage] = useState('');

  const fetchDocuments = async () => {
    try {
      const res = await fetch('/api/documents');
      const data = await res.json();
      setDocuments(data || []);
    } catch (e) {
      console.error('Error fetching documents', e);
      // Fallback for visual testing if API is unreachable
      setDocuments([{ id: 'DOC-MOCK-1', title: 'Arrêté Présidentiel (Mock UI)', status: 'PENDING_SIG' }]);
    }
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  const handleSign = async () => {
    if (!pinCode) {
      setMessage('Veuillez entrer votre code PIN Smartcard.');
      return;
    }
    
    try {
      const res = await fetch('/api/sign', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          document_id: selectedDoc.id,
          signer_id: 'MIN-001',
          pin_code: pinCode
        })
      });

      if (res.ok) {
        setMessage('✅ Document signé et scellé avec succès.');
        setSelectedDoc(null);
        setPinCode('');
        fetchDocuments();
      } else {
        setMessage('❌ Erreur de signature (Code PIN invalide ?)');
      }
    } catch (e) {
      setMessage('❌ Erreur de connexion au service PKI.');
    }
  };

  return (
    <div style={{ maxWidth: '800px', margin: '40px auto', padding: '20px', background: '#fff', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)' }}>
      <header style={{ borderBottom: '2px solid #0056b3', paddingBottom: '10px', marginBottom: '20px' }}>
        <h1 style={{ color: '#0056b3', margin: 0 }}>Parapheur Électronique National</h1>
        <p style={{ color: '#666', margin: '5px 0 0 0' }}>Espace Sécurisé - Hautes Instances de l'État</p>
      </header>

      {message && (
        <div style={{ padding: '10px', marginBottom: '20px', background: message.includes('✅') ? '#d4edda' : '#f8d7da', borderRadius: '4px' }}>
          {message}
        </div>
      )}

      {!selectedDoc ? (
        <div>
          <h3>Documents en attente de signature</h3>
          {documents.length === 0 && <p>Aucun document en attente.</p>}
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {documents.map(doc => (
              <li key={doc.id} style={{ display: 'flex', justifyContent: 'space-between', padding: '15px', border: '1px solid #ddd', marginBottom: '10px', borderRadius: '4px' }}>
                <div>
                  <strong>{doc.title}</strong>
                  <div style={{ fontSize: '0.85em', color: '#666' }}>ID: {doc.id} | Statut: {doc.status}</div>
                </div>
                {doc.status === 'PENDING_SIG' ? (
                  <button onClick={() => setSelectedDoc(doc)} style={{ background: '#0056b3', color: 'white', border: 'none', padding: '10px 15px', borderRadius: '4px', cursor: 'pointer' }}>
                    Ouvrir pour Signature
                  </button>
                ) : (
                  <span style={{ color: 'green', fontWeight: 'bold', padding: '10px' }}>SIGNÉ ✓</span>
                )}
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <div style={{ border: '1px solid #ccc', padding: '20px', borderRadius: '8px', background: '#fafafa' }}>
          <h3>Signature Cryptographique (Smartcard)</h3>
          <p>Vous êtes sur le point de sceller électroniquement le document <strong>{selectedDoc.title}</strong>.</p>
          <div style={{ marginBottom: '15px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Code PIN (Simulé : 1234) :</label>
            <input 
              type="password" 
              value={pinCode} 
              onChange={e => setPinCode(e.target.value)} 
              style={{ padding: '10px', width: '100%', boxSizing: 'border-box' }}
              placeholder="Saisissez le code PIN de la Smartcard..."
            />
          </div>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button onClick={handleSign} style={{ background: '#28a745', color: 'white', border: 'none', padding: '10px 20px', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold' }}>
              Valider et Sceller (QES)
            </button>
            <button onClick={() => setSelectedDoc(null)} style={{ background: '#6c757d', color: 'white', border: 'none', padding: '10px 20px', borderRadius: '4px', cursor: 'pointer' }}>
              Annuler
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default App;
