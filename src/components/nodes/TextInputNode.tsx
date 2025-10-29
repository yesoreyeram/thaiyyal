import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";

type TextInputNodeData = {
  text?: string;
  label?: string;
};

export function TextInputNode({ id, data }: NodeProps<TextInputNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, text } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-emerald-700 to-emerald-800 text-white shadow-lg rounded-lg border border-emerald-600 hover:border-emerald-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{data?.label || "Text Input"}</div>
      <input
        value={data?.text ?? ""}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-emerald-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-emerald-400 focus:outline-none"
        placeholder="Enter text..."
        aria-label="Text input value"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}
