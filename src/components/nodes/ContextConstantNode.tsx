import { NodeProps, useReactFlow } from "reactflow";
import React from "react";

type ContextConstantNodeData = {
  context_name?: string;
  context_value?: string | number;
  label?: string;
};

export function ContextConstantNode({ id, data }: NodeProps<ContextConstantNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const context_name = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, context_name } } : n
      )
    );
  };

  const onValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const context_value = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, context_value } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-amber-700 to-amber-800 text-white shadow-lg rounded-lg border-2 border-amber-400 hover:border-amber-300 transition-all">
      <div className="flex items-center gap-1 mb-1">
        <span className="text-sm">🔒</span>
        <div className="text-xs font-semibold text-amber-200">
          {String(data?.label || "Constant")}
        </div>
      </div>
      <div className="space-y-1">
        <input
          value={String(data?.context_name ?? "")}
          type="text"
          onChange={onNameChange}
          className="w-full text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-amber-400 focus:outline-none"
          placeholder="Name"
          aria-label="Constant name"
        />
        <input
          value={String(data?.context_value ?? "")}
          type="text"
          onChange={onValueChange}
          className="w-full text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-amber-400 focus:outline-none"
          placeholder="Value"
          aria-label="Constant value"
        />
      </div>
      <div className="mt-1 text-xs text-amber-300 font-mono">
        {'{{ const.'}
        {String(data?.context_name || '...')}
        {' }}'}
      </div>
    </div>
  );
}
