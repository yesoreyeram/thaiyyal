// This file contains all the remaining node implementations
import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";

// ===== CONTROL FLOW NODES =====

type ForEachNodeData = {
  max_iterations?: number;
  label?: string;
};

export function ForEachNode({ id, data }: NodeProps<ForEachNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_iterations = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_iterations } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-amber-700 to-amber-800 text-white shadow-lg rounded-lg border border-amber-600 hover:border-amber-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "For Each")}</div>
      <input
        value={Number(data?.max_iterations ?? 1000)}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-amber-400 focus:outline-none"
        placeholder="Max iterations"
        aria-label="Max iterations"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type WhileLoopNodeData = {
  condition?: string;
  max_iterations?: number;
  label?: string;
};

export function WhileLoopNode({ id, data }: NodeProps<WhileLoopNodeData>) {
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

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-amber-600 to-amber-700 text-white shadow-lg rounded-lg border border-amber-500 hover:border-amber-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "While Loop")}</div>
      <input
        value={String(data?.condition ?? ">0")}
        type="text"
        onChange={onConditionChange}
        className="w-24 text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-amber-400 focus:outline-none"
        placeholder="Condition"
        aria-label="Loop condition"
      />
      <input
        value={Number(data?.max_iterations ?? 1000)}
        type="number"
        onChange={onMaxIterChange}
        className="mt-1 w-24 text-xs border border-amber-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-amber-400 focus:outline-none"
        placeholder="Max iterations"
        aria-label="Max iterations"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

// ===== STATE & MEMORY NODES =====

type VariableNodeData = {
  var_name?: string;
  var_op?: string;
  label?: string;
};

export function VariableNode({ id, data }: NodeProps<VariableNodeData>) {
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
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, var_op } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-sky-600 to-sky-700 text-white shadow-lg rounded-lg border border-sky-500 hover:border-sky-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Variable")}</div>
      <input
        value={String(data?.var_name ?? "")}
        type="text"
        onChange={onNameChange}
        className="w-24 text-xs border border-sky-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-sky-400 focus:outline-none"
        placeholder="Variable name"
      />
      <select
        value={String(data?.var_op ?? "get")}
        onChange={onOpChange}
        className="w-24 text-xs border border-sky-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-sky-400 focus:outline-none"
      >
        <option value="get">Get</option>
        <option value="set">Set</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type ExtractNodeData = {
  field?: string;
  fields?: string[];
  label?: string;
};

export function ExtractNode({ id, data }: NodeProps<ExtractNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, field } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-sky-700 to-sky-800 text-white shadow-lg rounded-lg border border-sky-600 hover:border-sky-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Extract")}</div>
      <input
        value={String(data?.field ?? "")}
        type="text"
        onChange={onChange}
        className="w-28 text-xs border border-sky-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-sky-400 focus:outline-none"
        placeholder="Field name"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type TransformNodeData = {
  transform_type?: string;
  label?: string;
};

export function TransformNode({ id, data }: NodeProps<TransformNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const transform_type = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, transform_type } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-sky-800 to-sky-900 text-white shadow-lg rounded-lg border border-sky-700 hover:border-sky-600 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Transform")}</div>
      <select
        value={String(data?.transform_type ?? "to_array")}
        onChange={onChange}
        className="w-28 text-xs border border-sky-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-sky-400 focus:outline-none"
      >
        <option value="to_array">To Array</option>
        <option value="to_object">To Object</option>
        <option value="flatten">Flatten</option>
        <option value="keys">Keys</option>
        <option value="values">Values</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type AccumulatorNodeData = {
  accum_op?: string;
  initial_value?: unknown;
  label?: string;
};

export function AccumulatorNode({ id, data }: NodeProps<AccumulatorNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const accum_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, accum_op } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-indigo-600 to-indigo-700 text-white shadow-lg rounded-lg border border-indigo-500 hover:border-indigo-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Accumulator")}</div>
      <select
        value={String(data?.accum_op ?? "sum")}
        onChange={onChange}
        className="w-24 text-xs border border-indigo-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-indigo-400 focus:outline-none"
      >
        <option value="sum">Sum</option>
        <option value="product">Product</option>
        <option value="concat">Concat</option>
        <option value="array">Array</option>
        <option value="count">Count</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type CounterNodeData = {
  counter_op?: string;
  delta?: number;
  initial_value?: number;
  label?: string;
};

