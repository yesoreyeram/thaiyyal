import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type BooleanInputNodeData = {
  boolean_value?: boolean;
  label?: string;
};

type Props = NodePropsWithOptions<BooleanInputNodeData>;

const nodeInfo = getNodeInfo("booleanInputNode");

export function BooleanInputNode(props: Props) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const boolean_value = e.target.checked;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, boolean_value } } : n
      )
    );
  };
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Boolean")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="w-36 flex items-center gap-2 px-1">
        <input
          checked={data?.boolean_value ?? false}
          type="checkbox"
          onChange={onChange}
          className="text-xs border px-1.5 py-0.5 rounded focus:ring-1 focus:outline-none dark:scheme-dark border-gray-600  bg-gray-900 text-white focus:ring-blue-400 "
          aria-label="Boolean value"
          style={{ colorScheme: "dark" }}
        />
        <span className="text-xs text-white font-extralight">
          {data?.boolean_value ? "True" : "False"}
        </span>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
