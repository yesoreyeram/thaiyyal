import React from "react";
import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TextInputNodeData = {
  text?: string;
  label?: string;
};

export function TextInputNode(props: NodePropsWithOptions<TextInputNodeData>) {
  const { id, data, onShowOptions } = props;
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
      id={id}
      title={String(data?.label || "Text Input")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <textarea
        value={String(data?.text ?? "")}
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
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
