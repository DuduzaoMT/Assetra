import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { useAuthentication } from "../context/AuthContext";
import { Layout } from "../components/layout/Layout";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import { Card, CardContent, CardHeader } from "../components/ui/Card";
import { UserPlus, Eye, EyeOff } from "lucide-react";
import {
  isValidEmail,
  validatePassword,
  validateName,
} from "../utils/validation";
import { toast } from "react-hot-toast";

export function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { performUserRegistration } = useAuthentication();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm();

  const passwordValue = watch("password");

  const onSubmit = async (formData) => {
    // Remove password confirmation from data sent to backend
    const { confirmPassword, ...registrationData } = formData;

    // Sanitize inputs
    registrationData.username = registrationData.username?.trim();
    registrationData.email = registrationData.email?.trim().toLowerCase();

    setIsLoading(true);
    try {
      await performUserRegistration(registrationData);
      navigate("/");
    } catch (error) {
      // Mostra erro amigável ao utilizador
      toast.error(
        error.message?.replace(/^rpc error: code = Unknown desc = /, "") ||
          "Erro ao criar conta",
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-md mx-auto mt-8">
        <Card>
          <CardHeader>
            <div className="text-center space-y-2">
              <div className="flex justify-center">
                <UserPlus className="w-8 h-8 text-blue-600" />
              </div>
              <h1 className="text-2xl font-bold text-gray-900">
                Create Account
              </h1>
              <p className="text-gray-600">
                Fill in the details to create your account
              </p>
            </div>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <Input
                label="Username"
                type="text"
                autoComplete="username"
                {...register("username", {
                  required: "Username is required",
                  validate: (value) => {
                    const validation = validateName(value);
                    return validation.isValid || validation.error;
                  },
                })}
                error={errors.username?.message}
                placeholder="Enter your username"
              />

              <Input
                label="Email Address"
                type="email"
                autoComplete="email"
                {...register("email", {
                  required: "Email address is required",
                  validate: (value) =>
                    isValidEmail(value) || "Please enter a valid email address",
                })}
                error={errors.email?.message}
                placeholder="Enter your email address"
              />

              <div className="relative">
                <Input
                  label="Password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="new-password"
                  {...register("password", {
                    required: "Password is required",
                    validate: (value) => {
                      const validation = validatePassword(value);
                      if (!validation.isValid) {
                        return validation.errors[0]; // Return first error
                      }
                      return true;
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
              </div>

              <div className="relative">
                <Input
                  label="Confirm Password"
                  type={showConfirmPassword ? "text" : "password"}
                  {...register("confirmPassword", {
                    required: "Password confirmation is required",
                    validate: (value) =>
                      value === passwordValue || "Passwords do not match",
                  })}
                  error={errors.confirmPassword?.message}
                  placeholder="Confirm your password"
                />
                <button
                  type="button"
                  className="password-eye-btn"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
                </button>
              </div>

              <Button
                type="submit"
                className="w-full btn-primary"
                loading={isLoading}
                disabled={isLoading}
              >
                {isLoading ? "Creating account..." : "Create Account"}
              </Button>
            </form>

            <div className="mt-6 text-center space-y-2">
              <p className="text-sm text-gray-600">
                Already have an account?{" "}
                <Link
                  to="/login"
                  className="text-blue-600 hover:text-blue-500 font-medium"
                >
                  Sign in here
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
