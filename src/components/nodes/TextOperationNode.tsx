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
    <div className="p-2 bg-green-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Text Operation</div>
      <select
        value={data?.text_op ?? "uppercase"}
        onChange={onOpChange}
        className="mt-1 w-40 border px-2 py-1 rounded text-black text-sm"
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
          className="mt-1 w-40 border px-2 py-1 rounded text-black text-sm"
          placeholder="Separator..."
        />
      )}
      {data?.text_op === "repeat" && (
        <input
          value={data?.repeat_n ?? 1}
          type="number"
          onChange={onRepeatNChange}
          className="mt-1 w-40 border px-2 py-1 rounded text-black text-sm"
          placeholder="Repeat count..."
        />
      )}
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
