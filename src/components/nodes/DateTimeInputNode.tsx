import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type Props = NodePropsWithOptions<{ datetime_value?: string }>;

const nodeInfo = getNodeInfo("datetimeInputNode");

export function DateTimeInputNode(props: Props) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const datetime_value = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, datetime_value } } : n
      )
    );
  };
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "DateTime")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2" />
      <input
        value={String(data?.datetime_value ?? "")}
        type="datetime-local"
        onChange={onChange}
        className="w-full text-xs border px-1.5 py-0.5 rounded focus:ring-1 focus:outline-none dark:scheme-dark border-gray-600 bg-gray-900 text-white focus:ring-blue-400"
        aria-label="DateTime value"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2" />
    </NodeWrapper>
  );
}
