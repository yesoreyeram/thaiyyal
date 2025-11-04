/**
 * ChunkNode Component
 * 
 * Splits an array into chunks of specified size.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ChunkNodeData = {
  size?: number;
  label?: string;
};

/**
 * ChunkNode React Component
 * 
 * This component renders a visual node in the workflow editor that splits an array into chunks
 * 
 * @param {NodePropsWithOptions<ChunkNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ChunkNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ChunkNode({ id, data, onShowOptions }: NodePropsWithOptions<ChunkNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const size = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, size } } : n
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

  const nodeInfo = getNodeInfo("chunkNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Chunk")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-pink-600 to-pink-700 text-white shadow-lg rounded-lg border border-pink-500 hover:border-pink-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={Number(data?.size ?? 3)}
        type="number"
        onChange={onChange}
        className="w-20 text-xs border border-pink-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-400 focus:outline-none"
        placeholder="Size"
        aria-label="Chunk size"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
