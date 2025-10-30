import React, { useState, useCallback } from "react";
import { NodeResizer } from "reactflow";

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
  const [isVisible, setIsVisible] = useState(false);

  const handleResize = useCallback(
    (event: any, params: any) => {
      if (onResize) {
        onResize(params.width, params.height);
      }
    },
    [onResize]
  );

  return (
    <NodeResizer
      minWidth={minWidth}
      minHeight={minHeight}
      maxWidth={maxWidth}
      maxHeight={maxHeight}
      isVisible={isVisible}
      onResize={handleResize}
      lineClassName="border-blue-400"
      handleClassName="w-2 h-2 bg-blue-400 rounded-full"
    />
  );
}