export function CounterNode({ id, data }: NodeProps<CounterNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onOpChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const counter_op = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, counter_op } } : n
      )
    );
  };

  const onDeltaChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const delta = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, delta } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-indigo-700 to-indigo-800 text-white shadow-lg rounded-lg border border-indigo-600 hover:border-indigo-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Counter")}</div>
      <select
        value={String(data?.counter_op ?? "increment")}
        onChange={onOpChange}
        className="w-24 text-xs border border-indigo-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
      >
        <option value="increment">Increment</option>
        <option value="decrement">Decrement</option>
        <option value="reset">Reset</option>
        <option value="get">Get</option>
      </select>
      {(data?.counter_op === "increment" || data?.counter_op === "decrement") && (
        <input
          value={Number(data?.delta ?? 1)}
          type="number"
          onChange={onDeltaChange}
          className="w-24 text-xs border border-indigo-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
          placeholder="Delta"
        />
      )}
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

// ===== ADVANCED CONTROL FLOW NODES =====

type SwitchNodeData = {
  cases?: Array<{ when: string; value?: unknown; output_path?: string }>;
  default_path?: string;
  label?: string;
};

export function SwitchNode({ id, data }: NodeProps<SwitchNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onDefaultPathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const default_path = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, default_path } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-orange-600 to-orange-700 text-white shadow-lg rounded-lg border border-orange-500 hover:border-orange-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Switch")}</div>
      <input
        value={String(data?.default_path ?? "default")}
        type="text"
        onChange={onDefaultPathChange}
        className="w-24 text-xs border border-orange-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-orange-400 focus:outline-none"
        placeholder="Default path"
      />
      <div className="text-xs mt-1">Cases: {data?.cases?.length ?? 0}</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type ParallelNodeData = {
  max_concurrency?: number;
  label?: string;
};

export function ParallelNode({ id, data }: NodeProps<ParallelNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_concurrency = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_concurrency } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-orange-700 to-orange-800 text-white shadow-lg rounded-lg border border-orange-600 hover:border-orange-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Parallel")}</div>
      <input
        value={Number(data?.max_concurrency ?? 10)}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-orange-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-orange-500 focus:outline-none"
        placeholder="Max concurrency"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type JoinNodeData = {
  join_strategy?: string;
  timeout?: string;
  label?: string;
};

export function JoinNode({ id, data }: NodeProps<JoinNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const join_strategy = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, join_strategy } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-orange-800 to-orange-900 text-white shadow-lg rounded-lg border border-orange-700 hover:border-orange-600 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Join")}</div>
      <select
        value={String(data?.join_strategy ?? "all")}
        onChange={onChange}
        className="w-24 text-xs border border-orange-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-orange-600 focus:outline-none"
      >
        <option value="all">All</option>
        <option value="any">Any</option>
        <option value="first">First</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type SplitNodeData = {
  paths?: string[];
  label?: string;
};

export function SplitNode({ data }: NodeProps<SplitNodeData>) {
  return (
    <div className="px-3 py-2 bg-gradient-to-br from-pink-600 to-pink-700 text-white shadow-lg rounded-lg border border-pink-500 hover:border-pink-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Split")}</div>
      <div className="text-xs mt-1">Paths: {data?.paths?.length ?? 2}</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type DelayNodeData = {
  duration?: string;
  label?: string;
};

