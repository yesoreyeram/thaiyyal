"use client";
import React, { useEffect } from "react";

interface ToastProps {
  message: string;
  type?: "success" | "error" | "info" | "warning";
  isVisible: boolean;
  onClose: () => void;
  duration?: number;
}

export function Toast({
  message,
  type = "info",
  isVisible,
  onClose,
  duration = 3000,
}: ToastProps) {
  useEffect(() => {
    if (isVisible && duration > 0) {
      const timer = setTimeout(() => {
        onClose();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [isVisible, duration, onClose]);

  if (!isVisible) return null;

  const getTypeStyles = () => {
    switch (type) {
      case "success":
        return "bg-green-600 border-green-500 text-white";
      case "error":
        return "bg-red-600 border-red-500 text-white";
      case "warning":
        return "bg-yellow-600 border-yellow-500 text-white";
      case "info":
      default:
        return "bg-blue-600 border-blue-500 text-white";
    }
  };

  const getIcon = () => {
    switch (type) {
      case "success":
        return "✅";
      case "error":
        return "❌";
      case "warning":
        return "⚠️";
      case "info":
      default:
        return "ℹ️";
    }
  };

  return (
    <div className="fixed top-20 right-6 z-60 animate-slide-in">
      <div
        className={`${getTypeStyles()} px-6 py-4 rounded-xl border-2 shadow-2xl flex items-center gap-3 min-w-[300px] max-w-md backdrop-blur-sm`}
      >
        <span className="text-2xl">{getIcon()}</span>
        <span className="flex-1 font-medium">{message}</span>
        <button
          onClick={onClose}
          className="text-white/80 hover:text-white transition-colors ml-2"
          aria-label="Close notification"
        >
          ✕
        </button>
      </div>
    </div>
  );
}
