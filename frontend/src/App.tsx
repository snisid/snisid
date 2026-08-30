import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { AppLayout } from './components/layout/AppLayout';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { RealTimeMetrics } from './pages/Dashboard/RealTimeMetrics';
import { CaptureInterface } from './components/biometrics/CaptureInterface';
import { IdentitiesPage } from './pages/IdentitiesPage';
import { AuditPage } from './pages/AuditPage';
import { GlossaryPage } from './pages/GlossaryPage';
import SigintDashboard from './pages/SigintDashboard';
import HumintDashboard from './pages/HumintDashboard';
import AirDefenseCOP from './pages/AirDefenseCOP';
import MilC2Dashboard from './pages/MilC2Dashboard';
import BioSurveillanceDashboard from './pages/BioSurveillanceDashboard';
import ExecProtectionDashboard from './pages/ExecProtectionDashboard';
import TransportSecurityDashboard from './pages/TransportSecurityDashboard';
import RadiationSafetyDashboard from './pages/RadiationSafetyDashboard';
import FusionCenterDashboard from './pages/FusionCenterDashboard';
import CounterintelDashboard from './pages/CounterintelDashboard';
import CriticalInfraDashboard from './pages/CriticalInfraDashboard';
import FisaCourtDashboard from './pages/FisaCourtDashboard';
import ClassificationDashboard from './pages/ClassificationDashboard';

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

            {/* National Security Dashboards */}
            <Route path="/dashboard/sigint" element={<SigintDashboard />} />
            <Route path="/dashboard/humint" element={<HumintDashboard />} />
            <Route path="/dashboard/air-defense" element={<AirDefenseCOP />} />
            <Route path="/dashboard/mil-c2" element={<MilC2Dashboard />} />
            <Route path="/dashboard/bio-surveillance" element={<BioSurveillanceDashboard />} />
            <Route path="/dashboard/exec-protection" element={<ExecProtectionDashboard />} />
            <Route path="/dashboard/transport-security" element={<TransportSecurityDashboard />} />
            <Route path="/dashboard/radiation-safety" element={<RadiationSafetyDashboard />} />
            <Route path="/dashboard/fusion-center" element={<FusionCenterDashboard />} />
            <Route path="/dashboard/counterintel" element={<CounterintelDashboard />} />
            <Route path="/dashboard/critical-infrastructure" element={<CriticalInfraDashboard />} />
            <Route path="/dashboard/fisa-court" element={<FisaCourtDashboard />} />
            <Route path="/dashboard/classification" element={<ClassificationDashboard />} />

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

