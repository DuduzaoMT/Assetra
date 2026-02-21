import apiClient from "./api";

// In-memory storage for access token (not accessible by JavaScript when page reloads)
let accessToken = null;

export const authenticationService = {
  // Get the current access token from memory
  getAccessToken() {
    return accessToken;
  },

  // Set the access token in memory
  setAccessToken(token) {
    accessToken = token;
  },

  // Clear the access token from memory
  clearAccessToken() {
    accessToken = null;
  },

  // Refresh the access token using the refresh token from cookies
  async refreshAccessToken() {
    try {
      const response = await apiClient.post("/refresh-token");
      const newAccessToken = response.data.access_token;
      this.setAccessToken(newAccessToken);
      return newAccessToken;
    } catch (error) {
      // If refresh fails, clear everything
      this.clearAccessToken();
      localStorage.removeItem("currentUser");
      // Re-throw without wrapping - let the caller decide what to do
      throw error;
    }
  },

  // Sign in user with email and password
  async signIn(loginCredentials) {
    try {
      const response = await apiClient.post("/signin", {
        email: loginCredentials.email,
        password: loginCredentials.password,
      });

      const userData = response.data.user;
      const token = response.data.access_token;

      // Store access token in memory
      this.setAccessToken(token);

      // Save user data to localStorage (NOT the token)
      this.saveUserData(userData);

      return { user: userData };
    } catch (error) {
      throw new Error(error.response?.data?.error || "Failed to sign in");
    }
  },

  // Sign up new user
  async signUp(userRegistrationData) {
    try {
      const response = await apiClient.post("/signup", {
        name: userRegistrationData.username,
        email: userRegistrationData.email,
        password: userRegistrationData.password,
      });

      const userData = response.data.user;
      const token = response.data.access_token;

      // Store access token in memory
      this.setAccessToken(token);

      // Save user data to localStorage (NOT the token)
      this.saveUserData(userData);

      return { user: userData };
    } catch (error) {
      throw new Error(
        error.response?.data?.error || "Failed to create account",
      );
    }
  },

  // Update user profile
  async updateUser(userId, updateData) {
    try {
      const response = await apiClient.put(`/users/${userId}`, {
        id: userId,
        name: updateData.username,
        email: updateData.email,
        ...(updateData.password && { password: updateData.password }),
      });
      return response.data;
    } catch (error) {
      throw new Error(
        error.response?.data?.error || "Failed to update profile",
      );
    }
  },

  // Delete user account
  async deleteUser(userId) {
    try {
      const response = await apiClient.delete(`/users/${userId}`);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || "Failed to delete user");
    }
  },

  // Get specific user by ID
  async getUser(userId) {
    try {
      const response = await apiClient.get(`/users/${userId}`);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || "Failed to get user");
    }
  },

  // Sign out user - call backend to clear cookies
  async signOut() {
    try {
      await apiClient.post("/logout");
      this.clearAccessToken();
      localStorage.removeItem("currentUser");
    } catch (error) {
      // Even if logout fails, clear local data
      this.clearAccessToken();
      localStorage.removeItem("currentUser");
    }
  },

  // Check if user is authenticated by checking if user data exists
  isUserAuthenticated() {
    const currentUserData = localStorage.getItem("currentUser");
    return currentUserData !== null;
  },

  // Get current authenticated user data
  getCurrentAuthenticatedUser() {
    const currentUserData = localStorage.getItem("currentUser");
    return currentUserData ? JSON.parse(currentUserData) : null;
  },

  // Alias for compatibility
  getCurrentUser() {
    return this.getCurrentAuthenticatedUser();
  },

  // Save user data to localStorage (tokens are in httpOnly cookies)
  saveUserData(userData) {
    localStorage.setItem("currentUser", JSON.stringify(userData));
  },
};

// Default export for compatibility
export default authenticationService;
