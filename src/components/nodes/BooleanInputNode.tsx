import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type BooleanInputNodeData = {
  value?: boolean;
  label?: string;
};

export function BooleanInputNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<BooleanInputNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.checked;
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

  const nodeInfo = getNodeInfo("booleanInputNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Boolean")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-indigo-700 to-indigo-800 text-white shadow-lg rounded-lg border border-indigo-600 hover:border-indigo-500 transition-all"
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex items-center gap-2 px-1">
        <input
          checked={data?.value ?? false}
          type="checkbox"
          onChange={onChange}
          className="w-4 h-4 text-indigo-600 bg-gray-900 border-indigo-600 rounded focus:ring-indigo-500 focus:ring-2"
          aria-label="Boolean value"
        />
        <span className="text-xs text-gray-300">
          {data?.value ? "True" : "False"}
        </span>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
