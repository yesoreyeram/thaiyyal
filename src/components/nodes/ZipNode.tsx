/**
 * ZipNode Component
 * 
 * Combines multiple arrays element-wise.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ZipNodeData = {
  label?: string;
};

/**
 * ZipNode React Component
 * 
 * This component renders a visual node in the workflow editor that combines multiple arrays element-wise
 * 
 * @param {NodePropsWithOptions<ZipNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ZipNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ZipNode({ id, data, onShowOptions }: NodePropsWithOptions<ZipNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("zipNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Zip")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-yellow-600 to-yellow-700 text-white shadow-lg rounded-lg border border-yellow-500 hover:border-yellow-400 transition-all"
    >
      <Handle 
        type="target" 
        position={Position.Left} 
        id="array1"
        style={{ top: '30%' }}
        className="w-2 h-2 bg-blue-400"
        title="First array"
      />
      <Handle 
        type="target" 
        position={Position.Left} 
        id="array2"
        style={{ top: '70%' }}
        className="w-2 h-2 bg-blue-400"
        title="Second array"
      />
      <div className="text-xs py-1">Combine arrays</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
