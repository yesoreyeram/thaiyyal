import { Handle, Position } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TransposeNodeData = {
  label?: string;
};

export function TransposeNode(props: NodePropsWithOptions<TransposeNodeData>) {
  const { id, data, onShowOptions } = props;
  const nodeInfo = getNodeInfo("transposeNode");
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Transpose")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs py-1 w-36">Transpose matrix</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
