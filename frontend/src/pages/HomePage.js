import React from "react";
import { Link } from "react-router-dom";
import { useAuthentication } from "../context/AuthContext";
import { Layout } from "../components/layout/Layout";
import { Button } from "../components/ui/Button";
import { Card, CardContent } from "../components/ui/Card";
import { Users, Shield, Zap, Globe } from "lucide-react";

export function HomePage() {
  const { isUserAuthenticated, user } = useAuthentication();

  const platformFeatures = [
    {
      icon: <Users className="w-8 h-8 text-gray-dark" />,
      title: "User Management",
      description: "Complete user management system with access control.",
    },
    {
      icon: <Shield className="w-8 h-8 text-gray-dark" />,
      title: "Advanced Security",
      description: "JWT authentication and encryption for maximum security.",
    },
    {
      icon: <Zap className="w-8 h-8 text-gray-dark" />,
      title: "High Performance",
      description: "Fast and efficient API built with Go for high performance.",
    },
    {
      icon: <Globe className="w-8 h-8 text-gray-dark" />,
      title: "Scalable",
      description: "Architecture ready for Kubernetes and cloud environments.",
    },
  ];

  return (
    <Layout>
      <div className="space-y-12">
        {/* Hero Section */}
        <div className="text-center space-y-6">
          <h1 className="text-4xl md:text-6xl font-bold text-gray-dark">
            Welcome to <span className="hero-gradient">Assetra</span>
          </h1>
          <p className="text-xl text-gray max-w-3xl mx-auto">
            A modern and secure platform for user management and authentication,
            built with development best practices.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            {isUserAuthenticated ? (
              <div className="space-y-4 text-center">
                <p className="text-lg text-gray-700">
                  Hello,{" "}
                  <span className="font-semibold">
                    {user?.name || user?.username}
                  </span>
                  ! You are logged in.
                </p>
              </div>
            ) : (
              <div className="space-y-6">
                <div className="flex gap-4 justify-center items-center">
                  <Link to="/login">
                    <Button
                      size="lg"
                      className="btn-primary px-8 py-3 text-base font-semibold rounded-xl shadow-sm"
                    >
                      Sign In
                    </Button>
                  </Link>
                  <Link to="/register">
                    <Button
                      variant="outline"
                      size="lg"
                      className="btn-outline px-8 py-3 text-base font-semibold rounded-xl shadow-sm"
                    >
                      Create Account
                    </Button>
                  </Link>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Features Section */}
        <div className="space-y-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-dark mb-4">Features</h2>
            <p className="text-lg text-gray">
              Discover the main characteristics of our platform
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {platformFeatures.map((feature, index) => (
              <Card
                key={index}
                className="text-center hover:shadow-lg transition-shadow"
              >
                <CardContent className="p-6">
                  <div className="feature-icon icon-bounce flex justify-center mb-4">
                    {feature.icon}
                  </div>
                  <h3 className="text-lg font-semibold text-gray-dark mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-gray">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Stats Section */}
        <div className="card stats-card bg-gradient text-green-800 p-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
            <div>
              <div className="text-3xl font-bold mb-2">100%</div>
              <div className="opacity-75">Secure</div>
            </div>
            <div>
              <div className="text-3xl font-bold mb-2">24/7</div>
              <div className="opacity-75">Availability</div>
            </div>
            <div>
              <div className="text-3xl font-bold mb-2">∞</div>
              <div className="opacity-75">Scalability</div>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}
