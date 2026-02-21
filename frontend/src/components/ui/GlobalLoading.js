import React from "react";
import { useLoading } from "../../context/LoadingContext";

const GlobalLoading = () => {
  const { isLoading } = useLoading();

  if (!isLoading) return null;

  return (
    <div className="global-loading-overlay">
      <div className="global-loading-content">
        <div className="spinner"></div>
        <p className="loading-text">Loading...</p>
      </div>
    </div>
  );
};

export default GlobalLoading;
