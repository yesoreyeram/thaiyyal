/**
 * DelayNode Component
 * 
 * Pauses execution for a specified duration.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type DelayNodeData = {
  duration?: string;
  label?: string;
};

/**
 * DelayNode React Component
 * 
 * This component renders a visual node in the workflow editor that pauses execution
 * 
 * @param {NodePropsWithOptions<DelayNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {DelayNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function DelayNode({ id, data }: NodeProps<DelayNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const duration = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, duration } } : n
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
        {String(data?.label || "Delay")}
      </div>
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
    </div>
  );
}
