/**
 * ExpressionNode Component
 *
 * Applies a user-provided expression to transform input data.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ExpressionNodeData = {
  expression?: string;
  label?: string;
};

/**
 * ExpressionNode React Component
 *
 * This component renders a visual node in the workflow editor that applies an expression to input data
 *
 * @param {NodePropsWithOptions<ExpressionNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {ExpressionNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function ExpressionNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<ExpressionNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("expression");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Expression")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={String(data?.expression ?? "input * 2")}
        type="text"
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        placeholder="input * 2"
        aria-label="Expression"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
