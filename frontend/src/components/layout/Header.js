import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuthentication } from "../../context/AuthContext";
import { LogOut, User, Users, Home, LogIn, UserPlus } from "lucide-react";

export function Header() {
  const { isUserAuthenticated, performUserLogout } = useAuthentication();
  const navigate = useNavigate();

  const handleLogout = () => {
    performUserLogout();
    navigate("/");
  };

  return (
    <header className="header">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo/Brand */}
          <div className="flex items-center">
            <Link
              to="/"
              className="text-2xl font-bold hero-gradient hover:opacity-80 transition-opacity"
            >
              ✨ Assetra
            </Link>
          </div>

          {/* Navigation */}
          <nav className="flex items-center space-x-8">
            <Link to="/" className="nav-link flex items-center">
              <Home className="w-4 h-4 nav-icon-space" />
              Home
            </Link>

            {isUserAuthenticated && (
              <>
                <Link to="/users" className="nav-link flex items-center">
                  <Users className="w-4 h-4 nav-icon-space" />
                  Users
                </Link>
                <Link to="/profile" className="nav-link flex items-center">
                  <User className="w-4 h-4 nav-icon-space" />
                  Profile
                </Link>
              </>
            )}
          </nav>

          {/* Authentication */}
          <div className="flex items-center space-x-8">
            {isUserAuthenticated ? (
              <button
                onClick={handleLogout}
                className="nav-link-gray flex items-center"
              >
                <LogOut className="w-4 h-4 nav-icon-space" />
                Logout
              </button>
            ) : (
              <>
                <Link to="/login" className="nav-link-gray flex items-center">
                  <LogIn className="w-4 h-4 nav-icon-space" />
                  Login
                </Link>
                <Link
                  to="/register"
                  className="nav-link-green flex items-center"
                >
                  <UserPlus className="w-4 h-4 nav-icon-space" />
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
