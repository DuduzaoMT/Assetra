import React, { createContext, useContext, useReducer } from "react";

const LoadingContext = createContext();

const initialLoadingState = {
  activeRequests: 0,
  isLoading: false,
};

function loadingReducer(state, action) {
  switch (action.type) {
    case "START_REQUEST":
      return {
        ...state,
        activeRequests: state.activeRequests + 1,
        isLoading: true,
      };
    case "END_REQUEST":
      const newActiveRequests = Math.max(0, state.activeRequests - 1);
      return {
        ...state,
        activeRequests: newActiveRequests,
        isLoading: newActiveRequests > 0,
      };
    case "RESET_LOADING":
      return initialLoadingState;
    default:
      return state;
  }
}

export function LoadingProvider({ children }) {
  const [loadingState, dispatchLoading] = useReducer(
    loadingReducer,
    initialLoadingState
  );

  const startLoading = () => {
    dispatchLoading({ type: "START_REQUEST" });
  };

  const stopLoading = () => {
    dispatchLoading({ type: "END_REQUEST" });
  };

  const resetLoading = () => {
    dispatchLoading({ type: "RESET_LOADING" });
  };

  const loadingContextValue = {
    isLoading: loadingState.isLoading,
    activeRequests: loadingState.activeRequests,
    startLoading,
    stopLoading,
    resetLoading,
  };

  return (
    <LoadingContext.Provider value={loadingContextValue}>
      {children}
    </LoadingContext.Provider>
  );
}

export function useLoading() {
  const loadingContext = useContext(LoadingContext);
  if (!loadingContext) {
    throw new Error("useLoading must be used within a LoadingProvider");
  }
  return loadingContext;
}
