/**
 * CounterNode Component
 *
 * Increments or decrements a counter.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type CounterNodeData = {
  counter_op?: string;
  delta?: number;
  initial_value?: number;
  label?: string;
};

/**
 * CounterNode React Component
 *
 * This component renders a visual node in the workflow editor that manages a counter
 *
 * @param {NodePropsWithOptions<CounterNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {CounterNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function CounterNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<CounterNodeData>) {
  const { setNodes } = useReactFlow();

  const onOpChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const counter_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, counter_op } } : n
      )
    );
  };

  const onDeltaChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const delta = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, delta } } : n))
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("counterNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Counter")}
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
        value={String(data?.counter_op ?? "increment")}
        onChange={onOpChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
      >
        <option value="increment">Increment</option>
        <option value="decrement">Decrement</option>
        <option value="reset">Reset</option>
        <option value="get">Get</option>
      </select>
      {(data?.counter_op === "increment" ||
        data?.counter_op === "decrement") && (
        <input
          value={Number(data?.delta ?? 1)}
          type="number"
          onChange={onDeltaChange}
          className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
          placeholder="Delta"
        />
      )}
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}

// ===== ADVANCED CONTROL FLOW NODES =====
