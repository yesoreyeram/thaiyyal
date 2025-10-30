import React, { useState, useCallback, useRef, useEffect } from "react";

interface NodeResizeHandleProps {
  onResize?: (width: number, height: number) => void;
  minWidth?: number;
  minHeight?: number;
  maxWidth?: number;
  maxHeight?: number;
}

export function NodeResizeHandle({
  onResize,
  minWidth = 150,
  minHeight = 80,
  maxWidth = 500,
  maxHeight = 500,
}: NodeResizeHandleProps) {
  const [isDragging, setIsDragging] = useState(false);
  const startPosRef = useRef({ x: 0, y: 0 });
  const startSizeRef = useRef({ width: 0, height: 0 });
  const nodeRef = useRef<HTMLDivElement | null>(null);

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
    startPosRef.current = { x: e.clientX, y: e.clientY };
    
    // Find the parent node element
    const nodeElement = (e.target as HTMLElement).closest('.react-flow__node');
    if (nodeElement) {
      nodeRef.current = nodeElement as HTMLDivElement;
      const rect = nodeElement.getBoundingClientRect();
      startSizeRef.current = { width: rect.width, height: rect.height };
    }
  }, []);

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      if (!nodeRef.current) return;
      
      const deltaX = e.clientX - startPosRef.current.x;
      const deltaY = e.clientY - startPosRef.current.y;
      
      let newWidth = startSizeRef.current.width + deltaX;
      let newHeight = startSizeRef.current.height + deltaY;
      
      // Apply constraints
      newWidth = Math.max(minWidth, Math.min(maxWidth, newWidth));
      newHeight = Math.max(minHeight, Math.min(maxHeight, newHeight));
      
      // Apply styles directly for smooth resizing
      nodeRef.current.style.width = `${newWidth}px`;
      nodeRef.current.style.height = `${newHeight}px`;
    };

    const handleMouseUp = () => {
      setIsDragging(false);
      if (nodeRef.current && onResize) {
        const rect = nodeRef.current.getBoundingClientRect();
        onResize(rect.width, rect.height);
      }
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging, minWidth, minHeight, maxWidth, maxHeight, onResize]);

  return (
    <div
      onMouseDown={handleMouseDown}
      className={`absolute bottom-0 right-0 w-4 h-4 cursor-nwse-resize group ${
        isDragging ? 'z-50' : 'z-10'
      }`}
      aria-label="Resize node"
      title="Drag to resize"
    >
      <div className="absolute bottom-0.5 right-0.5 w-3 h-3 opacity-40 group-hover:opacity-70 transition-opacity">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 16 16"
          fill="currentColor"
          className="w-full h-full text-gray-400"
        >
          <path d="M15 14.5a.5.5 0 0 1-.5.5h-1a.5.5 0 0 1 0-1h1a.5.5 0 0 1 .5.5ZM15 11.5a.5.5 0 0 1-.5.5h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 1 .5.5ZM15 8.5a.5.5 0 0 1-.5.5h-7a.5.5 0 0 1 0-1h7a.5.5 0 0 1 .5.5Z" />
        </svg>
      </div>
    </div>
  );
}
