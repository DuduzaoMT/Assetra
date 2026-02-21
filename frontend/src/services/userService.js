import apiClient from "./api";

export const userManagementService = {
  // Fetch all users from the server
  async getAllUsers() {
    try {
      const response = await apiClient.get("/users");
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to fetch users");
    }
  },

  // Fetch specific user by ID
  async getUserById(userId) {
    try {
      const response = await apiClient.get(`/users/${userId}`);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to fetch user");
    }
  },

  // Update user information
  async updateUserData(userId, userUpdateData) {
    try {
      const response = await apiClient.put(`/users/${userId}`, userUpdateData);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to update user");
    }
  },

  // Delete user by ID
  async deleteUserById(userId) {
    try {
      const response = await apiClient.delete(`/users/${userId}`);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to delete user");
    }
  },
};
