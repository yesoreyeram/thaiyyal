/**
 * UniqueNode Component
 * 
 * Removes duplicate elements from an array.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type UniqueNodeData = {
  by_field?: string;
  label?: string;
};

/**
 * UniqueNode React Component
 * 
 * This component renders a visual node in the workflow editor that removes duplicate elements
 * 
 * @param {NodePropsWithOptions<UniqueNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {UniqueNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function UniqueNode({ id, data, onShowOptions }: NodePropsWithOptions<UniqueNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const by_field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, by_field } } : n
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

  const nodeInfo = getNodeInfo("uniqueNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Unique")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-fuchsia-600 to-fuchsia-700 text-white shadow-lg rounded-lg border border-fuchsia-500 hover:border-fuchsia-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.by_field ?? "")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-fuchsia-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-fuchsia-400 focus:outline-none"
        placeholder="Field (optional)"
        aria-label="Unique by field"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
