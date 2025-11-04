/**
 * RangeNode Component
 * 
 * Generates a sequence of numbers from start to end.
 */

import React from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type RangeNodeData = {
  start?: number;
  end?: number;
  step?: number;
  label?: string;
};

/**
 * RangeNode React Component
 * 
 * This component renders a visual node in the workflow editor that generates a sequence of numbers
 * 
 * @param {NodePropsWithOptions<RangeNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {RangeNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function RangeNode({ id, data, onShowOptions }: NodePropsWithOptions<RangeNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onStartChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const start = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, start } } : n
      )
    );
  };

  const onEndChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const end = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, end } } : n
      )
    );
  };

  const onStepChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const step = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, step } } : n
      )
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("rangeNode");
  

  return (
    <NodeWrapper
      title={String(data?.label || "Range")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-green-600 to-green-700 text-white shadow-lg rounded-lg border border-green-500 hover:border-green-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex flex-col gap-0.5">
        <div className="flex items-center gap-0.5">
          <input
            value={Number(data?.start ?? 0)}
            type="number"
            onChange={onStartChange}
            className="w-12 text-[10px] leading-tight border border-green-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-green-400 focus:outline-none"
            placeholder="Start"
            aria-label="Start value"
          />
          <span className="text-[10px]">to</span>
          <input
            value={Number(data?.end ?? 10)}
            type="number"
            onChange={onEndChange}
            className="w-12 text-[10px] leading-tight border border-green-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-green-400 focus:outline-none"
            placeholder="End"
            aria-label="End value"
          />
        </div>
        <input
          value={Number(data?.step ?? 1)}
          type="number"
          onChange={onStepChange}
          className="w-full text-[10px] leading-tight border border-green-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-green-400 focus:outline-none"
          placeholder="Step: 1"
          aria-label="Step value"
        />
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
