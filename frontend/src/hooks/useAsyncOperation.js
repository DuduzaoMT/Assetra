import { useState } from "react";
import { useLoading } from "../context/LoadingContext";

/**
 * Custom hook for controlling loading in specific requests
 * Can be used when you need manual control beyond global loading
 */
export function useAsyncOperation() {
  const [isLocalLoading, setIsLocalLoading] = useState(false);
  const [error, setError] = useState(null);
  const { startLoading, stopLoading } = useLoading();

  const executeAsync = async (asyncFunction, options = {}) => {
    const {
      useGlobalLoading = true,
      useLocalLoading = false,
      onSuccess,
      onError,
    } = options;

    try {
      if (useLocalLoading) setIsLocalLoading(true);
      if (useGlobalLoading) startLoading();

      setError(null);

      const result = await asyncFunction();

      if (onSuccess) onSuccess(result);

      return { success: true, data: result };
    } catch (err) {
      setError(err.message);

      if (onError) onError(err);

      return { success: false, error: err.message };
    } finally {
      if (useLocalLoading) setIsLocalLoading(false);
      if (useGlobalLoading) stopLoading();
    }
  };

  return {
    executeAsync,
    isLocalLoading,
    error,
    clearError: () => setError(null),
  };
}
