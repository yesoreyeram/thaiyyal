/**
 * TimeoutNode Component
 * 
 * Limits execution time.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

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
export function TimeoutNode({ id, data }: NodeProps<TimeoutNodeData>) {
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

  return (
    <div className="px-2 py-1 bg-gray-800 text-white shadow-lg rounded border border-gray-700 hover:border-gray-600 transition-all">
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs font-semibold mb-1 text-gray-200">
        {String(data?.label || "Timeout")}
      </div>
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
    </div>
  );
}
