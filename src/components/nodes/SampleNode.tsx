import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SampleNodeData = {
  count?: number;
  label?: string;
};

export function SampleNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SampleNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const count = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, count } } : n))
    );
  };

  const nodeInfo = getNodeInfo("sampleNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Sample")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={Number(data?.count ?? 1)}
        type="number"
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="Count"
        aria-label="Sample count"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
