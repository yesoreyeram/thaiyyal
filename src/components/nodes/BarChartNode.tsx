import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type BarChartNodeData = {
  label?: string;
  bar_color?: string;
  bar_width?: string;
  show_values?: boolean;
  max_bars?: number;
  orientation?: "vertical" | "horizontal";
};

export function BarChartNode({ id, data, ...props }: NodeProps<BarChartNodeData>) {
  const { setNodes } = useReactFlow();

  const onColorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const bar_color = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, bar_color } } : n
      )
    );
  };

  const onWidthChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const bar_width = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, bar_width } } : n
      )
    );
  };

  const onShowValuesChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const show_values = e.target.checked;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, show_values } } : n
      )
    );
  };

  const onMaxBarsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_bars = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_bars } } : n
      )
    );
  };

  const onOrientationChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const orientation = e.target.value as "vertical" | "horizontal";
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, orientation } } : n
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

  const nodeInfo = getNodeInfo("barChartNode");
  // Type assertion is consistent with other nodes in the codebase (see AllNodes.tsx, ArrayNodes.tsx)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const onShowOptions = (props as any).onShowOptions;

  return (
    <NodeWrapper
      title={String(data?.label || "Bar Chart")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex flex-col gap-1">
        {/* Orientation selector */}
        <select
          value={String(data?.orientation ?? "vertical")}
          onChange={onOrientationChange}
          className="w-28 text-[10px] leading-tight border border-gray-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          aria-label="Chart orientation"
        >
          <option value="vertical">Vertical</option>
          <option value="horizontal">Horizontal</option>
        </select>

        {/* Bar width selector */}
        <select
          value={String(data?.bar_width ?? "medium")}
          onChange={onWidthChange}
          className="w-28 text-[10px] leading-tight border border-gray-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
          aria-label="Bar width"
        >
          <option value="thin">Thin</option>
          <option value="medium">Medium</option>
          <option value="thick">Thick</option>
        </select>

        {/* Bar color picker */}
        <div className="flex items-center gap-1">
          <input
            type="color"
            value={String(data?.bar_color ?? "#3b82f6")}
            onChange={onColorChange}
            className="w-8 h-5 border border-gray-600 rounded bg-gray-900 cursor-pointer"
            aria-label="Bar color"
          />
          <span className="text-[9px] text-gray-400">Color</span>
        </div>

        {/* Max bars input */}
        <input
          value={data?.max_bars ?? 20}
          type="number"
          onChange={onMaxBarsChange}
          className="w-28 text-[10px] leading-tight border border-gray-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-blue-400 focus:outline-none"
          placeholder="Max bars"
          min="1"
          max="100"
          aria-label="Maximum number of bars"
        />

        {/* Show values checkbox */}
        <label className="flex items-center gap-1">
          <input
            type="checkbox"
            checked={data?.show_values ?? true}
            onChange={onShowValuesChange}
            className="w-3 h-3 text-blue-600 bg-gray-900 border-gray-600 rounded focus:ring-blue-500"
            aria-label="Show values on bars"
          />
          <span className="text-[9px] text-gray-300">Show values</span>
        </label>
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
