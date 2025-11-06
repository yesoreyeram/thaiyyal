import { ReactNode, useState, useCallback, useRef } from "react";
import { useReactFlow } from "reactflow";
import { NodeTopBar } from "./NodeTopBar";
import { NodeDescriptionModal } from "./NodeDescriptionModal";
import { NodeResizeHandle } from "./NodeResizeHandle";

export interface NodeWrapperProps {
  id: string;
  title: string;
  children: ReactNode;
  nodeInfo?: {
    description: string;
    inputs?: string[];
    outputs?: string[];
  };
  onShowOptions?: (x: number, y: number) => void;
  onDelete?: () => void;
  className?: string;
  enableResize?: boolean;
  onOpenInfo?: () => void;
}

export function NodeWrapper(props: NodeWrapperProps) {
  const {
    id,
    title,
    children,
    nodeInfo,
    onShowOptions,
    onDelete,
    onOpenInfo,
    enableResize = false,
  } = props;
  const [showInfo, setShowInfo] = useState(false);
  const nodeRef = useRef<HTMLDivElement>(null);
  const { setNodes } = useReactFlow();

  const handleShowInfo = useCallback(() => {
    if (onOpenInfo) {
      onOpenInfo(); // Close palette if open
    }
    setShowInfo(true);
  }, [onOpenInfo]);

  const handleCloseInfo = useCallback(() => {
    setShowInfo(false);
  }, []);

  const handleShowOptions = useCallback(
    (x: number, y: number) => {
      if (nodeRef.current && onShowOptions) {
        // Get node position to calculate relative coordinates
        const nodeRect = nodeRef.current.getBoundingClientRect();
        // Position menu relative to node's bottom-right area
        const relativeX = nodeRect.left + nodeRect.width;
        const relativeY = y;
        onShowOptions(relativeX, relativeY);
      }
    },
    [onShowOptions]
  );

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  return (
    <div
      ref={nodeRef}
      className={`relative bg-gray-800 text-white shadow-lg rounded border border-gray-700 hover:border-gray-600 transition-all`}
    >
      <div className="px-2 py-1">
        <NodeTopBar
          title={title}
          onInfo={nodeInfo ? handleShowInfo : undefined}
          onOptions={handleShowOptions}
          onTitleChange={handleTitleChange}
          onDelete={onDelete}
        />
        {children}
      </div>
      {enableResize && <NodeResizeHandle onResize={() => {}} />}
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
