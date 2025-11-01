"use client";
import React, { useEffect, useRef } from "react";

interface JSONPayloadModalProps {
  isOpen: boolean;
  onClose: () => void;
  payload: object;
}

export function JSONPayloadModal({ isOpen, onClose, payload }: JSONPayloadModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    const handleClickOutside = (e: MouseEvent) => {
      if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      document.addEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.removeEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "";
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const jsonString = JSON.stringify(payload, null, 2);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(jsonString);
    // Could add a toast notification here
  };

  const exportToFile = () => {
    const blob = new Blob([jsonString], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = 'workflow.json';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div
        ref={modalRef}
        className="bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-4xl max-h-[85vh] flex flex-col overflow-hidden"
      >
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-700 flex items-center justify-between bg-gray-800/50">
          <div className="flex items-center gap-3">
            <span className="text-2xl">📋</span>
            <h2 className="text-xl font-semibold text-white">Workflow JSON Payload</h2>
          </div>
          
          <div className="flex items-center gap-2">
            <button
              onClick={exportToFile}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-2"
              title="Export to JSON file"
              aria-label="Export to JSON file"
            >
              <span>💾</span>
              <span>Export</span>
            </button>
            
            <button
              onClick={copyToClipboard}
              className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-2"
              title="Copy to Clipboard"
              aria-label="Copy to Clipboard"
            >
              <span>📋</span>
              <span>Copy</span>
            </button>
            
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm font-medium rounded-lg transition-colors"
              title="Close (ESC)"
              aria-label="Close modal"
            >
              ✕
            </button>
          </div>
        </div>
        
        {/* Content */}
        <div className="flex-1 overflow-auto p-6 custom-scrollbar">
          <pre className="text-sm text-gray-300 font-mono leading-relaxed bg-gray-950 p-4 rounded-lg border border-gray-800">
            {jsonString}
          </pre>
        </div>
        
        {/* Footer */}
        <div className="px-6 py-3 border-t border-gray-700 bg-gray-800/50 flex items-center justify-between text-xs text-gray-500">
          <span>
            {payload && 'nodes' in payload && Array.isArray(payload.nodes) 
              ? `${payload.nodes.length} nodes` 
              : '0 nodes'}
            {' • '}
            {payload && 'edges' in payload && Array.isArray(payload.edges) 
              ? `${payload.edges.length} connections` 
              : '0 connections'}
          </span>
          <span>Press ESC to close</span>
        </div>
      </div>
    </div>
  );
}
