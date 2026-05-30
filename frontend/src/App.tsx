import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { AppLayout } from './components/layout/AppLayout';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { RealTimeMetrics } from './pages/Dashboard/RealTimeMetrics';
import { CaptureInterface } from './components/biometrics/CaptureInterface';
import { IdentitiesPage } from './pages/IdentitiesPage';
import { AuditPage } from './pages/AuditPage';
import { GlossaryPage } from './pages/GlossaryPage';

// Biometric page wrapper
const BiometricPage = () => (
  <div className="flex flex-col items-center justify-center h-full">
    <CaptureInterface onCapture={(blob) => console.log('Captured', blob)} />
  </div>
);

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Protected Routes Wrapper */}
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout><Outlet /></AppLayout>}>
            {/* Dashboard — accessible à tous les agents connectés */}
            <Route path="/" element={<RealTimeMetrics />} />

            {/* Identités — RBAC géré côté API */}
            <Route path="/identities" element={<IdentitiesPage />} />

            {/* Biométrie — opérateur uniquement */}
            <Route element={<ProtectedRoute requiredRoles={['OPERATOR']} />}>
              <Route path="/biometrics" element={<BiometricPage />} />
            </Route>

            {/* Audit & Investigation forensique — auditeur uniquement */}
            <Route element={<ProtectedRoute requiredRoles={['AUDITOR']} />}>
              <Route path="/audit" element={<AuditPage />} />
            </Route>

            {/* Glossaire Technique — accessible à tous */}
            <Route path="/glossary" element={<GlossaryPage />} />
          </Route>
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;

