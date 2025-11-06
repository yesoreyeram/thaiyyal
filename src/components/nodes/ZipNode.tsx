import { Handle, Position } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ZipNodeData = {
  label?: string;
};

export function ZipNode(props: NodePropsWithOptions<ZipNodeData>) {
  const { id, data, onShowOptions } = props;
  const nodeInfo = getNodeInfo("zipNode");
  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Zip")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        id="array1"
        style={{ top: "30%" }}
        className="w-2 h-2 bg-blue-400"
        title="First array"
      />
      <Handle
        type="target"
        position={Position.Left}
        id="array2"
        style={{ top: "70%" }}
        className="w-2 h-2 bg-blue-400"
        title="Second array"
      />
      <div className="text-xs py-1 w-36">Combine arrays</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
