import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type HttpNodeData = {
  url?: string;
  label?: string;
};

export function HttpNode({ id, data, ...props }: NodeProps<HttpNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const url = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, url } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("httpNode");
  const onShowOptions = (props as any).onShowOptions;

  return (
    <NodeWrapper
      title={String(data?.label || "HTTP Request")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className="bg-gradient-to-br from-purple-700 to-purple-800 text-white shadow-lg rounded-lg border border-purple-600 hover:border-purple-500 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.url ?? "")}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-purple-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-400 focus:outline-none"
        placeholder="https://..."
        aria-label="HTTP URL"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
