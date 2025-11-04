import React, { useCallback } from "react";
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
  const handleResize = useCallback(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (_event: any, params: any) => {
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
      onResize={handleResize}
      lineClassName="border-blue-400"
      handleClassName="w-2 h-2 bg-blue-400 rounded-full"
    />
  );
}
