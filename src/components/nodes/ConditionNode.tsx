import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

const nodeInfo = getNodeInfo("conditionNode");

type Props = NodePropsWithOptions<{
  condition?: string;
  true_path?: string;
  false_path?: string;
}>;

export function ConditionNode(props: Props) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
      )
    );
  };
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Condition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2" />
      <input
        value={String(data?.condition ?? ">0")}
        type="text"
        onChange={onConditionChange}
        className="w-full text-xs border px-1.5 py-0.5 rounded focus:ring-1 focus:outline-none dark:scheme-dark border-gray-600 bg-gray-900 text-white focus:ring-blue-400"
        placeholder="node.id.value > 100"
        aria-label="Condition expression"
      />
      <Handle
        type="source"
        position={Position.Right}
        id="true"
        style={{ top: "30%" }}
        className="w-2 h-2"
        title="True path"
      />
      <Handle
        type="source"
        position={Position.Right}
        id="false"
        style={{ top: "70%" }}
        className="w-2 h-2 bg-red-500"
        title="False path"
      />
    </NodeWrapper>
  );
}
