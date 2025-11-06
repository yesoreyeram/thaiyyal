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

export function CounterNode(props: NodePropsWithOptions<CounterNodeData>) {
  const { id, data, onShowOptions } = props;
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
  const nodeInfo = getNodeInfo("counterNode");
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Counter")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex items-center gap-0.5 w-36">
        <select
          value={String(data?.counter_op ?? "increment")}
          onChange={onOpChange}
          className="w-22 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
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
            className="w-14 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
            placeholder="Delta"
          />
        )}
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}

// ===== ADVANCED CONTROL FLOW NODES =====
