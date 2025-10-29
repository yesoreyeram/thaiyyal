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
    <div className="p-2 bg-yellow-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">For Each</div>
      <input
        value={data?.max_iterations ?? 1000}
        type="number"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Max iterations"
      />
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-yellow-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">While Loop</div>
      <input
        value={data?.condition ?? ">0"}
        type="text"
        onChange={onConditionChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Condition"
      />
      <input
        value={data?.max_iterations ?? 100}
        type="number"
        onChange={onMaxIterChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Max iterations"
      />
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-blue-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Variable</div>
      <input
        value={data?.var_name ?? ""}
        type="text"
        onChange={onNameChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Variable name"
      />
      <select
        value={data?.var_op ?? "get"}
        onChange={onOpChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="get">Get</option>
        <option value="set">Set</option>
      </select>
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-blue-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Extract</div>
      <input
        value={data?.field ?? ""}
        type="text"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Field name"
      />
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-blue-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Transform</div>
      <select
        value={data?.transform_type ?? "to_array"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="to_array">To Array</option>
        <option value="to_object">To Object</option>
        <option value="flatten">Flatten</option>
        <option value="keys">Keys</option>
        <option value="values">Values</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

type AccumulatorNodeData = {
  accum_op?: string;
  initial_value?: any;
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
    <div className="p-2 bg-indigo-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Accumulator</div>
      <select
        value={data?.accum_op ?? "sum"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="sum">Sum</option>
        <option value="product">Product</option>
        <option value="concat">Concat</option>
        <option value="array">Array</option>
        <option value="count">Count</option>
      </select>
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-indigo-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Counter</div>
      <select
        value={data?.counter_op ?? "increment"}
        onChange={onOpChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="increment">Increment</option>
        <option value="decrement">Decrement</option>
        <option value="reset">Reset</option>
        <option value="get">Get</option>
      </select>
      {(data?.counter_op === "increment" || data?.counter_op === "decrement") && (
        <input
          value={data?.delta ?? 1}
          type="number"
          onChange={onDeltaChange}
          className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
          placeholder="Delta"
        />
      )}
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

// ===== ADVANCED CONTROL FLOW NODES =====

type SwitchNodeData = {
  cases?: Array<{ when: string; value?: any; output_path?: string }>;
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
    <div className="p-2 bg-orange-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Switch</div>
      <input
        value={data?.default_path ?? "default"}
        type="text"
        onChange={onDefaultPathChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Default path"
      />
      <div className="text-xs mt-1">Cases: {data?.cases?.length ?? 0}</div>
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-orange-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Parallel</div>
      <input
        value={data?.max_concurrency ?? 10}
        type="number"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Max concurrency"
      />
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-orange-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Join</div>
      <select
        value={data?.join_strategy ?? "all"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="all">All</option>
        <option value="any">Any</option>
        <option value="first">First</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

type SplitNodeData = {
  paths?: string[];
  label?: string;
};

export function SplitNode({ id, data }: NodeProps<SplitNodeData>) {
  const { setNodes } = useReactFlow();

  return (
    <div className="p-2 bg-pink-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Split</div>
      <div className="text-xs mt-1">Paths: {data?.paths?.length ?? 2}</div>
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-pink-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Delay</div>
      <input
        value={data?.duration ?? "1s"}
        type="text"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="1s, 100ms..."
      />
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-pink-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Cache</div>
      <select
        value={data?.cache_op ?? "get"}
        onChange={onOpChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="get">Get</option>
        <option value="set">Set</option>
        <option value="delete">Delete</option>
      </select>
      <input
        value={data?.cache_key ?? ""}
        type="text"
        onChange={onKeyChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Cache key"
      />
      {data?.cache_op === "set" && (
        <input
          value={data?.ttl ?? "5m"}
          type="text"
          onChange={onTTLChange}
          className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
          placeholder="TTL (5m, 1h)"
        />
      )}
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-red-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Retry</div>
      <input
        value={data?.max_attempts ?? 3}
        type="number"
        onChange={onAttemptsChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Max attempts"
      />
      <select
        value={data?.backoff_strategy ?? "exponential"}
        onChange={onStrategyChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="exponential">Exponential</option>
        <option value="linear">Linear</option>
        <option value="constant">Constant</option>
      </select>
      <input
        value={data?.initial_delay ?? "1s"}
        type="text"
        onChange={onInitialDelayChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="Initial delay"
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

type TryCatchNodeData = {
  fallback_value?: any;
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
    <div className="p-2 bg-red-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Try-Catch</div>
      <label className="flex items-center gap-1 mt-1">
        <input
          type="checkbox"
          checked={data?.continue_on_error ?? true}
          onChange={onContinueChange}
          className="text-sm"
        />
        <span className="text-xs">Continue on error</span>
      </label>
      <Handle type="source" position={Position.Right} />
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
    <div className="p-2 bg-red-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Timeout</div>
      <input
        value={data?.timeout ?? "30s"}
        type="text"
        onChange={onTimeoutChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
        placeholder="30s, 5m..."
      />
      <select
        value={data?.timeout_action ?? "error"}
        onChange={onActionChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black text-sm"
      >
        <option value="error">Error</option>
        <option value="continue_with_partial">Continue with partial</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
