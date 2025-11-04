/**
 * ExtractNode Component
 * 
 * Extracts a field from an object or array.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type ExtractNodeData = {
  field?: string;
  fields?: string[];
  label?: string;
};

/**
 * ExtractNode React Component
 * 
 * This component renders a visual node in the workflow editor that extracts fields from data
 * 
 * @param {NodePropsWithOptions<ExtractNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ExtractNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ExtractNode({ id, data }: NodeProps<ExtractNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, field } } : n))
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
        {String(data?.label || "Extract")}
      </div>
      <input
        value={String(data?.field ?? "")}
        type="text"
        onChange={onChange}
        className="w-28 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Field name"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </div>
  );
}
