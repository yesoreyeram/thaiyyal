"use client";
import React, { useState, useRef, useEffect } from "react";

interface WorkflowNavBarProps {
  workflowTitle: string;
  onTitleChange: (title: string) => void;
  onSave: () => void;
  onShowJSON: () => void;
  onDelete: () => void;
  onRun: () => void;
  onExport: () => void;
  onImport: (data: { nodes: unknown[]; edges: unknown[] }) => void;
  hasUnsavedChanges?: boolean;
}

export function WorkflowNavBar({
  workflowTitle,
  onTitleChange,
  onSave,
  onShowJSON,
  onDelete,
  onRun,
  onExport,
  onImport,
  hasUnsavedChanges = false,
}: WorkflowNavBarProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(workflowTitle);
  const inputRef = useRef<HTMLInputElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

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

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const json = JSON.parse(event.target?.result as string);
        if (json.nodes && json.edges && Array.isArray(json.nodes) && Array.isArray(json.edges)) {
          onImport(json);
        } else {
          alert('Invalid workflow file format. Expected JSON with "nodes" and "edges" arrays.');
        }
      } catch {
        alert('Failed to parse JSON file. Please ensure it is a valid JSON workflow file.');
      }
    };
    reader.readAsText(file);
    
    // Reset input so the same file can be selected again
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
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
            className="px-2 py-1 bg-gray-800 border border-gray-600 rounded text-white text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[200px]"
            aria-label="Workflow title"
          />
        ) : (
          <div className="flex items-center gap-2">
            <button
              onClick={() => setIsEditing(true)}
              className="px-2 py-1 text-white text-sm font-medium hover:bg-gray-800 rounded transition-colors flex items-center gap-2 group"
              title="Click to edit workflow title"
              aria-label="Edit workflow title"
            >
              <span>{workflowTitle}</span>
              <span className="text-xs opacity-0 group-hover:opacity-100 transition-opacity">
                ✏️
              </span>
            </button>
            <button
              onClick={onDelete}
              className="p-1.5 hover:bg-gray-800 text-red-400 hover:text-red-300 rounded transition-all"
              title="Delete Workflow"
              aria-label="Delete Workflow"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                className="w-4 h-4"
              >
                <path
                  fillRule="evenodd"
                  d="M5 3.25V4H2.75a.75.75 0 0 0 0 1.5h.3l.815 8.15A1.5 1.5 0 0 0 5.357 15h5.285a1.5 1.5 0 0 0 1.493-1.35l.815-8.15h.3a.75.75 0 0 0 0-1.5H11v-.75A2.25 2.25 0 0 0 8.75 1h-1.5A2.25 2.25 0 0 0 5 3.25Zm2.25-.75a.75.75 0 0 0-.75.75V4h3v-.75a.75.75 0 0 0-.75-.75h-1.5ZM6.05 6a.75.75 0 0 1 .787.713l.275 5.5a.75.75 0 0 1-1.498.075l-.275-5.5A.75.75 0 0 1 6.05 6Zm3.9 0a.75.75 0 0 1 .712.787l-.275 5.5a.75.75 0 0 1-1.498-.075l.275-5.5a.75.75 0 0 1 .786-.713Z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          </div>
        )}

        {hasUnsavedChanges && (
          <span className="text-xs text-yellow-500 flex items-center gap-1">
            <span className="w-2 h-2 bg-yellow-500 rounded-full animate-pulse"></span>
            Unsaved changes
          </span>
        )}
      </div>

      <div className="flex items-center gap-2">
        <button
          onClick={onSave}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="Save Workflow (Ctrl+S)"
          aria-label="Save Workflow"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="w-4 h-4"
          >
            <path
              fillRule="evenodd"
              d="M13.78 2.22a.75.75 0 0 1 0 1.06l-7.25 7.25a.75.75 0 0 1-1.06 0L2.22 7.28a.75.75 0 0 1 1.06-1.06L6 8.94l6.72-6.72a.75.75 0 0 1 1.06 0Z"
              clipRule="evenodd"
            />
          </svg>
          <span>Save</span>
        </button>

        <button
          onClick={onExport}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="Export Workflow as JSON"
          aria-label="Export Workflow"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="w-4 h-4"
          >
            <path d="M8.75 2.75a.75.75 0 0 0-1.5 0v5.69L5.03 6.22a.75.75 0 0 0-1.06 1.06l3.5 3.5a.75.75 0 0 0 1.06 0l3.5-3.5a.75.75 0 0 0-1.06-1.06L8.75 8.44V2.75Z" />
            <path d="M3.5 9.75a.75.75 0 0 0-1.5 0v1.5A2.75 2.75 0 0 0 4.75 14h6.5A2.75 2.75 0 0 0 14 11.25v-1.5a.75.75 0 0 0-1.5 0v1.5c0 .69-.56 1.25-1.25 1.25h-6.5c-.69 0-1.25-.56-1.25-1.25v-1.5Z" />
          </svg>
          <span>Export</span>
        </button>

        <button
          onClick={handleImportClick}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="Import Workflow from JSON"
          aria-label="Import Workflow"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="w-4 h-4"
          >
            <path d="M7.25 10.25a.75.75 0 0 0 1.5 0V4.56l2.22 2.22a.75.75 0 1 0 1.06-1.06l-3.5-3.5a.75.75 0 0 0-1.06 0l-3.5 3.5a.75.75 0 0 0 1.06 1.06l2.22-2.22v5.69Z" />
            <path d="M3.5 9.75a.75.75 0 0 0-1.5 0v1.5A2.75 2.75 0 0 0 4.75 14h6.5A2.75 2.75 0 0 0 14 11.25v-1.5a.75.75 0 0 0-1.5 0v1.5c0 .69-.56 1.25-1.25 1.25h-6.5c-.69 0-1.25-.56-1.25-1.25v-1.5Z" />
          </svg>
          <span>Import</span>
        </button>
        
        <input
          ref={fileInputRef}
          type="file"
          accept=".json,application/json"
          onChange={handleFileChange}
          className="hidden"
          aria-label="Import workflow file"
        />

        <button
          onClick={onShowJSON}
          className="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-all text-sm font-medium flex items-center gap-1.5"
          title="Show JSON Payload"
          aria-label="Show JSON Payload"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="w-4 h-4"
          >
            <path
              fillRule="evenodd"
              d="M4.5 2A1.5 1.5 0 0 0 3 3.5v9A1.5 1.5 0 0 0 4.5 14h7a1.5 1.5 0 0 0 1.5-1.5v-9A1.5 1.5 0 0 0 11.5 2h-7Zm.75 3.5a.75.75 0 0 0 0 1.5h5.5a.75.75 0 0 0 0-1.5h-5.5ZM5.25 9a.75.75 0 0 1 .75-.75h4a.75.75 0 0 1 0 1.5H6A.75.75 0 0 1 5.25 9Z"
              clipRule="evenodd"
            />
          </svg>
          <span>JSON</span>
        </button>

        <button
          onClick={onRun}
          className="px-4 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded transition-all flex items-center gap-1.5 shadow-sm hover:shadow-md"
          title="Run Workflow"
          aria-label="Run Workflow"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="w-4 h-4"
          >
            <path d="M3 3.732a1.5 1.5 0 0 1 2.305-1.265l6.706 4.267a1.5 1.5 0 0 1 0 2.531l-6.706 4.268A1.5 1.5 0 0 1 3 12.267V3.732Z" />
          </svg>
          <span>Run</span>
        </button>
      </div>
    </div>
  );
}
