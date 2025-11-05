/**
 * AccumulatorNode Component
 *
 * Accumulates values across executions.
 */

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

/**
 * AccumulatorNode React Component
 *
 * This component renders a visual node in the workflow editor that accumulates values across executions
 *
 * @param {NodePropsWithOptions<AccumulatorNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {AccumulatorNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function AccumulatorNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<AccumulatorNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const accum_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, accum_op } } : n
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

  const nodeInfo = getNodeInfo("accumulatorNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Accumulator")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={String(data?.accum_op ?? "sum")}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
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
