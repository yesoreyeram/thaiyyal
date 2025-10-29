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
    <div className="p-2 bg-yellow-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Condition</div>
      <input
        value={data?.condition ?? ">0"}
        type="text"
        onChange={onConditionChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder=">100, <50..."
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
