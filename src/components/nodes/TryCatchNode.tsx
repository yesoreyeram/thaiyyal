/**
 * TryCatchNode Component
 * 
 * Handles errors gracefully.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type TryCatchNodeData = {
  fallback_value?: unknown;
  continue_on_error?: boolean;
  error_output_path?: string;
  label?: string;
};

/**
 * TryCatchNode React Component
 * 
 * This component renders a visual node in the workflow editor that handles errors gracefully
 * 
 * @param {NodePropsWithOptions<TryCatchNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {TryCatchNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function TryCatchNode({ id, data }: NodeProps<TryCatchNodeData>) {
  const { setNodes } = useReactFlow();

  const onContinueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const continue_on_error = e.target.checked;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, continue_on_error } } : n
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
        {String(data?.label || "Try-Catch")}
      </div>
      <label className="flex items-center gap-1 mt-1">
        <input
          type="checkbox"
          checked={data?.continue_on_error ?? true}
          onChange={onContinueChange}
          className="text-sm"
        />
        <span className="text-xs">Continue on error</span>
      </label>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </div>
  );
}
