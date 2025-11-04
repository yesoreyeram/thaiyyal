import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";


type ConditionNodeData = {
  condition?: string;
  true_path?: string;
  false_path?: string;
  label?: string;
};

export function ConditionNode({ id, data, onShowOptions }: NodePropsWithOptions<ConditionNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
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

  const nodeInfo = getNodeInfo("conditionNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Condition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-amber-600 to-amber-700 text-white shadow-lg rounded-lg border border-amber-500 hover:border-amber-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.condition ?? ">0")}
        type="text"
        onChange={onConditionChange}
        className="w-28 text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-amber-400 focus:outline-none"
        placeholder="node.id.value > 100"
        aria-label="Condition expression"
      />
      {/* True path handle (top right) */}
      <Handle 
        type="source" 
        position={Position.Right} 
        id="true"
        style={{ top: '30%' }}
        className="w-2 h-2 bg-green-500"
        title="True path"
      />
      {/* False path handle (bottom right) */}
      <Handle 
        type="source" 
        position={Position.Right} 
        id="false"
        style={{ top: '70%' }}
        className="w-2 h-2 bg-red-500"
        title="False path"
      />
    </NodeWrapper>
  );
}
