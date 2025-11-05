/**
 * TimeoutNode Component
 *
 * Limits execution time.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TimeoutNodeData = {
  timeout?: string;
  timeout_action?: string;
  label?: string;
};

/**
 * TimeoutNode React Component
 *
 * This component renders a visual node in the workflow editor that limits execution time
 *
 * @param {NodePropsWithOptions<TimeoutNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {TimeoutNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function TimeoutNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<TimeoutNodeData>) {
  const { setNodes } = useReactFlow();

  const onTimeoutChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const timeout = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, timeout } } : n))
    );
  };

  const onActionChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const timeout_action = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, timeout_action } } : n
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

  const nodeInfo = getNodeInfo("timeoutNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Timeout")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.timeout ?? "30s")}
        type="text"
        onChange={onTimeoutChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="30s, 5m..."
      />
      <select
        value={String(data?.timeout_action ?? "error")}
        onChange={onActionChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
      >
        <option value="error">Error</option>
        <option value="continue_with_partial">Continue with partial</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
