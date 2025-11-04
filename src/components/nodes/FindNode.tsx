/**
 * FindNode Component
 * 
 * Returns the first element matching a condition.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type FindNodeData = {
  expression?: string;
  label?: string;
};

/**
 * FindNode React Component
 * 
 * This component renders a visual node in the workflow editor that finds the first element matching a condition
 * 
 * @param {NodePropsWithOptions<FindNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {FindNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function FindNode({ id, data, onShowOptions }: NodePropsWithOptions<FindNodeData>) {
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

  const nodeInfo = getNodeInfo("findNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Find")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-sky-600 to-sky-700 text-white shadow-lg rounded-lg border border-sky-500 hover:border-sky-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item.id == 1")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-sky-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-sky-400 focus:outline-none"
        placeholder="item.id == 1"
        aria-label="Find expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
