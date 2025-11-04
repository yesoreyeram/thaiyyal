/**
 * MapNode Component
 * 
 * Transforms each element in an array using an expression. Similar to Array.prototype.map().
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type MapNodeData = {
  expression?: string;
  label?: string;
};

/**
 * MapNode React Component
 * 
 * This component renders a visual node in the workflow editor that transforms each element in an array using an expression
 * 
 * @param {NodePropsWithOptions<MapNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {MapNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function MapNode({ id, data, onShowOptions }: NodePropsWithOptions<MapNodeData>) {
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

  const nodeInfo = getNodeInfo("mapNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Map")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-cyan-600 to-cyan-700 text-white shadow-lg rounded-lg border border-cyan-500 hover:border-cyan-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item * 2")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-cyan-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-cyan-400 focus:outline-none"
        placeholder="item * 2"
        aria-label="Map expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
