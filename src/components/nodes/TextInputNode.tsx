import React from "react";
import { TextArea } from "../ui/NodeTextArea";
import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TextInputNodeData = {
  text?: string;
  label?: string;
};

export function TextInputNode(props: NodePropsWithOptions<TextInputNodeData>) {
  const { id, data, onShowOptions } = props;
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, text } } : n))
    );
  };

  const nodeInfo = getNodeInfo("textInputNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Text Input")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <TextArea text={String(data?.text ?? "")} onChange={onChange} />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
