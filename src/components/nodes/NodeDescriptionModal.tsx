import React, { useRef, useEffect } from "react";
import { createPortal } from "react-dom";

interface NodeDescriptionModalProps {
  title: string;
  description: string;
  inputs?: string[];
  outputs?: string[];
  onClose: () => void;
}

export function NodeDescriptionModal({
  title,
  description,
  inputs = [],
  outputs = [],
  onClose,
}: NodeDescriptionModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
        onClose();
      }
    };
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [onClose]);

  const modalContent = (
    <div className="fixed inset-0 bg-black/30 backdrop-blur-[1px] flex items-center justify-center z-50">
      <div
        ref={modalRef}
        className="bg-gray-900 border border-gray-700 rounded-lg shadow-2xl overflow-hidden w-96"
        role="dialog"
        aria-label="Node information"
      >
        {/* Header with different background and border */}
        <div className="flex items-start justify-between px-4 py-3 bg-gray-800 border-b border-gray-700">
          <h3 className="text-sm font-semibold text-white">{title}</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-200 transition-colors"
            aria-label="Close"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="w-4 h-4"
            >
              <path d="M5.28 4.22a.75.75 0 0 0-1.06 1.06L6.94 8l-2.72 2.72a.75.75 0 1 0 1.06 1.06L8 9.06l2.72 2.72a.75.75 0 1 0 1.06-1.06L9.06 8l2.72-2.72a.75.75 0 0 0-1.06-1.06L8 6.94 5.28 4.22Z" />
            </svg>
          </button>
        </div>
        
        {/* Content */}
        <div className="p-4">
          <p className="text-xs text-gray-300 mb-3 leading-relaxed">{description}</p>
          {inputs.length > 0 && (
            <div className="mb-3">
              <h4 className="text-xs font-semibold text-gray-200 mb-1.5 px-2 py-1 bg-gray-800 border border-gray-700 rounded">
                Inputs
              </h4>
              <ul className="text-xs text-gray-300 space-y-1 mt-1.5">
                {inputs.map((input, i) => (
                  <li key={i} className="flex items-start gap-1.5">
                    <span className="text-blue-400 flex-shrink-0">▸</span>
                    <span>{input}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
          {outputs.length > 0 && (
            <div>
              <h4 className="text-xs font-semibold text-gray-200 mb-1.5 px-2 py-1 bg-gray-800 border border-gray-700 rounded">
                Outputs
              </h4>
              <ul className="text-xs text-gray-300 space-y-1 mt-1.5">
                {outputs.map((output, i) => (
                  <li key={i} className="flex items-start gap-1.5">
                    <span className="text-green-400 flex-shrink-0">▸</span>
                    <span>{output}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  );

  // Render modal at document body level using portal
  return typeof document !== 'undefined' 
    ? createPortal(modalContent, document.body)
    : null;
}
