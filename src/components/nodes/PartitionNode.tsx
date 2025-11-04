/**
 * PartitionNode Component
 * 
 * Splits an array into two groups based on a condition.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type PartitionNodeData = {
  expression?: string;
  label?: string;
};

/**
 * PartitionNode React Component
 * 
 * This component renders a visual node in the workflow editor that splits an array based on a condition
 * 
 * @param {NodePropsWithOptions<PartitionNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {PartitionNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function PartitionNode({ id, data, onShowOptions }: NodePropsWithOptions<PartitionNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
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

  const nodeInfo = getNodeInfo("partitionNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Partition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-orange-600 to-orange-700 text-white shadow-lg rounded-lg border border-orange-500 hover:border-orange-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item > 0")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-orange-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-orange-400 focus:outline-none"
        placeholder="item > 0"
        aria-label="Partition expression"
      />
      <Handle 
        type="source" 
        position={Position.Right} 
        id="true"
        style={{ top: '30%' }}
        className="w-2 h-2 bg-green-500"
        title="Match path"
      />
      <Handle 
        type="source" 
        position={Position.Right} 
        id="false"
        style={{ top: '70%' }}
        className="w-2 h-2 bg-red-500"
        title="No match path"
      />
    </NodeWrapper>
  );
}
