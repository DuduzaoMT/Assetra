import React, { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import { AuthenticationProvider } from "./context/AuthContext";
import { LoadingProvider, useLoading } from "./context/LoadingContext";
import { ProtectedRoute } from "./components/ProtectedRoute";
import GlobalLoading from "./components/ui/GlobalLoading";
import { setupLoadingInterceptors } from "./services/api";
import "./App.css";

// Application Pages
import { HomePage } from "./pages/HomePage";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { UsersPage } from "./pages/UsersPage";
import { ProfilePage } from "./pages/ProfilePage";

// Component to setup loading interceptors
function LoadingSetup() {
  const loadingContext = useLoading();

  useEffect(() => {
    setupLoadingInterceptors(loadingContext);
  }, [loadingContext]);

  return null;
}

function App() {
  return (
    <LoadingProvider>
      <AuthenticationProvider>
        <Router>
          <div className="App">
            <LoadingSetup />
            <Routes>
              {/* Public Routes - accessible without authentication */}
              <Route path="/" element={<HomePage />} />
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />

              {/* Protected Routes - require authentication */}
              <Route
                path="/users"
                element={
                  <ProtectedRoute>
                    <UsersPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/profile"
                element={
                  <ProtectedRoute>
                    <ProfilePage />
                  </ProtectedRoute>
                }
              />
            </Routes>

            {/* Global Loading Overlay */}
            <GlobalLoading />

            {/* Global Toast Notifications */}
            <Toaster
              position="top-right"
              toastOptions={{
                duration: 4000,
                style: {
                  background: "#363636",
                  color: "#fff",
                },
                success: {
                  duration: 3000,
                  theme: {
                    primary: "green",
                    secondary: "black",
                  },
                },
              }}
            />
          </div>
        </Router>
      </AuthenticationProvider>
    </LoadingProvider>
  );
}

export default App;
