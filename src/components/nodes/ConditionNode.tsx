import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ConditionNodeData = {
  condition?: string;
  true_path?: string;
  false_path?: string;
  label?: string;
};

export function ConditionNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<ConditionNodeData>) {
  const { setNodes } = useReactFlow();

  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("conditionNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Condition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.condition ?? ">0")}
        type="text"
        onChange={onConditionChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="node.id.value > 100"
        aria-label="Condition expression"
      />
      {/* True path handle (top right) */}
      <Handle
        type="source"
        position={Position.Right}
        id="true"
        style={{ top: "30%" }}
        className="w-2 h-2 bg-green-500"
        title="True path"
      />
      {/* False path handle (bottom right) */}
      <Handle
        type="source"
        position={Position.Right}
        id="false"
        style={{ top: "70%" }}
        className="w-2 h-2 bg-red-500"
        title="False path"
      />
    </NodeWrapper>
  );
}
