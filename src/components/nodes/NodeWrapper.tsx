import React, { ReactNode, useState, useCallback, useRef } from "react";
import { NodeTopBar } from "./NodeTopBar";
import { NodeInfoPopup } from "./NodeInfoPopup";
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
  className?: string;
  enableResize?: boolean;
}

export function NodeWrapper({
  title,
  children,
  nodeInfo,
  onShowOptions,
  className = "",
  enableResize = true,
}: NodeWrapperProps) {
  const [showInfo, setShowInfo] = useState(false);
  const [infoPosition, setInfoPosition] = useState({ x: 0, y: 0 });
  const nodeRef = useRef<HTMLDivElement>(null);

  const handleShowInfo = useCallback((e?: React.MouseEvent) => {
    if (e) {
      const rect = (e.target as HTMLElement).getBoundingClientRect();
      setInfoPosition({ x: rect.right + 8, y: rect.top });
    }
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
      <div className="px-3 py-2">
        <NodeTopBar
          title={title}
          onInfo={nodeInfo ? () => handleShowInfo() : undefined}
          onOptions={handleShowOptions}
        />
        {children}
      </div>
      {enableResize && <NodeResizeHandle />}
      {showInfo && nodeInfo && (
        <NodeInfoPopup
          title={title}
          description={nodeInfo.description}
          inputs={nodeInfo.inputs}
          outputs={nodeInfo.outputs}
          onClose={handleCloseInfo}
          x={infoPosition.x}
          y={infoPosition.y}
        />
      )}
    </div>
  );
}
