/**
 * ExtractNode Component
 *
 * Extracts a field from an object or array.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

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
export function ExtractNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<ExtractNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, field } } : n))
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("extractNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Extract")}
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
    </NodeWrapper>
  );
}