export function DelayNode({ id, data }: NodeProps<DelayNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const duration = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, duration } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-pink-700 to-pink-800 text-white shadow-lg rounded-lg border border-pink-600 hover:border-pink-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Delay")}</div>
      <input
        value={String(data?.duration ?? "1s")}
        type="text"
        onChange={onChange}
        className="w-24 text-xs border border-pink-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-500 focus:outline-none"
        placeholder="1s, 100ms..."
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type CacheNodeData = {
  cache_op?: string;
  cache_key?: string;
  ttl?: string;
  label?: string;
};

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
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, ttl } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-pink-800 to-pink-900 text-white shadow-lg rounded-lg border border-pink-700 hover:border-pink-600 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Cache")}</div>
      <select
        value={String(data?.cache_op ?? "get")}
        onChange={onOpChange}
        className="w-24 text-xs border border-pink-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-600 focus:outline-none"
      >
        <option value="get">Get</option>
        <option value="set">Set</option>
        <option value="delete">Delete</option>
      </select>
      <input
        value={String(data?.cache_key ?? "")}
        type="text"
        onChange={onKeyChange}
        className="w-24 text-xs border border-pink-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-600 focus:outline-none"
        placeholder="Cache key"
      />
      {data?.cache_op === "set" && (
        <input
          value={String(data?.ttl ?? "5m")}
          type="text"
          onChange={onTTLChange}
          className="w-24 text-xs border border-pink-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-600 focus:outline-none"
          placeholder="TTL (5m, 1h)"
        />
      )}
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

// ===== ERROR HANDLING & RESILIENCE NODES =====

type RetryNodeData = {
  max_attempts?: number;
  backoff_strategy?: string;
  initial_delay?: string;
  max_delay?: string;
  multiplier?: number;
  label?: string;
};

export function RetryNode({ id, data }: NodeProps<RetryNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onAttemptsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max_attempts = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, max_attempts } } : n
      )
    );
  };

  const onStrategyChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const backoff_strategy = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, backoff_strategy } } : n
      )
    );
  };

  const onInitialDelayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const initial_delay = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, initial_delay } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-red-600 to-red-700 text-white shadow-lg rounded-lg border border-red-500 hover:border-red-400 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Retry")}</div>
      <input
        value={Number(data?.max_attempts ?? 3)}
        type="number"
        onChange={onAttemptsChange}
        className="w-24 text-xs border border-red-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-red-400 focus:outline-none"
        placeholder="Max attempts"
      />
      <select
        value={String(data?.backoff_strategy ?? "exponential")}
        onChange={onStrategyChange}
        className="w-24 text-xs border border-red-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-red-400 focus:outline-none"
      >
        <option value="exponential">Exponential</option>
        <option value="linear">Linear</option>
        <option value="constant">Constant</option>
      </select>
      <input
        value={String(data?.initial_delay ?? "1s")}
        type="text"
        onChange={onInitialDelayChange}
        className="w-24 text-xs border border-red-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-red-400 focus:outline-none"
        placeholder="Initial delay"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type TryCatchNodeData = {
  fallback_value?: unknown;
  continue_on_error?: boolean;
  error_output_path?: string;
  label?: string;
};

export function TryCatchNode({ id, data }: NodeProps<TryCatchNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onContinueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const continue_on_error = e.target.checked;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, continue_on_error } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-red-700 to-red-800 text-white shadow-lg rounded-lg border border-red-600 hover:border-red-500 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Try-Catch")}</div>
      <label className="flex items-center gap-1 mt-1">
        <input
          type="checkbox"
          checked={data?.continue_on_error ?? true}
          onChange={onContinueChange}
          className="text-sm"
        />
        <span className="text-xs">Continue on error</span>
      </label>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}

type TimeoutNodeData = {
  timeout?: string;
  timeout_action?: string;
  label?: string;
};

export function TimeoutNode({ id, data }: NodeProps<TimeoutNodeData>) {
  const { setNodes } = useReactFlow();
  
  const onTimeoutChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const timeout = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, timeout } } : n
      )
    );
  };

  const onActionChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const timeout_action = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, timeout_action } } : n
      )
    );
  };

  return (
    <div className="px-3 py-2 bg-gradient-to-br from-red-800 to-red-900 text-white shadow-lg rounded-lg border border-red-700 hover:border-red-600 transition-all">
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs font-semibold mb-1 text-gray-200">{String(data?.label || "Timeout")}</div>
      <input
        value={String(data?.timeout ?? "30s")}
        type="text"
        onChange={onTimeoutChange}
        className="w-24 text-xs border border-red-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-red-600 focus:outline-none"
        placeholder="30s, 5m..."
      />
      <select
        value={String(data?.timeout_action ?? "error")}
        onChange={onActionChange}
        className="w-24 text-xs border border-red-700 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-red-600 focus:outline-none"
      >
        <option value="error">Error</option>
        <option value="continue_with_partial">Continue with partial</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </div>
  );
}
