import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";

type TextOperationNodeData = {
  text_op?: string;
  separator?: string;
  repeat_n?: number;
  label?: string;
};

export function TextOperationNode({ id, data }: NodeProps<TextOperationNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onOpChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const text_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, text_op } } : n))
    );
  };

  const onSeparatorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const separator = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, separator } } : n))
    );
  };

  const onRepeatNChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const repeat_n = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, repeat_n } } : n))
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-emerald-600 to-emerald-700 text-white shadow-lg rounded-lg border border-emerald-500 hover:border-emerald-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{data?.label || "Text Operation"}</div>
      <select
        value={data?.text_op ?? "uppercase"}
        onChange={onOpChange}
        className="w-32 text-xs border border-emerald-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-emerald-400 focus:outline-none"
        aria-label="Text operation type"
      >
        <option value="uppercase">Uppercase</option>
        <option value="lowercase">Lowercase</option>
        <option value="titlecase">Title Case</option>
        <option value="camelcase">Camel Case</option>
        <option value="inversecase">Inverse Case</option>
        <option value="concat">Concatenate</option>
        <option value="repeat">Repeat</option>
      </select>
      {data?.text_op === "concat" && (
        <input
          value={data?.separator ?? ""}
          type="text"
          onChange={onSeparatorChange}
          className="mt-1 w-32 text-xs border border-emerald-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-emerald-400 focus:outline-none"
          placeholder="Separator..."
          aria-label="Separator"
        />
      )}
      {data?.text_op === "repeat" && (
        <input
          value={data?.repeat_n ?? 1}
          type="number"
          onChange={onRepeatNChange}
          className="mt-1 w-32 text-xs border border-emerald-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-emerald-400 focus:outline-none"
          placeholder="Repeat count..."
          aria-label="Repeat count"
        />
      )}
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}
