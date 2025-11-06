import { Handle, Position } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ReverseNodeData = {
  label?: string;
};

export function ReverseNode(props: NodePropsWithOptions<ReverseNodeData>) {
  const { id, data, onShowOptions } = props;
  const nodeInfo = getNodeInfo("reverseNode");
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Reverse")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs py-1 w-36">Reverse array order</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
