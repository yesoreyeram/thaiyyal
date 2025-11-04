/**
 * ReverseNode Component
 * 
 * Reverses the order of array elements.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ReverseNodeData = {
  label?: string;
};

/**
 * ReverseNode React Component
 * 
 * This component renders a visual node in the workflow editor that reverses array element order
 * 
 * @param {NodePropsWithOptions<ReverseNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ReverseNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ReverseNode({ id, data, onShowOptions }: NodePropsWithOptions<ReverseNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("reverseNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Reverse")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-rose-600 to-rose-700 text-white shadow-lg rounded-lg border border-rose-500 hover:border-rose-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs py-1">Reverse array order</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
