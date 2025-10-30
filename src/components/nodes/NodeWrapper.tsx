import React, { ReactNode, useState, useCallback, useRef } from "react";
import { NodeTopBar } from "./NodeTopBar";
import { NodeDescriptionModal } from "./NodeDescriptionModal";
import { NodeResizeHandle } from "./NodeResizeHandle";

export interface NodeWrapperProps {
  title: string;
  children: ReactNode;
  nodeInfo?: {
    description: string;
    inputs?: string[];
    outputs?: string[];
  };
  onShowOptions?: (x: number, y: number) => void;
  onTitleChange?: (newTitle: string) => void;
  className?: string;
  enableResize?: boolean;
}

export function NodeWrapper({
  title,
  children,
  nodeInfo,
  onShowOptions,
  onTitleChange,
  className = "",
  enableResize = true,
}: NodeWrapperProps) {
  const [showInfo, setShowInfo] = useState(false);
  const nodeRef = useRef<HTMLDivElement>(null);

  const handleShowInfo = useCallback(() => {
    setShowInfo(true);
  }, []);

  const handleCloseInfo = useCallback(() => {
    setShowInfo(false);
  }, []);

  const handleShowOptions = useCallback((x: number, y: number) => {
    if (nodeRef.current && onShowOptions) {
      // Get node position to calculate relative coordinates
      const nodeRect = nodeRef.current.getBoundingClientRect();
      // Position menu relative to node's bottom-right area
      const relativeX = nodeRect.left + nodeRect.width;
      const relativeY = y;
      onShowOptions(relativeX, relativeY);
    }
  }, [onShowOptions]);

  return (
    <div ref={nodeRef} className={`relative ${className}`}>
      <div className="px-2.5 py-1.5">
        <NodeTopBar
          title={title}
          onInfo={nodeInfo ? handleShowInfo : undefined}
          onOptions={handleShowOptions}
          onTitleChange={onTitleChange}
        />
        {children}
      </div>
      {enableResize && <NodeResizeHandle />}
      {showInfo && nodeInfo && (
        <NodeDescriptionModal
          title={title}
          description={nodeInfo.description}
          inputs={nodeInfo.inputs}
          outputs={nodeInfo.outputs}
          onClose={handleCloseInfo}
        />
      )}
    </div>
  );
}
