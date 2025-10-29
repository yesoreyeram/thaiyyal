import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";

type ConditionNodeData = {
  condition?: string;
  true_path?: string;
  false_path?: string;
  label?: string;
};

export function ConditionNode({ id, data }: NodeProps<ConditionNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-amber-600 to-amber-700 text-white shadow-lg rounded-lg border border-amber-500 hover:border-amber-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Condition")}</div>
      <input
        value={String(data?.condition ?? ">0")}
        type="text"
        onChange={onConditionChange}
        className="w-28 text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-amber-400 focus:outline-none"
        placeholder=">100, <50..."
        aria-label="Condition expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}
