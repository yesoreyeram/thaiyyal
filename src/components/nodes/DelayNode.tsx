/**
 * DelayNode Component
 *
 * Pauses execution for a specified duration.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type DelayNodeData = {
  duration?: string;
  label?: string;
};

export function DelayNode(props: NodePropsWithOptions<DelayNodeData>) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const duration = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, duration } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("delayNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Delay")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.duration ?? "1s")}
        type="text"
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="1s, 100ms..."
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
