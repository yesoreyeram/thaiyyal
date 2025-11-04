/**
 * CacheNode Component
 * 
 * Caches data for reuse.
 */

import React from "react";
import { Handle, Position, useReactFlow, NodeProps } from "reactflow";

type CacheNodeData = {
  cache_op?: string;
  cache_key?: string;
  ttl?: string;
  label?: string;
};

/**
 * CacheNode React Component
 * 
 * This component renders a visual node in the workflow editor that caches data for reuse
 * 
 * @param {NodePropsWithOptions<CacheNodeData>} props - Component props
 * @param {string} props.id - Unique identifier for this node instance
 * @param {CacheNodeData} props.data - Node configuration data
 * @param {function} [props.onShowOptions] - Callback to show the options context menu
 * @returns {JSX.Element} A rendered node component
 */
export function CacheNode({ id, data }: NodeProps<CacheNodeData>) {
  const { setNodes } = useReactFlow();

  const onOpChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const cache_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, cache_op } } : n
      )
    );
  };

  const onKeyChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const cache_key = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, cache_key } } : n
      )
    );
  };

  const onTTLChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const ttl = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, ttl } } : n))
    );
  };

  return (
    <div className="px-2 py-1 bg-gray-800 text-white shadow-lg rounded border border-gray-700 hover:border-gray-600 transition-all">
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <div className="text-xs font-semibold mb-1 text-gray-200">
        {String(data?.label || "Cache")}
      </div>
      <select
        value={String(data?.cache_op ?? "get")}
        onChange={onOpChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
      >
        <option value="get">Get</option>
        <option value="set">Set</option>
        <option value="delete">Delete</option>
      </select>
      <input
        value={String(data?.cache_key ?? "")}
        type="text"
        onChange={onKeyChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Cache key"
      />
      {data?.cache_op === "set" && (
        <input
          value={String(data?.ttl ?? "5m")}
          type="text"
          onChange={onTTLChange}
          className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
          placeholder="TTL (5m, 1h)"
        />
      )}
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </div>
  );
}

// ===== ERROR HANDLING & RESILIENCE NODES =====
