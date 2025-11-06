/**
 * SortNode Component
 *
 * Sorts array elements based on a field or comparison expression.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SortNodeData = {
  field?: string;
  order?: string;
  label?: string;
};

/**
 * SortNode React Component
 *
 * This component renders a visual node in the workflow editor that sorts array elements
 *
 * @param {NodePropsWithOptions<SortNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {SortNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function SortNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SortNodeData>) {
  const { setNodes } = useReactFlow();

  const onFieldChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, field } } : n))
    );
  };

  const onOrderChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const order = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, order } } : n))
    );
  };

  const nodeInfo = getNodeInfo("sortNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Sort")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex items-center gap-0.5 w-36">
        <input
          value={String(data?.field ?? "")}
          type="text"
          onChange={onFieldChange}
          className="w-16 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Field"
          aria-label="Sort field"
        />
        <select
          value={String(data?.order ?? "asc")}
          onChange={onOrderChange}
          className="w-16 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        >
          <option value="asc">Asc</option>
          <option value="desc">Desc</option>
        </select>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
