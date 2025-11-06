/**
 * WhileLoopNode Component
 *
 * Continues executing while a condition is true.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type WhileLoopNodeData = {
  condition?: string;
  max_iterations?: number;
  label?: string;
};

/**
 * WhileLoopNode React Component
 *
 * This component renders a visual node in the workflow editor that repeats execution while a condition is true
 *
 * @param {NodePropsWithOptions<WhileLoopNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {WhileLoopNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function WhileLoopNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<WhileLoopNodeData>) {
  const { setNodes } = useReactFlow();

  const onConditionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const condition = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, condition } } : n
      )
    );
  };

  const onMaxIterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_iterations = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_iterations } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("whileLoopNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "While Loop")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="flex flex-col gap-0.5">
        <input
          value={String(data?.condition ?? ">0")}
          type="text"
          onChange={onConditionChange}
          className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Condition"
          aria-label="Loop condition"
        />
        <input
          value={Number(data?.max_iterations ?? 100)}
          type="number"
          onChange={onMaxIterChange}
          className="w-36 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Max iter"
          aria-label="Max iterations"
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

// ===== STATE & MEMORY NODES =====
