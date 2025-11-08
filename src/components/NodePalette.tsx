"use client";
import React, { useState, useEffect, useRef } from "react";

export interface NodeConfig {
  type: string;
  label: string;
  color: string;
  defaultData: Record<string, unknown>;
}

export interface NodeCategory {
  name: string;
  nodes: NodeConfig[];
}

interface NodePaletteProps {
  isOpen: boolean;
  onClose: () => void;
  categories: NodeCategory[];
  onAddNode: (type: string, defaultData: Record<string, unknown>) => void;
}

export function NodePalette({
  isOpen,
  onClose,
  categories,
  onAddNode,
}: NodePaletteProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(
    new Set(categories.map((c) => c.name))
  );
  const searchInputRef = useRef<HTMLInputElement>(null);
  const paletteRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isOpen) {
      // Focus search input when palette opens
      setTimeout(() => searchInputRef.current?.focus(), 100);
    }
  }, [isOpen]);

  const toggleCategory = (categoryName: string) => {
    setExpandedCategories((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(categoryName)) {
        newSet.delete(categoryName);
      } else {
        newSet.add(categoryName);
      }
      return newSet;
    });
  };

  const filteredCategories = categories
    .map((category) => ({
      ...category,
      nodes: category.nodes.filter(
        (node) =>
          node.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
          node.type.toLowerCase().includes(searchQuery.toLowerCase())
      ),
    }))
    .filter((category) => category.nodes.length > 0);

  const handleAddNode = (
    type: string,
    defaultData: Record<string, unknown>
  ) => {
    onAddNode(type, defaultData);
    setSearchQuery(""); // Clear search after adding
  };

  const onDragStart = (
    event: React.DragEvent,
    type: string,
    defaultData: Record<string, unknown>
  ) => {
    event.dataTransfer.setData(
      "application/reactflow",
      JSON.stringify({ type, defaultData })
    );
    event.dataTransfer.effectAllowed = "move";
  };

  if (!isOpen) return null;

  return (
    <div
      ref={paletteRef}
      className="h-full bg-white dark:bg-black border-r border-gray-300 dark:border-gray-700 w-64 flex flex-col"
    >
      {/* Header with Search */}
      <div className="bg-white dark:bg-black border-b border-gray-300 dark:border-gray-700 p-3">
        <div className="flex items-center justify-between mb-2">
          <div className="text-sm font-bold text-black dark:text-white">Nodes</div>
          <button
            onClick={onClose}
            className="text-gray-600 dark:text-gray-400 hover:text-black dark:hover:text-white transition-colors"
            title="Close Sidebar"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={2}
              stroke="currentColor"
              className="w-4 h-4"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M18.75 19.5l-7.5-7.5 7.5-7.5m-6 15L5.25 12l7.5-7.5"
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
          className="w-full px-3 py-1.5 text-sm bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 rounded text-black dark:text-white placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-gray-400 dark:focus:ring-gray-600"
        />
      </div>

      {/* Categories and Nodes */}
      <div className="flex-1 overflow-y-auto">
        {filteredCategories.length === 0 ? (
          <div className="p-4 text-center text-gray-500 text-sm">
            No nodes found matching &ldquo;{searchQuery}&rdquo;
          </div>
        ) : (
          filteredCategories.map((category) => (
            <div
              key={category.name}
              className="border-b border-gray-200 dark:border-gray-800 last:border-b-0"
            >
              <button
                onClick={() => toggleCategory(category.name)}
                className="w-full p-3 flex items-center justify-between hover:bg-gray-100 dark:hover:bg-gray-900 transition-colors"
              >
                <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide">
                  {category.name}
                </div>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className={`w-4 h-4 text-gray-500 transition-transform ${
                    expandedCategories.has(category.name) ? "rotate-180" : ""
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
                    <div
                      key={config.type + config.label}
                      draggable
                      onDragStart={(e) =>
                        onDragStart(e, config.type, config.defaultData)
                      }
                      onClick={() =>
                        handleAddNode(config.type, config.defaultData)
                      }
                      className="bg-gray-100 dark:bg-gray-900 hover:bg-gray-200 dark:hover:bg-gray-800 text-black dark:text-white px-3 py-2 rounded text-sm transition-all text-left flex items-center gap-2 cursor-move border border-gray-300 dark:border-gray-700"
                      title="Click to add or drag to canvas"
                    >
                      <span className="text-xs">⋮⋮</span>
                      <span>{config.label}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Footer hint */}
      <div className="border-t border-gray-300 dark:border-gray-800 px-3 py-2 text-xs text-gray-600 dark:text-gray-500">
        Click to add • Drag to place
      </div>
    </div>
  );
}
