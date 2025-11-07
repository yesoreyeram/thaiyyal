"use client";
import React, { useState, useRef, useEffect } from "react";

interface PlaygroundResultsPanelProps {
  isOpen: boolean;
  isLoading: boolean;
  result: unknown;
  error: string | null;
  onClose: () => void;
  height: number;
  onHeightChange: (height: number) => void;
}

export function PlaygroundResultsPanel({
  isOpen,
  isLoading,
  result,
  error,
  onClose,
  height,
  onHeightChange,
}: PlaygroundResultsPanelProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [activeTab, setActiveTab] = useState<"response" | "headers">(
    "response"
  );
  const dragStartY = useRef(0);
  const dragStartHeight = useRef(0);

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = dragStartY.current - e.clientY;
      const newHeight = Math.max(
        100,
        Math.min(600, dragStartHeight.current + deltaY)
      );
      onHeightChange(newHeight);
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isDragging, onHeightChange]);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    dragStartY.current = e.clientY;
    dragStartHeight.current = height;
  };

  if (!isOpen) return null;

  const resultData = result as {
    status?: number;
    statusText?: string;
    headers?: Record<string, string>;
    data?: unknown;
  } | null;

  return (
    <div
      className="bg-gray-900 border-t border-gray-800 flex flex-col"
      style={{ height: `${height}px` }}
    >
      {/* Resize Handle */}
      <div
        className={`h-1 bg-gray-800 hover:bg-blue-500 cursor-ns-resize transition-colors ${
          isDragging ? "bg-blue-500" : ""
        }`}
        onMouseDown={handleMouseDown}
      />

      {/* Header */}
      <div className="h-10 bg-gray-950 border-b border-gray-800 flex items-center justify-between px-4">
        <div className="flex items-center gap-3">
          <span className="text-sm font-semibold text-gray-200">
            Response
          </span>
          {isLoading && (
            <div className="flex items-center gap-2 text-xs text-gray-400">
              <div className="w-3 h-3 border-2 border-blue-500 border-t-transparent rounded-full animate-spin" />
              <span>Executing request...</span>
            </div>
          )}
          {!isLoading && resultData && (
            <span
              className={`text-xs font-medium px-2 py-0.5 rounded ${
                resultData.status && resultData.status < 400
                  ? "bg-green-900/30 text-green-400"
                  : "bg-red-900/30 text-red-400"
              }`}
            >
              {resultData.status} {resultData.statusText}
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={onClose}
            className="px-2 py-1 text-gray-400 hover:text-white transition-colors"
            title="Close Panel"
          >
            ✕
          </button>
        </div>
      </div>

      {/* Tabs */}
      {!isLoading && resultData && (
        <div className="h-10 bg-gray-950 border-b border-gray-800 flex items-center px-4 gap-1">
          <button
            onClick={() => setActiveTab("response")}
            className={`px-3 py-1.5 text-sm rounded transition-colors ${
              activeTab === "response"
                ? "bg-gray-800 text-white"
                : "text-gray-400 hover:text-white hover:bg-gray-800/50"
            }`}
          >
            Response
          </button>
          <button
            onClick={() => setActiveTab("headers")}
            className={`px-3 py-1.5 text-sm rounded transition-colors ${
              activeTab === "headers"
                ? "bg-gray-800 text-white"
                : "text-gray-400 hover:text-white hover:bg-gray-800/50"
            }`}
          >
            Headers
          </button>
        </div>
      )}

      {/* Content */}
      <div className="flex-1 overflow-auto p-4">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full text-gray-400">
            <div className="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4" />
            <p className="text-sm">Executing HTTP request...</p>
            <p className="text-xs text-gray-500 mt-2">
              This may take a few moments
            </p>
          </div>
        )}

        {!isLoading && error && (
          <div className="bg-red-900/20 border border-red-800 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <div className="w-5 h-5 bg-red-500 rounded-full flex items-center justify-center shrink-0 mt-0.5">
                <span className="text-white text-xs font-bold">!</span>
              </div>
              <div className="flex-1">
                <h3 className="text-sm font-semibold text-red-400 mb-2">
                  Request Error
                </h3>
                <pre className="text-xs text-gray-300 bg-gray-950 rounded p-3 overflow-x-auto">
                  {error}
                </pre>
              </div>
            </div>
          </div>
        )}

        {!isLoading && resultData && activeTab === "response" && (
          <div className="space-y-4">
            {/* Response Body */}
            <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
              <h3 className="text-sm font-semibold text-gray-200 mb-3">
                Response Body
              </h3>
              <pre className="text-xs text-gray-300 bg-gray-950 rounded p-3 overflow-x-auto font-mono">
                {(() => {
                  try {
                    return JSON.stringify(resultData.data, null, 2);
                  } catch (error) {
                    return `Error formatting response: ${error instanceof Error ? error.message : "Unknown error"}`;
                  }
                })()}
              </pre>
            </div>
          </div>
        )}

        {!isLoading && resultData && activeTab === "headers" && (
          <div className="space-y-4">
            {/* Response Headers */}
            <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
              <h3 className="text-sm font-semibold text-gray-200 mb-3">
                Response Headers
              </h3>
              <div className="space-y-2">
                {resultData.headers ? (
                  Object.entries(resultData.headers).map(([key, value]) => (
                    <div
                      key={key}
                      className="bg-gray-950 rounded p-3 border border-gray-800 flex justify-between items-start gap-4"
                    >
                      <span className="text-xs font-semibold text-blue-400 break-all">
                        {key}:
                      </span>
                      <span className="text-xs text-gray-300 break-all text-right">
                        {value}
                      </span>
                    </div>
                  ))
                ) : (
                  <div className="text-xs text-gray-500">No headers available</div>
                )}
              </div>
            </div>
          </div>
        )}

        {!isLoading && !resultData && !error && (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <div className="text-4xl mb-4">▶️</div>
            <p className="text-sm">Click &ldquo;Run&rdquo; to execute the HTTP request</p>
            <p className="text-xs text-gray-600 mt-2">
              Results will appear here
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
