import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";


type FilterNodeData = {
  condition?: string;
  label?: string;
};

export function FilterNode({ id, data, onShowOptions }: NodePropsWithOptions<FilterNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
      )
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("filterNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Filter")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-purple-600 to-purple-700 text-white shadow-lg rounded-lg border border-purple-500 hover:border-purple-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.condition ?? "item.age > 0")}
        type="text"
        onChange={onConditionChange}
        className="w-32 text-xs border border-purple-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-purple-400 focus:outline-none"
        placeholder="item.age >= 18"
        aria-label="Filter condition"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
