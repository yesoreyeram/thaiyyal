"use client";
import React from "react";

interface WorkflowStatusBarProps {
  nodeCount: number;
  edgeCount: number;
  selectedCount?: number;
}

export function WorkflowStatusBar({ 
  nodeCount, 
  edgeCount, 
  selectedCount = 0 
}: WorkflowStatusBarProps) {
  return (
    <div className="h-7 bg-gray-900 border-t border-gray-800 flex items-center justify-between px-4 text-xs text-gray-400">
      <div className="flex items-center gap-4">
        <span className="flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>
          <span>{nodeCount} nodes</span>
        </span>
        <span className="flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
          <span>{edgeCount} connections</span>
        </span>
        {selectedCount > 0 && (
          <span className="flex items-center gap-1.5">
            <span className="w-1.5 h-1.5 bg-purple-500 rounded-full"></span>
            <span>{selectedCount} selected</span>
          </span>
        )}
      </div>
      <div className="flex items-center gap-2 text-gray-500">
        <span>Ready</span>
      </div>
    </div>
  );
}
