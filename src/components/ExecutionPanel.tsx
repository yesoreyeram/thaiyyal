"use client";
import React, { useState, useRef, useEffect } from "react";

export interface ExecutionResult {
  success: boolean;
  execution_time?: string;
  results?: {
    execution_id: string;
    node_results: Record<string, unknown>;
    final_output: unknown;
  };
  workflow_id?: string;
  workflow_name?: string;
  error?: string;
}

interface ExecutionPanelProps {
  isOpen: boolean;
  isLoading: boolean;
  result: ExecutionResult | null;
  error: string | null;
  onCancel: () => void;
  onClose: () => void;
  height: number;
  onHeightChange: (height: number) => void;
}

export function ExecutionPanel({
  isOpen,
  isLoading,
  result,
  error,
  onCancel,
  onClose,
  height,
  onHeightChange,
}: ExecutionPanelProps) {
  const [isDragging, setIsDragging] = useState(false);
  const dragStartY = useRef(0);
  const dragStartHeight = useRef(0);

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = dragStartY.current - e.clientY;
      const newHeight = Math.max(100, Math.min(600, dragStartHeight.current + deltaY));
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
            Execution Results
          </span>
          {isLoading && (
            <div className="flex items-center gap-2 text-xs text-gray-400">
              <div className="w-3 h-3 border-2 border-blue-500 border-t-transparent rounded-full animate-spin" />
              <span>Running workflow...</span>
            </div>
          )}
          {!isLoading && result && (
            <span className="text-xs text-gray-400">
              {result.success ? (
                <span className="text-green-400 flex items-center gap-1">
                  <span className="w-1.5 h-1.5 bg-green-400 rounded-full" />
                  Completed in {result.execution_time}
                </span>
              ) : (
                <span className="text-red-400 flex items-center gap-1">
                  <span className="w-1.5 h-1.5 bg-red-400 rounded-full" />
                  Failed
                </span>
              )}
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          {isLoading && (
            <button
              onClick={onCancel}
              className="px-3 py-1 text-xs bg-red-600 hover:bg-red-700 text-white rounded transition-colors"
              title="Cancel Execution"
            >
              Cancel
            </button>
          )}
          <button
            onClick={onClose}
            className="px-2 py-1 text-gray-400 hover:text-white transition-colors"
            title="Close Panel"
          >
            ✕
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-4">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full text-gray-400">
            <div className="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4" />
            <p className="text-sm">Executing workflow...</p>
            <p className="text-xs text-gray-500 mt-2">
              This may take a few moments
            </p>
          </div>
        )}

        {!isLoading && error && (
          <div className="bg-red-900/20 border border-red-800 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <div className="w-5 h-5 bg-red-500 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                <span className="text-white text-xs font-bold">!</span>
              </div>
              <div className="flex-1">
                <h3 className="text-sm font-semibold text-red-400 mb-2">
                  Execution Error
                </h3>
                <pre className="text-xs text-gray-300 bg-gray-950 rounded p-3 overflow-x-auto">
                  {error}
                </pre>
              </div>
            </div>
          </div>
        )}

        {!isLoading && result && result.success && (
          <div className="space-y-4">
            {/* Summary */}
            <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
              <h3 className="text-sm font-semibold text-gray-200 mb-3">
                Execution Summary
              </h3>
              <div className="grid grid-cols-2 gap-3 text-xs">
                <div>
                  <span className="text-gray-400">Execution ID:</span>
                  <span className="ml-2 text-gray-200 font-mono">
                    {result.results?.execution_id}
                  </span>
                </div>
                <div>
                  <span className="text-gray-400">Duration:</span>
                  <span className="ml-2 text-green-400">{result.execution_time}</span>
                </div>
                {result.workflow_id && (
                  <div>
                    <span className="text-gray-400">Workflow ID:</span>
                    <span className="ml-2 text-gray-200 font-mono">
                      {result.workflow_id}
                    </span>
                  </div>
                )}
                {result.workflow_name && (
                  <div>
                    <span className="text-gray-400">Workflow:</span>
                    <span className="ml-2 text-gray-200">{result.workflow_name}</span>
                  </div>
                )}
              </div>
            </div>

            {/* Final Output */}
            <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
              <h3 className="text-sm font-semibold text-gray-200 mb-3">
                Final Output
              </h3>
              <pre className="text-xs text-gray-300 bg-gray-950 rounded p-3 overflow-x-auto max-h-32">
                {JSON.stringify(result.results?.final_output, null, 2)}
              </pre>
            </div>

            {/* Node Results */}
            {result.results?.node_results && (
              <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
                <h3 className="text-sm font-semibold text-gray-200 mb-3">
                  Node Results
                </h3>
                <div className="space-y-2">
                  {Object.entries(result.results.node_results).map(
                    ([nodeId, value]) => (
                      <div
                        key={nodeId}
                        className="bg-gray-950 rounded p-3 border border-gray-800"
                      >
                        <div className="text-xs font-semibold text-blue-400 mb-2">
                          Node: {nodeId}
                        </div>
                        <pre className="text-xs text-gray-300 overflow-x-auto">
                          {JSON.stringify(value, null, 2)}
                        </pre>
                      </div>
                    )
                  )}
                </div>
              </div>
            )}
          </div>
        )}

        {!isLoading && !result && !error && (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <div className="text-4xl mb-4">▶️</div>
            <p className="text-sm">Click "Run" to execute the workflow</p>
            <p className="text-xs text-gray-600 mt-2">
              Results will appear here
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
