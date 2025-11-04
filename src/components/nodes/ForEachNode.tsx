/**
 * ForEachNode Component
 * 
 * Executes a subflow for each element in an input array.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ForEachNodeData = {
  max_iterations?: number;
  label?: string;
};

/**
 * ForEachNode React Component
 * 
 * This component renders a visual node in the workflow editor that iterates over array elements
 * 
 * @param {NodePropsWithOptions<ForEachNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ForEachNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ForEachNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<ForEachNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_iterations = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_iterations } } : n
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

  const nodeInfo = getNodeInfo("forEachNode");

  return (
    <NodeWrapper
      title={String(data?.label || "For Each")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={Number(data?.max_iterations ?? 1000)}
        type="number"
        onChange={onChange}
        className="w-20 text-[10px] leading-tight border border-gray-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="Max iter"
        aria-label="Max iterations"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
