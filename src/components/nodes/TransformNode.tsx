/**
 * TransformNode Component
 *
 * Transforms data using an expression.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type TransformNodeData = {
  transform_type?: string;
  label?: string;
};

/**
 * TransformNode React Component
 *
 * This component renders a visual node in the workflow editor that transforms data using expressions
 *
 * @param {NodePropsWithOptions<TransformNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {TransformNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function TransformNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<TransformNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const transform_type = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, transform_type } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("transformNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Transform")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={String(data?.transform_type ?? "to_array")}
        onChange={onChange}
        className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="to_array">To Array</option>
        <option value="to_object">To Object</option>
        <option value="flatten">Flatten</option>
        <option value="keys">Keys</option>
        <option value="values">Values</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
