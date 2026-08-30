import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from 'react-oidc-context';
import { useSecurityStore } from '../../store/securityStore';

interface ProtectedRouteProps {
  requiredRoles?: string[];
  requiredPermissions?: string[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ 
  requiredRoles = [], 
  requiredPermissions = [] 
}) => {
  const auth = useAuth();
  const { hasRole, hasPermission } = useSecurityStore();

  if (auth.isLoading) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-slate-900 text-white" aria-live="polite">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        <span className="sr-only">Chargement de l'authentification...</span>
      </div>
    );
  }

  if (auth.error) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-slate-900 text-red-500">
        <h1>Erreur d'authentification : {auth.error.message}</h1>
      </div>
    );
  }

  if (!auth.isAuthenticated) {
    // Redirect to login (Keycloak handles the actual login page)
    auth.signinRedirect();
    return null; // Return null while redirecting
  }

  // RBAC checks
  const roleAuthorized = requiredRoles.length === 0 || requiredRoles.some(role => hasRole(role));
  const permAuthorized = requiredPermissions.length === 0 || requiredPermissions.some(perm => hasPermission(perm));

  if (!roleAuthorized || !permAuthorized) {
    return (
      <div className="flex flex-col h-screen w-screen items-center justify-center bg-slate-900 text-white">
        <h1 className="text-3xl font-bold text-red-500 mb-4">Accès Refusé</h1>
        <p>Vous n'avez pas les autorisations nécessaires pour voir cette page.</p>
        <button 
          onClick={() => window.history.back()}
          className="mt-6 px-4 py-2 bg-blue-600 rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Retour
        </button>
      </div>
    );
  }

  // Render child routes
  return <Outlet />;
};
