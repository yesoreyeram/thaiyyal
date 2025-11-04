/**
 * ParseNode Component
 * 
 * Parses JSON or other formatted data.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type ParseNodeData = {
  input_type?: string;
  label?: string;
};

/**
 * ParseNode React Component
 * 
 * This component renders a visual node in the workflow editor that parses formatted data
 * 
 * @param {NodePropsWithOptions<ParseNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ParseNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ParseNode({ id, data }: NodeProps<ParseNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const input_type = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, input_type } } : n
      )
    );
  };

  return (
    <div className="px-2 py-1 bg-gradient-to-br from-purple-700 to-purple-800 text-white shadow-lg rounded border border-purple-600 hover:border-purple-500 transition-all">
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs font-semibold mb-1 text-purple-100">
        {String(data?.label || "Parse")}
      </div>
      <select
        value={String(data?.input_type ?? "AUTO")}
        onChange={onChange}
        className="w-24 text-xs border border-purple-600 px-2 py-1 rounded bg-purple-900 text-white focus:ring-2 focus:ring-purple-400 focus:outline-none"
      >
        <option value="AUTO">Auto</option>
        <option value="JSON">JSON</option>
        <option value="CSV">CSV</option>
        <option value="TSV">TSV</option>
        <option value="YAML">YAML</option>
        <option value="XML">XML</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </div>
  );
}
