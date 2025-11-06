/**
 * VariableNode Component
 *
 * Gets or sets a workflow variable.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type VariableNodeData = {
  var_name?: string;
  var_op?: string;
  label?: string;
};

/**
 * VariableNode React Component
 *
 * This component renders a visual node in the workflow editor that manages workflow variables
 *
 * @param {NodePropsWithOptions<VariableNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {VariableNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function VariableNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<VariableNodeData>) {
  const { setNodes } = useReactFlow();

  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const var_name = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, var_name } } : n
      )
    );
  };

  const onOpChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const var_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, var_op } } : n))
    );
  };

  const nodeInfo = getNodeInfo("variableNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Variable")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex items-center gap-0.5 w-36">
        <select
          value={String(data?.var_op ?? "get")}
          onChange={onOpChange}
          className="w-14 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
        >
          <option value="get">Get</option>
          <option value="set">Set</option>
        </select>
        <input
          value={String(data?.var_name ?? "")}
          type="text"
          onChange={onNameChange}
          className="w-22 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Name"
        />
      </div>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
