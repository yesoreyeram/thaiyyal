"use client";
import React, { useState, useRef, useEffect } from "react";

interface WorkflowNavBarProps {
  workflowTitle: string;
  onTitleChange: (title: string) => void;
  onSave: () => void;
  onShowJSON: () => void;
  onDelete: () => void;
  onRun: () => void;
  hasUnsavedChanges?: boolean;
}

export function WorkflowNavBar({
  workflowTitle,
  onTitleChange,
  onSave,
  onShowJSON,
  onDelete,
  onRun,
  hasUnsavedChanges = false,
}: WorkflowNavBarProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(workflowTitle);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  useEffect(() => {
    setEditValue(workflowTitle);
  }, [workflowTitle]);

  const handleSubmit = () => {
    if (editValue.trim()) {
      onTitleChange(editValue.trim());
    } else {
      setEditValue(workflowTitle);
    }
    setIsEditing(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSubmit();
    } else if (e.key === "Escape") {
      setEditValue(workflowTitle);
      setIsEditing(false);
    }
  };

  return (
    <div className="h-12 bg-gray-950 border-b border-gray-800 flex items-center justify-between px-4">
      <div className="flex items-center gap-3">
        {isEditing ? (
          <input
            ref={inputRef}
            type="text"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            onBlur={handleSubmit}
            onKeyDown={handleKeyDown}
            className="px-1 py-1 bg-gray-800 border border-gray-600 rounded-lg text-white text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[200px]"
            aria-label="Workflow title"
          />
        ) : (
          <button
            onClick={() => setIsEditing(true)}
            className="px-2 py-1 text-white text-sm font-medium hover:bg-gray-800 rounded-lg transition-colors flex items-center gap-2 group"
            title="Click to edit workflow title"
            aria-label="Edit workflow title"
          >
            <span>{workflowTitle}</span>
            <span className="text-xs opacity-0 group-hover:opacity-100 transition-opacity">
              ✏️
            </span>
          </button>
        )}

        {hasUnsavedChanges && (
          <span className="text-xs text-yellow-500 flex items-center gap-1">
            <span className="w-2 h-2 bg-yellow-500 rounded-full animate-pulse"></span>
            Unsaved changes
          </span>
        )}
      </div>

      <div className="flex items-center gap-3">
        <button
          onClick={onSave}
          className="p-2 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded-lg transition-all"
          title="Save Workflow (Ctrl+S)"
          aria-label="Save Workflow"
        >
          <span>💾</span>
        </button>

        <button
          onClick={onShowJSON}
          className="p-2 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded-lg transition-all"
          title="Show JSON Payload"
          aria-label="Show JSON Payload"
        >
          <span>📋</span>
        </button>

        <button
          onClick={onDelete}
          className="p-2 bg-gray-800 hover:bg-gray-700 text-red-400 hover:text-red-300 rounded-lg transition-all"
          title="Delete Workflow"
          aria-label="Delete Workflow"
        >
          <span>🗑️</span>
        </button>

        <button
          onClick={onRun}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-all flex items-center gap-2 hover:shadow-lg hover:shadow-purple-500/20"
          title="Run Workflow"
          aria-label="Run Workflow"
        >
          <span>▶️</span>
          <span>Run</span>
        </button>
      </div>
    </div>
  );
}
