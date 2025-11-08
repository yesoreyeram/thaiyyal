"use client";
import React, { useState, useEffect } from "react";
import { Workflow } from "../types/workflow";

interface OpenWorkflowModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (workflow: Workflow) => void;
}

export function OpenWorkflowModal({
  isOpen,
  onClose,
  onSelect,
}: OpenWorkflowModalProps) {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);

  const loadWorkflows = React.useCallback(() => {
    const saved = localStorage.getItem("thaiyyal_workflows");
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setWorkflows(parsed);
      } catch (e) {
        console.error("Failed to load workflows", e);
        setWorkflows([]);
      }
    } else {
      setWorkflows([]);
    }
  }, []);

  useEffect(() => {
    if (isOpen) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      loadWorkflows();
    }
  }, [isOpen, loadWorkflows]);

  const handleDelete = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    const updated = workflows.filter((w) => w.id !== id);
    setWorkflows(updated);
    localStorage.setItem("thaiyyal_workflows", JSON.stringify(updated));
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return (
      date.toLocaleDateString() +
      " " +
      date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
    );
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-black border border-gray-300 dark:border-gray-700 rounded-2xl shadow-2xl w-full max-w-3xl max-h-[80vh] flex flex-col overflow-hidden">
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-300 dark:border-gray-700 flex items-center justify-between bg-gray-50 dark:bg-gray-950">
          <div className="flex items-center gap-3">
            <span className="text-2xl">📂</span>
            <h2 className="text-xl font-semibold text-black dark:text-white">Open Workflow</h2>
          </div>

          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-200 dark:bg-gray-800 hover:bg-gray-300 dark:hover:bg-gray-700 text-black dark:text-white text-sm font-medium rounded-lg transition-colors"
            aria-label="Close modal"
          >
            ✕
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-auto p-6 custom-scrollbar">
          {workflows.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-6xl mb-4 opacity-50">📭</div>
              <p className="text-gray-600 dark:text-gray-400 text-lg mb-2">No saved workflows</p>
              <p className="text-gray-500 text-sm">
                Create and save a workflow to see it here
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {workflows.map((workflow) => (
                <div
                  key={workflow.id}
                  onClick={() => {
                    onSelect(workflow);
                    onClose();
                  }}
                  className="group p-4 bg-gray-50 dark:bg-gray-900 hover:bg-gray-100 dark:hover:bg-gray-800 border border-gray-300 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-600 rounded-lg cursor-pointer transition-all"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1 min-w-0">
                      <h3 className="text-black dark:text-white font-medium mb-1 truncate transition-colors">
                        {workflow.title}
                      </h3>
                      <div className="flex items-center gap-4 text-xs text-gray-600 dark:text-gray-400">
                        <span className="flex items-center gap-1">
                          <span className="w-1.5 h-1.5 bg-black dark:bg-white rounded-full"></span>
                          {workflow.data.nodes.length} nodes
                        </span>
                        <span className="flex items-center gap-1">
                          <span className="w-1.5 h-1.5 bg-black dark:bg-white rounded-full"></span>
                          {workflow.data.edges.length} connections
                        </span>
                        <span>Updated: {formatDate(workflow.updatedAt)}</span>
                      </div>
                    </div>

                    <button
                      onClick={(e) => handleDelete(workflow.id, e)}
                      className="ml-4 px-3 py-1.5 bg-gray-200 dark:bg-gray-800 hover:bg-gray-300 dark:hover:bg-gray-700 text-black dark:text-white text-sm rounded-lg transition-colors opacity-0 group-hover:opacity-100"
                      title="Delete workflow"
                      aria-label="Delete workflow"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-3 border-t border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-950 flex items-center justify-between text-xs text-gray-600 dark:text-gray-400">
          <span>
            {workflows.length} workflow{workflows.length !== 1 ? "s" : ""} saved
          </span>
          <span>Click to open • Press ESC to close</span>
        </div>
      </div>
    </div>
  );
}
