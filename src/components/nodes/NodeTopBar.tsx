import React, { useState, useRef, useEffect } from "react";

interface NodeTopBarProps {
  title: string;
  onInfo?: () => void;
  onOptions?: (x: number, y: number) => void;
  onTitleChange?: (newTitle: string) => void;
  compact?: boolean;
}

export function NodeTopBar({ title, onInfo, onOptions, onTitleChange, compact = true }: NodeTopBarProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(title);
  const optionsRef = useRef<HTMLButtonElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setEditValue(title);
  }, [title]);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

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

  const handleTitleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (onTitleChange) {
      setIsEditing(true);
    }
  };

  const handleTitleSubmit = () => {
    if (editValue.trim() && onTitleChange) {
      onTitleChange(editValue.trim());
    } else {
      setEditValue(title);
    }
    setIsEditing(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    e.stopPropagation();
    if (e.key === "Enter") {
      handleTitleSubmit();
    } else if (e.key === "Escape") {
      setEditValue(title);
      setIsEditing(false);
    }
  };

  return (
    <div className="flex items-center justify-between gap-1 mb-0.5 pb-0.5 border-b border-white/10 bg-black/20 -mx-2 -mt-1 px-1.5 pt-0.5 rounded-t">
      {isEditing ? (
        <input
          ref={inputRef}
          type="text"
          value={editValue}
          onChange={(e) => setEditValue(e.target.value)}
          onBlur={handleTitleSubmit}
          onKeyDown={handleKeyDown}
          className="text-[10px] leading-tight font-medium text-gray-100 flex-1 min-w-0 bg-gray-900 border border-gray-600 rounded px-1 py-0.5 focus:outline-none focus:ring-1 focus:ring-blue-500"
          onClick={(e) => e.stopPropagation()}
        />
      ) : (
        <div 
          className={`text-[10px] leading-tight font-medium text-gray-100 truncate flex-1 min-w-0 ${onTitleChange ? 'cursor-text hover:text-white' : ''}`}
          onClick={handleTitleClick}
          title={onTitleChange ? "Click to edit title" : title}
        >
          {title}
        </div>
      )}
      <div className="flex items-center gap-0.5 flex-shrink-0">
        {onInfo && (
          <button
            onClick={handleInfoClick}
            className="w-3.5 h-3.5 flex items-center justify-center rounded hover:bg-white/10 transition-colors text-gray-300 hover:text-white"
            aria-label="Show node information"
            title="Show information"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="w-2.5 h-2.5"
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
            className="w-3.5 h-3.5 flex items-center justify-center rounded hover:bg-white/10 transition-colors text-gray-300 hover:text-white"
            aria-label="Show node options"
            title="Show options"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="w-2.5 h-2.5"
            >
              <path d="M2 8a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM6.5 8a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM12.5 6.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3Z" />
            </svg>
          </button>
        )}
      </div>
    </div>
  );
}
