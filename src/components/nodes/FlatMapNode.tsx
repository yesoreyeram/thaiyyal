/**
 * FlatMapNode Component
 * 
 * Maps each element and flattens the result into a single array.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type FlatMapNodeData = {
  expression?: string;
  label?: string;
};

/**
 * FlatMapNode React Component
 * 
 * This component renders a visual node in the workflow editor that maps and flattens array elements
 * 
 * @param {NodePropsWithOptions<FlatMapNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {FlatMapNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function FlatMapNode({ id, data, onShowOptions }: NodePropsWithOptions<FlatMapNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
      )
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("flatMapNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "FlatMap")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-indigo-600 to-indigo-700 text-white shadow-lg rounded-lg border border-indigo-500 hover:border-indigo-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item.values")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-indigo-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
        placeholder="item.values"
        aria-label="FlatMap expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
