import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext'

interface ProtectedRouteProps {
  children: React.ReactNode; // The component to render if the user is authenticated
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  // Get userId from AuthContext
  const { userId, loading } = useAuth(); 
  if (loading) {
    return null; // Render nothing while determining authentication status
  }
  console.log('ProtectedRoute: userId is', userId);
  // If the user is not authenticated, redirect them to the login page
  if (!userId) {
    return <Navigate to="/login" replace />;
  }

  // If the user is authenticated, render the protected component
  return <>{children}</>;
};

export default ProtectedRoute;
