/**
 * SplitNode Component
 *
 * Splits data into multiple outputs.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SplitNodeData = {
  paths?: string[];
  label?: string;
};

/**
 * SplitNode React Component
 *
 * This component renders a visual node in the workflow editor that splits data to multiple outputs
 *
 * @param {NodePropsWithOptions<SplitNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {SplitNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function SplitNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SplitNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("splitNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Split")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs mt-1">Paths: {data?.paths?.length ?? 2}</div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
