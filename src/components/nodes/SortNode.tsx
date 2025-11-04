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
export function SortNode({ id, data, onShowOptions }: NodePropsWithOptions<SortNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onFieldChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, field } } : n
      )
    );
  };

  const onOrderChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const order = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, order } } : n
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

  const nodeInfo = getNodeInfo("sortNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Sort")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-lime-600 to-lime-700 text-white shadow-lg rounded-lg border border-lime-500 hover:border-lime-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex items-center gap-0.5">
        <input
          value={String(data?.field ?? "")}
          type="text"
          onChange={onFieldChange}
          className="w-16 text-[10px] leading-tight border border-lime-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-lime-400 focus:outline-none"
          placeholder="Field"
          aria-label="Sort field"
        />
        <select
          value={String(data?.order ?? "asc")}
          onChange={onOrderChange}
          className="w-14 text-[10px] leading-tight border border-lime-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-lime-400 focus:outline-none"
        >
          <option value="asc">Asc</option>
          <option value="desc">Desc</option>
        </select>
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
