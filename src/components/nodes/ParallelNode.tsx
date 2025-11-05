/**
 * ParallelNode Component
 *
 * Executes multiple paths simultaneously.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ParallelNodeData = {
  max_concurrency?: number;
  label?: string;
};

/**
 * ParallelNode React Component
 *
 * This component renders a visual node in the workflow editor that executes paths simultaneously
 *
 * @param {NodePropsWithOptions<ParallelNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ParallelNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ParallelNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<ParallelNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_concurrency = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_concurrency } } : n
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

  const nodeInfo = getNodeInfo("parallelNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Parallel")}
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
        value={Number(data?.max_concurrency ?? 10)}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Max concurrency"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
