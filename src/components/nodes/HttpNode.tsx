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
    <div className="p-2 bg-purple-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">HTTP Request</div>
      <input
        value={data?.url ?? ""}
        type="text"
        onChange={onChange}
        className="mt-1 w-48 border px-2 py-1 rounded text-black text-sm"
        placeholder="https://..."
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
