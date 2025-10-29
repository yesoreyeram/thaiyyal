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
    <div className="p-2 bg-green-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Text Input</div>
      <input
        value={data?.text ?? ""}
        type="text"
        onChange={onChange}
        className="mt-1 w-48 border px-2 py-1 rounded text-black"
        placeholder="Enter text..."
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
