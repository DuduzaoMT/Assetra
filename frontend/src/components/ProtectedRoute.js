import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuthentication } from "../context/AuthContext";

export function ProtectedRoute({ children }) {
  const { user } = useAuthentication();
  const currentLocation = useLocation();

  if (!user) {
    // Save current location to redirect back after login
    return <Navigate to="/login" state={{ from: currentLocation }} replace />;
  }

  return children;
}
