"use client";
import React from "react";
import { useRouter } from "next/navigation";

interface AppNavBarProps {
  onNewWorkflow: () => void;
  onOpenWorkflow: () => void;
}

export function AppNavBar({ onNewWorkflow, onOpenWorkflow }: AppNavBarProps) {
  const router = useRouter();

  const handleTitleClick = () => {
    router.push("/");
  };

  const handlePlaygroundClick = () => {
    router.push("/playground");
  };

  return (
    <div className="h-12 bg-gray-950 border-b border-gray-800 flex items-center justify-between px-4">
      <div className="flex items-center gap-3">
        <button
          onClick={handleTitleClick}
          className="text-l font-bold text-white flex items-center gap-2 hover:opacity-80 transition-opacity cursor-pointer"
          title="Go to Home"
        >
          <span className="text-l">⚡</span>
          <span className="bg-linear-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            Thaiyyal
          </span>
        </button>
        <button
          onClick={handlePlaygroundClick}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="HTTP Playground"
          aria-label="HTTP Playground"
        >
          <span>🧪</span>
          <span>Playground</span>
        </button>
      </div>
      <div className="flex items-center gap-3">
        <button
          onClick={onOpenWorkflow}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="Open Workflow"
          aria-label="Open Workflow"
        >
          <span>📂</span>
          <span>Open</span>
        </button>
        <button
          onClick={onNewWorkflow}
          className="px-4 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded transition-all flex items-center gap-1.5 shadow-sm hover:shadow-md"
          title="Create New Workflow"
          aria-label="Create New Workflow"
        >
          <span>⚡</span>
          <span>New</span>
        </button>
      </div>
    </div>
  );
}
