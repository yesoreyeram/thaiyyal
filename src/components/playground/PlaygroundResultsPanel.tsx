"use client";
import React, { useState, useRef, useEffect } from "react";

interface PlaygroundResultsPanelProps {
  isLoading: boolean;
  result: unknown;
  error: string | null;
  height: number;
  onHeightChange: (height: number) => void;
}

export function PlaygroundResultsPanel({
  isLoading,
  result,
  error,
  height,
  onHeightChange,
}: PlaygroundResultsPanelProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [activeTab, setActiveTab] = useState<
    "response" | "headers" | "preview" | "timing"
  >("response");
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

      {/* Header with Tabs */}
      <div className="h-10 bg-gray-950 border-b border-gray-800 flex items-center px-4 gap-2">
        <span className="text-sm font-semibold text-gray-200">Response</span>
        {!isLoading && resultData && (
          <>
            <span
              className={`text-xs font-medium px-2 py-0.5 rounded ${
                resultData.status && resultData.status < 400
                  ? "bg-green-900/30 text-green-400"
                  : "bg-red-900/30 text-red-400"
              }`}
            >
              {resultData.status} {resultData.statusText}
            </span>
            <div className="flex-1" />
            <div className="flex gap-1">
              <button
                onClick={() => setActiveTab("response")}
                className={`px-3 py-1 text-xs rounded transition-colors ${
                  activeTab === "response"
                    ? "bg-gray-800 text-white"
                    : "text-gray-400 hover:text-white hover:bg-gray-800/50"
                }`}
              >
                Response
              </button>
              <button
                onClick={() => setActiveTab("headers")}
                className={`px-3 py-1 text-xs rounded transition-colors ${
                  activeTab === "headers"
                    ? "bg-gray-800 text-white"
                    : "text-gray-400 hover:text-white hover:bg-gray-800/50"
                }`}
              >
                Headers
              </button>
              <button
                onClick={() => setActiveTab("preview")}
                className={`px-3 py-1 text-xs rounded transition-colors ${
                  activeTab === "preview"
                    ? "bg-gray-800 text-white"
                    : "text-gray-400 hover:text-white hover:bg-gray-800/50"
                }`}
              >
                Preview
              </button>
              <button
                onClick={() => setActiveTab("timing")}
                className={`px-3 py-1 text-xs rounded transition-colors ${
                  activeTab === "timing"
                    ? "bg-gray-800 text-white"
                    : "text-gray-400 hover:text-white hover:bg-gray-800/50"
                }`}
              >
                Timing
              </button>
            </div>
          </>
        )}
        {isLoading && (
          <>
            <div className="flex items-center gap-2 text-xs text-gray-400">
              <div className="w-3 h-3 border-2 border-blue-500 border-t-transparent rounded-full animate-spin" />
              <span>Executing request...</span>
            </div>
          </>
        )}
      </div>

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
          <div className="space-y-2">
            {(() => {
              try {
                const stringifiedData = JSON.stringify(resultData.data, null, 2);
                return (
                  <>
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-xs text-gray-400">Response Body</span>
                      <span className="text-xs text-gray-500">
                        Size: {JSON.stringify(resultData.data).length} bytes
                      </span>
                    </div>
                    <pre className="text-xs text-gray-300 bg-gray-950 rounded p-3 overflow-x-auto font-mono border border-gray-800">
                      {stringifiedData}
                    </pre>
                  </>
                );
              } catch (error) {
                return (
                  <pre className="text-xs text-red-400 bg-gray-950 rounded p-3 overflow-x-auto">
                    Error formatting response: {error instanceof Error ? error.message : "Unknown error"}
                  </pre>
                );
              }
            })()}
          </div>
        )}

        {!isLoading && resultData && activeTab === "headers" && (
          <div className="space-y-2">
            <div className="text-xs text-gray-400 mb-2">Response Headers</div>
            {resultData.headers && Object.keys(resultData.headers).length > 0 ? (
              <div className="space-y-1">
                {Object.entries(resultData.headers).map(([key, value]) => (
                  <div
                    key={key}
                    className="bg-gray-950 rounded p-2 border border-gray-800 flex gap-3 text-xs"
                  >
                    <span className="font-semibold text-blue-400 min-w-[150px]">
                      {key}:
                    </span>
                    <span className="text-gray-300 break-all">{value}</span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-xs text-gray-500 text-center py-4">
                No headers available
              </div>
            )}
          </div>
        )}

        {!isLoading && resultData && activeTab === "preview" && (
          <div className="space-y-2">
            <div className="text-xs text-gray-400 mb-2">Preview</div>
            <div className="bg-gray-950 rounded p-3 border border-gray-800">
              <div className="text-xs text-gray-300">
                <p>Status: {resultData.status} {resultData.statusText}</p>
                <p className="mt-2">
                  Preview functionality will render formatted response based on content type
                </p>
              </div>
            </div>
          </div>
        )}

        {!isLoading && resultData && activeTab === "timing" && (
          <div className="space-y-2">
            <div className="text-xs text-gray-400 mb-2">Timing</div>
            <div className="bg-gray-950 rounded p-3 border border-gray-800">
              <div className="space-y-2 text-xs">
                <div className="flex justify-between">
                  <span className="text-gray-400">Total Duration:</span>
                  <span className="text-gray-300">~1.5s (mock)</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">DNS Lookup:</span>
                  <span className="text-gray-300">-</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">TCP Connection:</span>
                  <span className="text-gray-300">-</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">TLS Handshake:</span>
                  <span className="text-gray-300">-</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Time to First Byte:</span>
                  <span className="text-gray-300">-</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Content Download:</span>
                  <span className="text-gray-300">-</span>
                </div>
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
