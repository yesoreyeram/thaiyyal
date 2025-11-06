import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type JoinNodeData = {
  join_strategy?: string;
  timeout?: string;
  label?: string;
};

export function JoinNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<JoinNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const join_strategy = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, join_strategy } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("joinNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Join")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={String(data?.join_strategy ?? "all")}
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="all">All</option>
        <option value="any">Any</option>
        <option value="first">First</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
