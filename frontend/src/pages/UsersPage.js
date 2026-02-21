import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Layout } from "../components/layout/Layout";
import { useUserManagement } from "../hooks/useUsers";
import { authenticationService } from "../services/authService";
import { Card, CardContent, CardHeader } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import { Users, Edit, Trash2, Search, X, Save } from "lucide-react";
import { Input } from "../components/ui/Input";
import toast from "react-hot-toast";

export function UsersPage() {
  const { usersList, usersError, removeUserFromSystem, refreshUsersList } =
    useUserManagement();
  const [searchQuery, setSearchQuery] = useState("");
  const [userBeingDeleted, setUserBeingDeleted] = useState(null);
  const [editingUser, setEditingUser] = useState(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [updateLoading, setUpdateLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm();

  const filteredUsersList = usersList.filter(
    (user) =>
      user.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.email?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleUserDeletion = async (userId) => {
    if (window.confirm("Are you sure you want to delete this user?")) {
      setUserBeingDeleted(userId);
      await removeUserFromSystem(userId);
      setUserBeingDeleted(null);
    }
  };

  const handleEditUser = (user) => {
    setEditingUser(user);
    setValue("username", user.name);
    setValue("email", user.email);
    setValue("password", ""); // Clear password field
    setIsEditModalOpen(true);
  };

  const handleCloseEditModal = () => {
    setIsEditModalOpen(false);
    setEditingUser(null);
    reset();
  };

  const onUpdateSubmit = async (data) => {
    if (!editingUser) return;

    setUpdateLoading(true);
    try {
      await authenticationService.updateUser(editingUser.id, {
        username: data.username,
        email: data.email,
        password: data.password || undefined, // Only send password if provided
      });

      toast.success("User updated successfully!");
      handleCloseEditModal();
      await refreshUsersList(); // Refresh the users list
    } catch (error) {
      toast.error(error.message || "Failed to update user");
    } finally {
      setUpdateLoading(false);
    }
  };

  const formatTimestamp = (timestamp) => {
    if (!timestamp) return "N/A";

    // Convert to number if it's a string
    const numTimestamp = Number(timestamp);

    // If timestamp is in seconds (less than 13 digits), convert to milliseconds
    const milliseconds =
      numTimestamp < 10000000000 ? numTimestamp * 1000 : numTimestamp;

    return new Date(milliseconds).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (usersError) {
    return (
      <Layout>
        <div className="text-center py-12">
          <div className="text-red-600 text-lg">{usersError}</div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-8">
        {/* Page Header */}
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center">
              <Users className="w-8 h-8 mr-2 text-blue-600" />
              Users
            </h1>
            <p className="text-gray-600 mt-1">Manage all platform users</p>
          </div>
        </div>

        {/* Search Section */}
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <Search className="text-gray-400 w-4 h-4" />
              <Input
                type="text"
                placeholder="Search users by name or email..."
                variant="outline"
                size="sm"
                className="btn-generic p-1 rounded"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </CardContent>
        </Card>

        {/* Users List */}
        <div className="grid gap-4">
          {filteredUsersList.length === 0 ? (
            <Card>
              <CardContent className="text-center py-12">
                <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-600">
                  {searchQuery ? "No users found" : "No users registered"}
                </p>
              </CardContent>
            </Card>
          ) : (
            filteredUsersList.map((userItem) => (
              <Card
                key={userItem.id}
                className="hover:shadow-md transition-shadow"
              >
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900">
                          {userItem.name}
                        </h3>
                        <p className="text-gray-600">{userItem.email}</p>
                        <div className="flex items-center space-x-4 mt-1 text-sm text-gray-500">
                          <span>ID: {userItem.id}</span>
                          {userItem.created && (
                            <span>
                              Created: {formatTimestamp(userItem.created)}
                            </span>
                          )}
                        </div>
                        {userItem.role && userItem.role.length > 0 && (
                          <div className="flex items-center space-x-2 mt-2">
                            {userItem.role.map((userRole, roleIndex) => (
                              <span
                                key={roleIndex}
                                className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                              >
                                {userRole}
                              </span>
                            ))}
                          </div>
                        )}
                      </div>
                    </div>

                    <div className="flex items-center gap-1">
                      <Button
                        variant="outline"
                        size="sm"
                        className="btn-generic px-2 py-1 text-xs font-medium rounded flex items-center h-7 min-w-0"
                        onClick={() => handleEditUser(userItem)}
                      >
                        <Edit className="w-3.5 h-3.5 mr-1" />
                        Edit
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        className="btn-danger px-2 py-1 text-xs font-medium rounded flex items-center h-7 min-w-0"
                        onClick={() => handleUserDeletion(userItem.id)}
                        loading={userBeingDeleted === userItem.id}
                        disabled={userBeingDeleted === userItem.id}
                      >
                        <Trash2 className="w-3.5 h-3.5 mr-1" />
                        Delete
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        {/* Statistics Section */}
        <Card>
          <CardHeader>
            <h3 className="text-lg font-semibold">Statistics</h3>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {usersList.length}
                </div>
                <div className="text-sm text-gray-600">Total Users</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">
                  {usersList.filter((u) => u.role?.includes("admin")).length}
                </div>
                <div className="text-sm text-gray-600">Administrators</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-purple-600">
                  {searchQuery ? filteredUsersList.length : usersList.length}
                </div>
                <div className="text-sm text-gray-600">
                  {searchQuery ? "Search Results" : "Active Users"}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Edit User Modal */}
      {isEditModalOpen && (
        <div>
          {/* Backdrop with blur effect */}
          <div className="modal-backdrop" onClick={handleCloseEditModal}></div>

          {/* Modal container */}
          <div className="modal-container">
            <Card className="modal-content">
              <CardHeader>
                <div className="flex justify-between items-center">
                  <h3 className="text-lg font-semibold">Edit User</h3>
                  <button
                    onClick={handleCloseEditModal}
                    className="modal-close-button"
                  >
                    <X className="w-5 h-5" />
                  </button>
                </div>
              </CardHeader>
              <CardContent>
                <form
                  onSubmit={handleSubmit(onUpdateSubmit)}
                  className="space-y-4"
                >
                  <Input
                    label="Username"
                    {...register("username", {
                      required: "Username is required",
                    })}
                    error={errors.username?.message}
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
                  />

                  <Input
                    label="New Password (optional)"
                    type="password"
                    {...register("password", {
                      minLength: {
                        value: 6,
                        message: "Password must be at least 6 characters",
                      },
                    })}
                    error={errors.password?.message}
                    placeholder="Leave blank to keep current password"
                  />

                  <div className="flex gap-4 pt-4">
                    <button
                      type="submit"
                      disabled={updateLoading}
                      className="nav-link-green flex items-center flex-1 justify-center"
                    >
                      <Save className="w-4 h-4 mr-2" />
                      {updateLoading ? "Saving..." : "Save Changes"}
                    </button>
                    <button
                      type="button"
                      onClick={handleCloseEditModal}
                      disabled={updateLoading}
                      className="nav-link-gray flex items-center flex-1 justify-center"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>
        </div>
      )}
    </Layout>
  );
}
