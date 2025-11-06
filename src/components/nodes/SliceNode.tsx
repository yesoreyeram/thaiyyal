/**
 * SliceNode Component
 *
 * Extracts a portion of an array from start to end indices.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SliceNodeData = {
  start?: number;
  end?: number;
  label?: string;
};

/**
 * SliceNode React Component
 *
 * This component renders a visual node in the workflow editor that extracts a portion of an array
 *
 * @param {NodePropsWithOptions<SliceNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {SliceNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function SliceNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SliceNodeData>) {
  const { setNodes } = useReactFlow();

  const onStartChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const start = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, start } } : n))
    );
  };

  const onEndChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const end = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, end } } : n))
    );
  };

  const nodeInfo = getNodeInfo("sliceNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Slice")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex items-center gap-0.5 w-36">
        <input
          value={Number(data?.start ?? 0)}
          type="number"
          onChange={onStartChange}
          className="w-14 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Start"
          aria-label="Start index"
        />
        <span className="text-[10px]">to</span>
        <input
          value={Number(data?.end ?? -1)}
          type="number"
          onChange={onEndChange}
          className="w-14 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="End"
          aria-label="End index"
        />
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
