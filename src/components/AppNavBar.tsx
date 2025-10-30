"use client";
import React from "react";

interface AppNavBarProps {
  onNewWorkflow: () => void;
  onOpenWorkflow: () => void;
}

export function AppNavBar({ onNewWorkflow, onOpenWorkflow }: AppNavBarProps) {
  return (
    <div className="h-14 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-6">
      <div className="flex items-center gap-3">
        <h1 className="text-xl font-bold text-white flex items-center gap-2">
          <span className="text-2xl">⚡</span>
          <span className="bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            Thaiyyal
          </span>
        </h1>
        <span className="text-xs text-gray-500 font-medium px-2 py-1 bg-gray-800 rounded-md">
          Workflow Builder
        </span>
      </div>
      
      <div className="flex items-center gap-2">
        <button
          onClick={onNewWorkflow}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-all flex items-center gap-2 hover:shadow-lg hover:shadow-blue-500/20"
          title="Create New Workflow"
          aria-label="Create New Workflow"
        >
          <span>📄</span>
          <span>New Workflow</span>
        </button>
        
        <button
          onClick={onOpenWorkflow}
          className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm font-medium rounded-lg transition-all flex items-center gap-2"
          title="Open Workflow"
          aria-label="Open Workflow"
        >
          <span>📂</span>
          <span>Open Workflow</span>
        </button>
      </div>
    </div>
  );
}
