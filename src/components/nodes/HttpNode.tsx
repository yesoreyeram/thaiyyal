import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";

type HttpNodeData = {
  url?: string;
  label?: string;
};

export function HttpNode({ id, data }: NodeProps<HttpNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const url = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, url } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-purple-700 to-purple-800 text-white shadow-lg rounded-lg border border-purple-600 hover:border-purple-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "HTTP Request")}</div>
      <input
        value={String(data?.url ?? "")}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-purple-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-400 focus:outline-none"
        placeholder="https://..."
        aria-label="HTTP URL"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}
