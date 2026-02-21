import axios from "axios";
import { toast } from "react-hot-toast";

const API_BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:9000";

let loadingContext = null;
let authServiceInstance = null;

// Flag to prevent multiple simultaneous refresh attempts
let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

// Create axios instance with default configuration
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000,
  withCredentials: true, // Important: enables cookies to be sent
});

export const setupLoadingInterceptors = (context) => {
  loadingContext = context;
};

// Setup auth service for token refresh
export const setupAuthService = (authService) => {
  authServiceInstance = authService;
};

// Request interceptor - add access token to headers
apiClient.interceptors.request.use(
  (requestConfig) => {
    if (loadingContext) {
      loadingContext.startLoading();
    }

    // Get access token from authService and add to headers
    if (authServiceInstance) {
      const token = authServiceInstance.getAccessToken();
      if (token) {
        requestConfig.headers.Authorization = `Bearer ${token}`;
      }
    }

    return requestConfig;
  },
  (requestError) => {
    if (loadingContext) {
      loadingContext.stopLoading();
    }
    return Promise.reject(requestError);
  },
);

// Response interceptor with automatic token refresh
apiClient.interceptors.response.use(
  (response) => {
    if (loadingContext) {
      loadingContext.stopLoading();
    }
    return response;
  },
  async (error) => {
    if (loadingContext) {
      loadingContext.stopLoading();
    }

    const originalRequest = error.config;

    // Handle 401 errors with token refresh
    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !originalRequest.url?.includes("/refresh-token")
    ) {
      if (isRefreshing) {
        // If already refreshing, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return apiClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      if (authServiceInstance) {
        try {
          // Try to refresh the access token
          const newToken = await authServiceInstance.refreshAccessToken();

          // Update the failed request with new token
          originalRequest.headers.Authorization = `Bearer ${newToken}`;

          // Process queued requests
          processQueue(null, newToken);

          isRefreshing = false;

          // Retry the original request
          return apiClient(originalRequest);
        } catch (refreshError) {
          // Refresh failed - clear auth and redirect to login
          processQueue(refreshError, null);
          isRefreshing = false;

          // Only show toast and redirect if we're not already on login page
          if (authServiceInstance) {
            authServiceInstance.clearAccessToken();
          }
          localStorage.removeItem("currentUser");

          if (window.location.pathname !== "/login") {
            toast.error("Session expired. Please login again.");
            window.location.href = "/login";
          }

          return Promise.reject(refreshError);
        }
      }
    }

    // If this is a 401 from /refresh-token itself, clear everything and redirect
    if (
      error.response?.status === 401 &&
      originalRequest.url?.includes("/refresh-token")
    ) {
      if (authServiceInstance) {
        authServiceInstance.clearAccessToken();
      }
      localStorage.removeItem("currentUser");
      // Don't show error or redirect if already on login page
      if (
        window.location.pathname !== "/login" &&
        window.location.pathname !== "/register"
      ) {
        window.location.href = "/login";
      }
      return Promise.reject(error);
    } else if (error.response?.status === 429) {
      toast.error("Too many requests. Please slow down.");
    } else if (error.response?.status >= 500) {
      toast.error("Server error. Please try again later.");
    }

    return Promise.reject(error);
  },
);

export default apiClient;
