import React, { useState } from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { useForm } from "react-hook-form";
import { useAuthentication } from "../context/AuthContext";
import { Layout } from "../components/layout/Layout";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { Card, CardContent, CardHeader } from "../components/ui/Card";
import { LogIn, Eye, EyeOff } from "lucide-react";
import { isValidEmail } from "../utils/validation";
import { toast } from "react-hot-toast";

export function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const { performUserLogin } = useAuthentication();
  const navigate = useNavigate();
  const location = useLocation();

  const redirectPath = location.state?.from?.pathname || "/";

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm();

  const onSubmit = async (formData) => {
    try {
      const loginResult = await performUserLogin({
        email: formData.email.trim().toLowerCase(),
        password: formData.password,
      });

      if (loginResult) {
        navigate(redirectPath, { replace: true });
      }
      // If login fails, the error already appears as toast in AuthContext
    } catch (error) {
      // Unexpected error will also appear as toast if necessary
      toast.error(error.message || "Login error");
    }
  };

  return (
    <Layout>
      <div className="login-center-container">
        <Card>
          <CardHeader>
            <div className="text-center space-y-2">
              <div className="flex justify-center">
                <LogIn className="w-8 h-8 text-blue-600" />
              </div>
              <h1 className="text-2xl font-bold text-gray-900">Sign In</h1>
              <p className="text-gray-600">
                Enter your credentials to access your account
              </p>
            </div>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <Input
                label="Email Address"
                type="email"
                autoComplete="email"
                {...register("email", {
                  required: "Email address is required",
                  validate: (value) =>
                    isValidEmail(value) || "Please enter a valid email address",
                  maxLength: {
                    value: 100,
                    message: "Email is too long",
                  },
                })}
                error={errors.email?.message}
                placeholder="Enter your email address"
              />

              <Input
                label="Password"
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                {...register("password", {
                  required: "Password is required",
                  minLength: {
                    value: 8,
                    message: "Password must be at least 8 characters",
                  },
                  maxLength: {
                    value: 128,
                    message: "Password is too long",
                  },
                })}
                error={errors.password?.message}
                placeholder="Enter your password"
              />
              <button
                type="button"
                className="password-eye-btn"
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? (
                  <EyeOff className="w-5 h-5" />
                ) : (
                  <Eye className="w-5 h-5" />
                )}
              </button>

              <Button type="submit" className="w-full btn-primary">
                Sign In
              </Button>
            </form>

            <div className="mt-6 text-center space-y-2">
              <p className="text-sm text-gray-600">
                Don't have an account?{" "}
                <Link
                  to="/register"
                  className="text-blue-600 hover:text-blue-500 font-medium"
                >
                  Sign up here
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
