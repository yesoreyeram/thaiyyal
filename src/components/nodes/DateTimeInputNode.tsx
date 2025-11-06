import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type DateTimeInputNodeData = {
  datetime_value?: string;
  label?: string;
};

export function DateTimeInputNode(
  props: NodePropsWithOptions<DateTimeInputNodeData>
) {
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

  const nodeInfo = getNodeInfo("datetimeInputNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "DateTime")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.datetime_value ?? "")}
        type="datetime-local"
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        aria-label="DateTime value"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
