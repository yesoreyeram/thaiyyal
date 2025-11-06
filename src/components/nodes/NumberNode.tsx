"use client";
import React from "react";
import { useReactFlow, Handle, Position } from "reactflow";
import { NodePropsWithOptions, getNodeInfo, NodeWrapper } from "./";

export function NumberNode({
  id,
  data,
  onShowOptions,
  onOpenInfo,
}: NodePropsWithOptions) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const v = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, value: v } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("numberNode");
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Number")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onOpenInfo={onOpenInfo}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={typeof data?.value === "number" ? data.value : 0}
        type="number"
        onChange={onChange}
        className="w-full text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
