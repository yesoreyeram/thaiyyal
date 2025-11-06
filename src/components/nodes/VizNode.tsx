"use client";
import React from "react";
import { useReactFlow, Handle, Position } from "reactflow";
import { NodePropsWithOptions, getNodeInfo, NodeWrapper } from "./";

export function VizNode(props: NodePropsWithOptions) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const mode = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, mode } } : n))
    );
  };

  const nodeInfo = getNodeInfo("vizNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Visualization")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={typeof data?.mode === "string" ? data.mode : "text"}
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="text">Text</option>
        <option value="table">Table</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
