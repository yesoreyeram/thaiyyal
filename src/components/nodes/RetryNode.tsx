/**
 * RetryNode Component
 *
 * Retries failed operations.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type RetryNodeData = {
  max_attempts?: number;
  backoff_strategy?: string;
  initial_delay?: string;
  max_delay?: string;
  multiplier?: number;
  label?: string;
};

/**
 * RetryNode React Component
 *
 * This component renders a visual node in the workflow editor that retries failed operations
 *
 * @param {NodePropsWithOptions<RetryNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {RetryNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function RetryNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<RetryNodeData>) {
  const { setNodes } = useReactFlow();

  const onAttemptsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_attempts = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_attempts } } : n
      )
    );
  };

  const onStrategyChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const backoff_strategy = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, backoff_strategy } } : n
      )
    );
  };

  const onInitialDelayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const initial_delay = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, initial_delay } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("retryNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Retry")}
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
          value={String(data?.backoff_strategy ?? "exponential")}
          onChange={onStrategyChange}
          className="w-22 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        >
          <option value="exponential">Exponential</option>
          <option value="linear">Linear</option>
          <option value="constant">Constant</option>
        </select>
        <input
          value={Number(data?.max_attempts ?? 3)}
          type="number"
          onChange={onAttemptsChange}
          className="w-14 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Max attempts"
        />
      </div>
      <div className="flex items-center gap-0.5 w-36 mt-1">
        <input
          value={String(data?.initial_delay ?? "1s")}
          type="text"
          onChange={onInitialDelayChange}
          className="w-36 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
          placeholder="Initial delay"
        />
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
