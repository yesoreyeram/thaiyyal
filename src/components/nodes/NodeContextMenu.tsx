import React, { useRef, useEffect } from "react";

interface NodeContextMenuProps {
  x: number;
  y: number;
  onClose: () => void;
  onDelete: () => void;
  onDuplicate?: () => void;
  onCopy?: () => void;
}

export function NodeContextMenu({
  x,
  y,
  onClose,
  onDelete,
  onDuplicate,
  onCopy,
}: NodeContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
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

  return (
    <div
      ref={menuRef}
      className="fixed bg-gray-800 border border-gray-700 rounded-lg shadow-2xl py-1 z-50 min-w-[160px]"
      style={{ left: x, top: y }}
      role="menu"
    >
      {onDuplicate && (
        <button
          onClick={() => {
            onDuplicate();
            onClose();
          }}
          className="w-full px-3 py-2 text-left text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors flex items-center gap-2"
          role="menuitem"
        >
          <span>📋</span>
          <span>Duplicate</span>
        </button>
      )}
      {onCopy && (
        <button
          onClick={() => {
            onCopy();
            onClose();
          }}
          className="w-full px-3 py-2 text-left text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors flex items-center gap-2"
          role="menuitem"
        >
          <span>📄</span>
          <span>Copy</span>
        </button>
      )}
      <div className="border-t border-gray-700 my-1"></div>
      <button
        onClick={() => {
          onDelete();
          onClose();
        }}
        className="w-full px-3 py-2 text-left text-sm text-red-400 hover:bg-red-900/20 hover:text-red-300 transition-colors flex items-center gap-2"
        role="menuitem"
      >
        <span>🗑️</span>
        <span>Delete</span>
      </button>
    </div>
  );
}
