import React, { createContext, useState, useContext, useEffect } from "react";
import authService from "../services/authService";
import { setupAuthService } from "../services/api";

const AuthContext = createContext({
  user: null,
  performUserLogin: async () => {},
  performUserRegistration: async () => {},
  performUserLogout: () => {},
});

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => authService.getCurrentUser());
  const [isInitializing, setIsInitializing] = useState(true);

  useEffect(() => {
    // Setup authService in api client for token management
    setupAuthService(authService);

    // On mount, check if we have user data but no access token in memory
    // This happens after page refresh (F5)
    const initializeAuth = async () => {
      const currentUser = authService.getCurrentUser();
      const token = authService.getAccessToken();

      if (currentUser && !token) {
        // User data exists but no token in memory - try to refresh silently
        try {
          await authService.refreshAccessToken();
          // Token refreshed successfully, user remains logged in
          console.log('[Auth] Session restored from refresh token');
        } catch (error) {
          // Refresh failed - this is expected if user has no valid refresh token
          // Clear everything silently (no error toast, this is normal on first visit)
          console.log('[Auth] No valid session, clearing user data');
          authService.clearAccessToken();
          localStorage.removeItem("currentUser");
          setUser(null);
        }
      }

      setIsInitializing(false);
    };

    initializeAuth();
  }, []);

  const performUserLogin = async (credentials) => {
    const { user: userData } = await authService.signIn(credentials);
    setUser(userData);
    return userData;
  };

  const performUserRegistration = async (userData) => {
    const { user: newUser } = await authService.signUp(userData);
    setUser(newUser);
    return newUser;
  };

  const performUserLogout = async () => {
    await authService.signOut();
    setUser(null);
  };

  // Don't render children until auth is initialized
  if (isInitializing) {
    return null; // Or a loading spinner
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isUserAuthenticated: !!user,
        performUserLogin,
        performUserRegistration,
        performUserLogout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

// Aliases for backward compatibility
export const useAuthentication = useAuth;
export const AuthenticationProvider = AuthProvider;

export default AuthContext;
