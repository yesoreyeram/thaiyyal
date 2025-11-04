/**
 * GroupByNode Component
 * 
 * Groups array elements by a specified key expression.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type GroupByNodeData = {
  key_field?: string;
  label?: string;
};

/**
 * GroupByNode React Component
 * 
 * This component renders a visual node in the workflow editor that groups array elements by a key
 * 
 * @param {NodePropsWithOptions<GroupByNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {GroupByNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function GroupByNode({ id, data, onShowOptions }: NodePropsWithOptions<GroupByNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const key_field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, key_field } } : n
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

  const nodeInfo = getNodeInfo("groupByNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Group By")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-violet-600 to-violet-700 text-white shadow-lg rounded-lg border border-violet-500 hover:border-violet-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.key_field ?? "category")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-violet-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-violet-400 focus:outline-none"
        placeholder="category"
        aria-label="Group by field"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
