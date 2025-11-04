/**
 * TransposeNode Component
 * 
 * Transposes a matrix (2D array).
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TransposeNodeData = {
  label?: string;
};

/**
 * TransposeNode React Component
 * 
 * This component renders a visual node in the workflow editor that transposes a matrix
 * 
 * @param {NodePropsWithOptions<TransposeNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {TransposeNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function TransposeNode({ id, data, onShowOptions }: NodePropsWithOptions<TransposeNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("transposeNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Transpose")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-red-600 to-red-700 text-white shadow-lg rounded-lg border border-red-500 hover:border-red-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs py-1">Transpose matrix</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
