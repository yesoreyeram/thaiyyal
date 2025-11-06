import { Handle, Position } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SplitNodeData = {
  paths?: string[];
  label?: string;
};

export function SplitNode(props: NodePropsWithOptions<SplitNodeData>) {
  const { id, data, onShowOptions } = props;
  const nodeInfo = getNodeInfo("splitNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Split")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="w-36 text-xs mt-1">Paths: {data?.paths?.length ?? 2}</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
