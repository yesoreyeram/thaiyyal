/**
 * SampleNode Component
 * 
 * Randomly samples elements from an array.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SampleNodeData = {
  count?: number;
  label?: string;
};

/**
 * SampleNode React Component
 * 
 * This component renders a visual node in the workflow editor that randomly samples array elements
 * 
 * @param {NodePropsWithOptions<SampleNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {SampleNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function SampleNode({ id, data, onShowOptions }: NodePropsWithOptions<SampleNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const count = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, count } } : n
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

  const nodeInfo = getNodeInfo("sampleNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Sample")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-blue-600 to-blue-700 text-white shadow-lg rounded-lg border border-blue-500 hover:border-blue-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={Number(data?.count ?? 1)}
        type="number"
        onChange={onChange}
        className="w-20 text-xs border border-blue-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Count"
        aria-label="Sample count"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
