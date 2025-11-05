import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type DateInputNodeData = {
  value?: string;
  label?: string;
};

export function DateInputNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<DateInputNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, value } } : n))
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("dateInputNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Date")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-cyan-700 to-cyan-800 text-white shadow-lg rounded-lg border border-cyan-600 hover:border-cyan-500 transition-all"
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.value ?? "")}
        type="date"
        onChange={onChange}
        className="w-36 text-xs border border-cyan-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-cyan-400 focus:outline-none"
        aria-label="Date value"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
