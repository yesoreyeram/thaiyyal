"use client";
import React, { useState, useEffect, useRef } from "react";

interface NodeConfig {
  type: string;
  label: string;
  color: string;
  defaultData: Record<string, unknown>;
}

interface NodeCategory {
  name: string;
  nodes: NodeConfig[];
}

interface NodePaletteProps {
  isOpen: boolean;
  onClose: () => void;
  categories: NodeCategory[];
  onAddNode: (type: string, defaultData: Record<string, unknown>) => void;
}

export function NodePalette({ isOpen, onClose, categories, onAddNode }: NodePaletteProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(
    new Set(categories.map(c => c.name))
  );
  const searchInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      // Focus search input when palette opens
      setTimeout(() => searchInputRef.current?.focus(), 100);
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose]);

  const toggleCategory = (categoryName: string) => {
    setExpandedCategories(prev => {
      const newSet = new Set(prev);
      if (newSet.has(categoryName)) {
        newSet.delete(categoryName);
      } else {
        newSet.add(categoryName);
      }
      return newSet;
    });
  };

  const filteredCategories = categories.map(category => ({
    ...category,
    nodes: category.nodes.filter(node =>
      node.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
      node.type.toLowerCase().includes(searchQuery.toLowerCase())
    ),
  })).filter(category => category.nodes.length > 0);

  const handleAddNode = (type: string, defaultData: Record<string, unknown>) => {
    onAddNode(type, defaultData);
    setSearchQuery(""); // Clear search after adding
  };

  if (!isOpen) return null;

  return (
    <div className="absolute left-4 bottom-4 z-10 bg-gray-900 border border-gray-700 rounded-lg shadow-2xl max-h-[calc(100vh-120px)] overflow-hidden w-64 flex flex-col">
      {/* Header with Search */}
      <div className="sticky top-0 bg-gray-900 border-b border-gray-700 p-3">
        <div className="flex items-center justify-between mb-2">
          <div className="text-sm font-bold text-white">Add Nodes</div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
            title="Close (ESC)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={2}
              stroke="currentColor"
              className="w-5 h-5"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        
        {/* Search Input */}
        <input
          ref={searchInputRef}
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search nodes..."
          className="w-full px-3 py-1.5 text-sm bg-gray-800 border border-gray-600 rounded text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Categories and Nodes */}
      <div className="flex-1 overflow-y-auto">
        {filteredCategories.length === 0 ? (
          <div className="p-4 text-center text-gray-500 text-sm">
            No nodes found matching "{searchQuery}"
          </div>
        ) : (
          filteredCategories.map((category) => (
            <div key={category.name} className="border-b border-gray-800 last:border-b-0">
              <button
                onClick={() => toggleCategory(category.name)}
                className="w-full p-3 flex items-center justify-between hover:bg-gray-800 transition-colors"
              >
                <div className="text-xs font-semibold text-gray-400 uppercase tracking-wide">
                  {category.name}
                </div>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className={`w-4 h-4 text-gray-500 transition-transform ${
                    expandedCategories.has(category.name) ? 'rotate-180' : ''
                  }`}
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M19.5 8.25l-7.5 7.5-7.5-7.5"
                  />
                </svg>
              </button>
              
              {expandedCategories.has(category.name) && (
                <div className="px-3 pb-3 flex flex-col gap-1">
                  {category.nodes.map((config) => (
                    <button
                      key={config.type}
                      onClick={() => handleAddNode(config.type, config.defaultData)}
                      className="bg-gray-800 hover:bg-gray-700 text-white px-3 py-2 rounded text-sm transition-all text-left flex items-center gap-2"
                    >
                      <span className="text-xs">+</span>
                      <span>{config.label}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
      </div>
      
      {/* Footer hint */}
      <div className="border-t border-gray-800 px-3 py-2 text-xs text-gray-500">
        Press <kbd className="px-1 py-0.5 bg-gray-800 rounded">ESC</kbd> to close
      </div>
    </div>
  );
}
