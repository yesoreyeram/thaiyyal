import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TextInputNodeData = {
  text?: string;
  label?: string;
};

export function TextInputNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<TextInputNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, text } } : n))
    );
  };

  const nodeInfo = getNodeInfo("textInputNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Text Input")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className="bg-gradient-to-br from-emerald-700 to-emerald-800 text-white shadow-lg rounded-lg border border-emerald-600 hover:border-emerald-500 transition-all"
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <textarea
        value={String(data?.text ?? "")}
        onChange={onChange}
        className="w-36 text-xs border border-emerald-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-emerald-400 focus:outline-none resize-none"
        placeholder="Enter text..."
        aria-label="Text input value"
        rows={3}
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
