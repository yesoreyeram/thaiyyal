import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type PartitionNodeData = {
  expression?: string;
  label?: string;
};

export function PartitionNode(props: NodePropsWithOptions<PartitionNodeData>) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const nodeInfo = getNodeInfo("partitionNode");
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
      )
    );
  };
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Partition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.expression ?? "item > 0")}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="item > 0"
        aria-label="Partition expression"
      />
      <Handle
        type="source"
        position={Position.Right}
        id="true"
        style={{ top: "30%" }}
        className="w-2 h-2 bg-green-500"
        title="Match path"
      />
      <Handle
        type="source"
        position={Position.Right}
        id="false"
        style={{ top: "70%" }}
        className="w-2 h-2 bg-red-500"
        title="No match path"
      />
    </NodeWrapper>
  );
}
