"use client";
import React from "react";

interface PlaygroundNavBarProps {
  onRun: () => void;
  isRunning: boolean;
}

export function PlaygroundNavBar({ onRun, isRunning }: PlaygroundNavBarProps) {
  return (
    <div className="h-12 bg-gray-950 border-b border-gray-800 flex items-center justify-between px-4">
      <div className="flex items-center gap-3">
        <h1 className="text-white text-sm font-medium">HTTP Playground</h1>
        <span className="text-xs text-gray-500">
          Test and validate HTTP requests
        </span>
      </div>
      <div className="flex items-center gap-2">
        <button
          onClick={onRun}
          disabled={isRunning}
          className={`px-4 py-1.5 ${
            isRunning
              ? "bg-gray-600 cursor-not-allowed"
              : "bg-blue-600 hover:bg-blue-700"
          } text-white text-sm font-medium rounded transition-all flex items-center gap-1.5 shadow-sm hover:shadow-md`}
          title="Run HTTP Request"
          aria-label="Run HTTP Request"
        >
          {isRunning ? (
            <>
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              <span>Running...</span>
            </>
          ) : (
            <>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                className="w-4 h-4"
              >
                <path d="M3 3.732a1.5 1.5 0 0 1 2.305-1.265l6.706 4.267a1.5 1.5 0 0 1 0 2.531l-6.706 4.268A1.5 1.5 0 0 1 3 12.267V3.732Z" />
              </svg>
              <span>Run</span>
            </>
          )}
        </button>
      </div>
    </div>
  );
}
