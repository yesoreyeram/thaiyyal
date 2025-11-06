import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SwitchNodeData = {
  cases?: Array<{ when: string; value?: unknown; output_path?: string }>;
  default_path?: string;
  label?: string;
};

export function SwitchNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SwitchNodeData>) {
  const { setNodes } = useReactFlow();

  const onDefaultPathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const default_path = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, default_path } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("switchNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Switch")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.default_path ?? "default")}
        type="text"
        onChange={onDefaultPathChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="Default path"
      />
      <div className="text-xs mt-1">Cases: {data?.cases?.length ?? 0}</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
