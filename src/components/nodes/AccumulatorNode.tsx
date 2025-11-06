import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type AccumulatorNodeData = {
  accum_op?: string;
  initial_value?: unknown;
  label?: string;
};

export function AccumulatorNode(
  props: NodePropsWithOptions<AccumulatorNodeData>
) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const accum_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, accum_op } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("accumulatorNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Accumulator")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={String(data?.accum_op ?? "sum")}
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="sum">Sum</option>
        <option value="product">Product</option>
        <option value="concat">Concat</option>
        <option value="array">Array</option>
        <option value="count">Count</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
