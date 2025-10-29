import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React, { useState } from "react";

type HttpPaginationNodeData = {
  base_url?: string;
  max_pages?: number;
  total_items?: number;
  page_size?: number;
  start_page?: number;
  page_param?: string;
  break_on_error?: boolean;
  label?: string;
};

export function HttpPaginationNode({ id, data }: NodeProps<HttpPaginationNodeData>) {
  const { setNodes } = useReactFlow();
  const [configMode, setConfigMode] = useState<'max_pages' | 'total_items'>(
    data?.max_pages ? 'max_pages' : 'total_items'
  );

  const updateNodeData = (updates: Partial<HttpPaginationNodeData>) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, ...updates } } : n
      )
    );
  };

  const onBaseUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateNodeData({ base_url: e.target.value });
  };

  const onMaxPagesChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    updateNodeData({ max_pages: isNaN(value) ? undefined : value });
  };

  const onTotalItemsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    updateNodeData({ total_items: isNaN(value) ? undefined : value });
  };

  const onPageSizeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    updateNodeData({ page_size: isNaN(value) ? undefined : value });
  };

  const onStartPageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value);
    updateNodeData({ start_page: isNaN(value) ? undefined : value });
  };

  const onPageParamChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateNodeData({ page_param: e.target.value || undefined });
  };

  const onBreakOnErrorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateNodeData({ break_on_error: e.target.checked });
  };

  const onConfigModeChange = (mode: 'max_pages' | 'total_items') => {
    setConfigMode(mode);
    if (mode === 'max_pages') {
      updateNodeData({ total_items: undefined, page_size: undefined });
    } else {
      updateNodeData({ max_pages: undefined });
    }
  };

  return (
    <div className="p-3 bg-indigo-600 text-white shadow-lg rounded-lg border-2 border-indigo-400 min-w-[280px]">
      <Handle type="target" position={Position.Left} />
      
      <div className="text-sm font-bold mb-2">HTTP Pagination</div>
      
      {/* Base URL */}
      <div className="mb-2">
        <label className="text-xs block mb-1">Base URL:</label>
        <input
          value={data?.base_url ?? ""}
          type="text"
          onChange={onBaseUrlChange}
          className="w-full border px-2 py-1 rounded text-black text-xs"
          placeholder="https://api.example.com/items"
        />
      </div>

      {/* Configuration Mode Toggle */}
      <div className="mb-2">
        <div className="flex gap-1 text-xs mb-1">
          <button
            onClick={() => onConfigModeChange('max_pages')}
            className={`px-2 py-1 rounded ${
              configMode === 'max_pages' 
                ? 'bg-indigo-800' 
                : 'bg-indigo-700 opacity-60'
            }`}
          >
            Max Pages
          </button>
          <button
            onClick={() => onConfigModeChange('total_items')}
            className={`px-2 py-1 rounded ${
              configMode === 'total_items' 
                ? 'bg-indigo-800' 
                : 'bg-indigo-700 opacity-60'
            }`}
          >
            Total Items
          </button>
        </div>
      </div>

      {/* Conditional Fields */}
      {configMode === 'max_pages' ? (
        <div className="mb-2">
          <label className="text-xs block mb-1">Max Pages:</label>
          <input
            value={data?.max_pages ?? ""}
            type="number"
            min="1"
            onChange={onMaxPagesChange}
            className="w-full border px-2 py-1 rounded text-black text-xs"
            placeholder="5"
          />
        </div>
      ) : (
        <>
          <div className="mb-2">
            <label className="text-xs block mb-1">Total Items:</label>
            <input
              value={data?.total_items ?? ""}
              type="number"
              min="1"
              onChange={onTotalItemsChange}
              className="w-full border px-2 py-1 rounded text-black text-xs"
              placeholder="50"
            />
          </div>
          <div className="mb-2">
            <label className="text-xs block mb-1">Page Size:</label>
            <input
              value={data?.page_size ?? ""}
              type="number"
              min="1"
              onChange={onPageSizeChange}
              className="w-full border px-2 py-1 rounded text-black text-xs"
              placeholder="10"
            />
          </div>
        </>
      )}

      {/* Advanced Options */}
      <details className="mb-2">
        <summary className="text-xs cursor-pointer mb-1">Advanced Options</summary>
        <div className="mt-2 space-y-2">
          <div>
            <label className="text-xs block mb-1">Start Page:</label>
            <input
              value={data?.start_page ?? ""}
              type="number"
              min="1"
              onChange={onStartPageChange}
              className="w-full border px-2 py-1 rounded text-black text-xs"
              placeholder="1 (default)"
            />
          </div>
          <div>
            <label className="text-xs block mb-1">Page Param:</label>
            <input
              value={data?.page_param ?? ""}
              type="text"
              onChange={onPageParamChange}
              className="w-full border px-2 py-1 rounded text-black text-xs"
              placeholder="page (default)"
            />
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={data?.break_on_error ?? true}
              onChange={onBreakOnErrorChange}
              className="w-4 h-4"
            />
            <label className="text-xs">Break on Error</label>
          </div>
        </div>
      </details>

      <div className="text-xs text-indigo-200 mt-2">
        Use {'{page}'} in URL for placeholder
      </div>

      <Handle type="source" position={Position.Right} />
    </div>
  );
}
