import { ChangeEvent } from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type GroupByNodeData = {
  key_field?: string;
  label?: string;
};

export function GroupByNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<GroupByNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    const key_field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, key_field } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("groupByNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Group By")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.key_field ?? "category")}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="category"
        aria-label="Group by field"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
