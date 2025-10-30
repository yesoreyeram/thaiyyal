import React, { useState, useRef, useEffect } from "react";

interface NodeTopBarProps {
  title: string;
  onInfo?: () => void;
  onOptions?: (x: number, y: number) => void;
  compact?: boolean;
}

export function NodeTopBar({ title, onInfo, onOptions, compact = true }: NodeTopBarProps) {
  const optionsRef = useRef<HTMLButtonElement>(null);

  const handleOptionsClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    if (onOptions && optionsRef.current) {
      const rect = optionsRef.current.getBoundingClientRect();
      onOptions(rect.right, rect.bottom);
    }
  };

  const handleInfoClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    onInfo?.();
  };

  return (
    <div className="flex items-center justify-between gap-1 mb-2 pb-2 border-b border-white/10 bg-black/20 -mx-3 -mt-2 px-3 pt-2 rounded-t-lg">
      <div className="text-xs font-semibold text-gray-100 truncate flex-1 min-w-0">
        {title}
      </div>
      <div className="flex items-center gap-0.5 flex-shrink-0">
        {onInfo && (
          <button
            onClick={handleInfoClick}
            className="w-5 h-5 flex items-center justify-center rounded hover:bg-white/10 transition-colors text-gray-300 hover:text-white"
            aria-label="Show node information"
            title="Show information"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="w-3.5 h-3.5"
            >
              <path
                fillRule="evenodd"
                d="M15 8A7 7 0 1 1 1 8a7 7 0 0 1 14 0ZM9 5a1 1 0 1 1-2 0 1 1 0 0 1 2 0ZM6.75 8a.75.75 0 0 0 0 1.5h.75v1.75a.75.75 0 0 0 1.5 0v-2.5A.75.75 0 0 0 8.25 8h-1.5Z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        )}
        {onOptions && (
          <button
            ref={optionsRef}
            onClick={handleOptionsClick}
            className="w-5 h-5 flex items-center justify-center rounded hover:bg-white/10 transition-colors text-gray-300 hover:text-white"
            aria-label="Show node options"
            title="Show options"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="w-3.5 h-3.5"
            >
              <path d="M2 8a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM6.5 8a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM12.5 6.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3Z" />
            </svg>
          </button>
        )}
      </div>
    </div>
  );
}
