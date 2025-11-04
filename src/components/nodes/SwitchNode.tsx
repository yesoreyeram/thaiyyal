/**
 * SwitchNode Component
 * 
 * Routes data based on case matching.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type SwitchNodeData = {
  cases?: Array<{ when: string; value?: unknown; output_path?: string }>;
  default_path?: string;
  label?: string;
};

/**
 * SwitchNode React Component
 * 
 * This component renders a visual node in the workflow editor that routes data based on cases
 * 
 * @param {NodePropsWithOptions<SwitchNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {SwitchNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function SwitchNode({ id, data }: NodeProps<SwitchNodeData>) {
  const { setNodes } = useReactFlow();

  const onDefaultPathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const default_path = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, default_path } } : n
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
        {String(data?.label || "Switch")}
      </div>
      <input
        value={String(data?.default_path ?? "default")}
        type="text"
        onChange={onDefaultPathChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Default path"
      />
      <div className="text-xs mt-1">Cases: {data?.cases?.length ?? 0}</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </div>
  );
}
