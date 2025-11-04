/**
 * ReduceNode Component
 * 
 * Reduces an array to a single value using an accumulator expression. Similar to Array.prototype.reduce().
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ReduceNodeData = {
  expression?: string;
  initial_value?: string;
  label?: string;
};

/**
 * ReduceNode React Component
 * 
 * This component renders a visual node in the workflow editor that reduces an array to a single value using an accumulator
 * 
 * @param {NodePropsWithOptions<ReduceNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ReduceNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ReduceNode({ id, data, onShowOptions }: NodePropsWithOptions<ReduceNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onExpressionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
      )
    );
  };

  const onInitialValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const initial_value = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, initial_value } } : n
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

  const nodeInfo = getNodeInfo("reduceNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Reduce")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-teal-600 to-teal-700 text-white shadow-lg rounded-lg border border-teal-500 hover:border-teal-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex flex-col gap-0.5">
        <input
          value={String(data?.expression ?? "acc + item")}
          type="text"
          onChange={onExpressionChange}
          className="w-32 text-[10px] leading-tight border border-teal-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-teal-400 focus:outline-none"
          placeholder="acc + item"
          aria-label="Reduce expression"
        />
        <input
          value={String(data?.initial_value ?? "0")}
          type="text"
          onChange={onInitialValueChange}
          className="w-32 text-[10px] leading-tight border border-teal-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-teal-400 focus:outline-none"
          placeholder="Initial: 0"
          aria-label="Initial value"
        />
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
