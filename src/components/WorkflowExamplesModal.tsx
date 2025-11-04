"use client";
import React, { useState, useMemo } from "react";
import { workflowExamples, WorkflowExample } from "../data/workflowExamples";

interface WorkflowExamplesModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (example: WorkflowExample) => void;
}

export function WorkflowExamplesModal({
  isOpen,
  onClose,
  onSelect,
}: WorkflowExamplesModalProps) {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  // Filter examples based on search term and selected tags
  const filteredExamples = useMemo(() => {
    return workflowExamples.filter((example) => {
      // Filter by search term
      const searchLower = searchTerm.toLowerCase();
      const matchesSearch =
        !searchTerm ||
        example.title.toLowerCase().includes(searchLower) ||
        example.description.toLowerCase().includes(searchLower) ||
        example.tags.some((tag) => tag.toLowerCase().includes(searchLower));

      // Filter by selected tags
      const matchesTags =
        selectedTags.length === 0 ||
        selectedTags.every((tag) => example.tags.includes(tag));

      return matchesSearch && matchesTags;
    });
  }, [searchTerm, selectedTags]);

  const toggleTag = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    );
  };

  const clearFilters = () => {
    setSearchTerm("");
    setSelectedTags([]);
  };

  const handleSelect = (example: WorkflowExample) => {
    onSelect(example);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-6xl max-h-[90vh] flex flex-col overflow-hidden">
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-700 flex items-center justify-between bg-gray-800/50">
          <div className="flex items-center gap-3">
            <span className="text-2xl">📚</span>
            <h2 className="text-xl font-semibold text-white">
              Workflow Examples Library
            </h2>
          </div>

          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm font-medium rounded-lg transition-colors"
            aria-label="Close modal"
          >
            ✕
          </button>
        </div>

        {/* Search and Filter Bar */}
        <div className="px-6 py-4 border-b border-gray-700 bg-gray-800/30">
          {/* Search Input */}
          <div className="mb-3">
            <input
              type="text"
              placeholder="Search examples by title, description, or tags..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full px-4 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Popular Tags */}
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm text-gray-400 mr-2">Popular tags:</span>
            {["api", "http", "json", "data", "branching", "variables", "error-handling", "performance"].map(
              (tag) => (
                <button
                  key={tag}
                  onClick={() => toggleTag(tag)}
                  className={`px-2.5 py-1 rounded-full text-xs font-medium transition-colors ${
                    selectedTags.includes(tag)
                      ? "bg-blue-600 text-white"
                      : "bg-gray-700 text-gray-300 hover:bg-gray-600"
                  }`}
                >
                  {tag}
                </button>
              )
            )}
            {(searchTerm || selectedTags.length > 0) && (
              <button
                onClick={clearFilters}
                className="px-2.5 py-1 rounded-full text-xs font-medium bg-red-600/20 text-red-400 hover:bg-red-600/30 transition-colors ml-2"
              >
                Clear filters
              </button>
            )}
          </div>
        </div>

        {/* Results Info */}
        <div className="px-6 py-2 bg-gray-800/20 text-sm text-gray-400">
          Showing {filteredExamples.length} of {workflowExamples.length} examples
          {selectedTags.length > 0 && (
            <span className="ml-2">
              (filtered by:{" "}
              {selectedTags.map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center gap-1 ml-1 px-2 py-0.5 bg-blue-600/20 text-blue-400 rounded"
                >
                  {tag}
                  <button
                    onClick={() => toggleTag(tag)}
                    className="hover:text-blue-300"
                  >
                    ×
                  </button>
                </span>
              ))}
              )
            </span>
          )}
        </div>

        {/* Content - Examples Grid */}
        <div className="flex-1 overflow-auto p-6 custom-scrollbar">
          {filteredExamples.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-6xl mb-4 opacity-50">🔍</div>
              <p className="text-gray-400 text-lg mb-2">No examples found</p>
              <p className="text-gray-600 text-sm">
                Try adjusting your search or filters
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredExamples.map((example) => (
                <div
                  key={example.id}
                  onClick={() => handleSelect(example)}
                  className="group p-4 bg-gray-800 hover:bg-gray-750 border border-gray-700 hover:border-blue-600/50 rounded-lg cursor-pointer transition-all hover:shadow-lg hover:shadow-blue-900/20"
                >
                  {/* Title */}
                  <h3 className="text-white font-medium mb-2 truncate group-hover:text-blue-400 transition-colors">
                    {example.title}
                  </h3>

                  {/* Description */}
                  <p className="text-gray-400 text-sm mb-3 line-clamp-3">
                    {example.description}
                  </p>

                  {/* Tags */}
                  <div className="flex flex-wrap gap-1.5">
                    {example.tags.slice(0, 4).map((tag) => (
                      <span
                        key={tag}
                        className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                          selectedTags.includes(tag)
                            ? "bg-blue-600/30 text-blue-300"
                            : "bg-gray-700 text-gray-400"
                        }`}
                      >
                        {tag}
                      </span>
                    ))}
                    {example.tags.length > 4 && (
                      <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-gray-700 text-gray-400">
                        +{example.tags.length - 4}
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-3 border-t border-gray-700 bg-gray-800/50 flex items-center justify-between text-xs text-gray-500">
          <span>
            {filteredExamples.length} example
            {filteredExamples.length !== 1 ? "s" : ""} available
          </span>
          <span>Click a card to load • Press ESC to close</span>
        </div>
      </div>
    </div>
  );
}
