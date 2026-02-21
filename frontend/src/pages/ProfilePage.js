import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useAuthentication } from "../context/AuthContext";
import { authenticationService } from "../services/authService";
import { Layout } from "../components/layout/Layout";
import { Input } from "../components/ui/Input";
import { Card, CardContent, CardHeader } from "../components/ui/Card";
import { User, Save, Eye, EyeOff, Trash2 } from "lucide-react";
import toast from "react-hot-toast";

export function ProfilePage() {
  const { user, performUserLogout } = useAuthentication();
  const [isEditing, setIsEditing] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm({
    defaultValues: {
      username: user?.name || user?.username || "",
      email: user?.email || "",
    },
  });

  const onSubmit = async (data) => {
    setLoading(true);
    try {
      await authenticationService.updateUser(user.id, {
        username: data.username,
        email: data.email,
        password: data.password || undefined, // Only send password if provided
      });

      toast.success("Profile updated successfully!");
      setIsEditing(false);

      // Update current user data in localStorage
      const updatedUser = {
        ...user,
        name: data.username,
        email: data.email,
      };
      authenticationService.saveUserData(updatedUser);
    } catch (error) {
      toast.error(error.message || "Failed to update profile");
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAccount = async () => {
    if (
      window.confirm(
        "Are you sure you want to delete your account? This action cannot be undone.",
      )
    ) {
      setDeleteLoading(true);
      try {
        await authenticationService.deleteUser(user.id);
        toast.success("Account deleted successfully");
        performUserLogout(); // Log out user after deletion
      } catch (error) {
        toast.error(error.message || "Failed to delete account");
      } finally {
        setDeleteLoading(false);
      }
    }
  };

  const handleCancel = () => {
    reset();
    setIsEditing(false);
  };

  /**
   * Format Unix timestamp to English (US) locale
   * @param {number|string} timestamp - Unix timestamp (seconds or milliseconds)
   * @returns {string} Formatted date string
   */
  const formatDate = (timestamp) => {
    if (!timestamp) return "N/A";

    // Convert to number if it's a string
    const numTimestamp = Number(timestamp);

    // If timestamp is in seconds (less than 13 digits), convert to milliseconds
    const milliseconds =
      numTimestamp < 10000000000 ? numTimestamp * 1000 : numTimestamp;

    return new Date(milliseconds).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <Layout>
      <div className="max-w-2xl mx-auto space-y-6">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 flex items-center justify-center">
            <User className="w-8 h-8 mr-2 text-blue-600" />
            My Profile
          </h1>
          <p className="text-gray-600 mt-1">Manage your personal information</p>
        </div>

        {/* Profile Card */}
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold">Personal Information</h2>
              {!isEditing && (
                <button
                  onClick={() => setIsEditing(true)}
                  className="nav-link-green flex items-center"
                >
                  Edit
                </button>
              )}
            </div>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                  label="Username"
                  {...register("username", {
                    required: "Username is required",
                  })}
                  error={errors.username?.message}
                  disabled={!isEditing}
                />

                <Input
                  label="Email"
                  type="email"
                  {...register("email", {
                    required: "Email is required",
                    pattern: {
                      value: /^\S+@\S+$/i,
                      message: "Invalid email",
                    },
                  })}
                  error={errors.email?.message}
                  disabled={!isEditing}
                />
              </div>

              {isEditing && (
                <div className="relative">
                  <Input
                    label="New Password (optional)"
                    type={showPassword ? "text" : "password"}
                    {...register("password", {
                      minLength: {
                        value: 6,
                        message: "Password must be at least 6 characters",
                      },
                    })}
                    error={errors.password?.message}
                    placeholder="Leave blank to keep current password"
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
              )}

              {isEditing && (
                <div className="flex gap-4 pt-4">
                  <button
                    type="submit"
                    disabled={loading}
                    className="nav-link-green flex items-center flex-1 justify-center"
                  >
                    <Save className="w-4 h-4 nav-icon-space" />
                    {loading ? "Saving..." : "Save Changes"}
                  </button>
                  <button
                    type="button"
                    onClick={handleCancel}
                    disabled={loading}
                    className="nav-link-gray flex items-center flex-1 justify-center"
                  >
                    Cancel
                  </button>
                </div>
              )}
            </form>
          </CardContent>
        </Card>

        {/* Account Info */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Account Information</h2>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  User ID
                </label>
                <div className="text-gray-900 font-mono text-sm bg-gray-50 p-2 rounded">
                  {user?.id || "N/A"}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Creation Date
                </label>
                <div className="text-gray-900 text-sm bg-gray-50 p-2 rounded">
                  {formatDate(user?.created)}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Last Updated
                </label>
                <div className="text-gray-900 text-sm bg-gray-50 p-2 rounded">
                  {formatDate(user?.updated)}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Roles
                </label>
                <div className="flex flex-wrap gap-2">
                  {user?.role && user.role.length > 0 ? (
                    user.role.map((role, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                      >
                        {role}
                      </span>
                    ))
                  ) : (
                    <span className="text-gray-500 text-sm">
                      No roles assigned
                    </span>
                  )}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Danger Zone */}
        <Card className="border-red-200">
          <CardHeader>
            <h2 className="text-xl font-semibold text-red-600">Danger Zone</h2>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <p className="text-sm text-gray-600">
                Once you delete your account, there is no going back. Please be
                certain.
              </p>
              <button
                onClick={handleDeleteAccount}
                disabled={deleteLoading}
                className="btn-danger flex items-center"
              >
                <Trash2 className="w-4 h-4 mr-2" />
                {deleteLoading ? "Deleting..." : "Delete Account"}
              </button>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
