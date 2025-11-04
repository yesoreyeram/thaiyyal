import { NodePropsWithOptions } from "./nodeTypes";
import { Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";


// ===== ARRAY OPERATION NODES =====

type MapNodeData = {
  expression?: string;
  label?: string;
};

export function MapNode({ id, data, onShowOptions }: NodePropsWithOptions<MapNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
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

  const nodeInfo = getNodeInfo("mapNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Map")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-cyan-600 to-cyan-700 text-white shadow-lg rounded-lg border border-cyan-500 hover:border-cyan-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item * 2")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-cyan-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-cyan-400 focus:outline-none"
        placeholder="item * 2"
        aria-label="Map expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type ReduceNodeData = {
  expression?: string;
  initial_value?: string;
  label?: string;
};

export function ReduceNode({ id, data, onShowOptions }: NodePropsWithOptions<ReduceNodeData>) {
  const { setNodes } = useReactFlow();

  const onExpressionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
      )
    );
  };

  const onInitialValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const initial_value = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, initial_value } } : n
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

  const nodeInfo = getNodeInfo("reduceNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Reduce")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-teal-600 to-teal-700 text-white shadow-lg rounded-lg border border-teal-500 hover:border-teal-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex flex-col gap-0.5">
        <input
          value={String(data?.expression ?? "acc + item")}
          type="text"
          onChange={onExpressionChange}
          className="w-32 text-[10px] leading-tight border border-teal-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-teal-400 focus:outline-none"
          placeholder="acc + item"
          aria-label="Reduce expression"
        />
        <input
          value={String(data?.initial_value ?? "0")}
          type="text"
          onChange={onInitialValueChange}
          className="w-32 text-[10px] leading-tight border border-teal-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-teal-400 focus:outline-none"
          placeholder="Initial: 0"
          aria-label="Initial value"
        />
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type SliceNodeData = {
  start?: number;
  end?: number;
  label?: string;
};

export function SliceNode({ id, data, onShowOptions }: NodePropsWithOptions<SliceNodeData>) {
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

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("sliceNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Slice")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-emerald-600 to-emerald-700 text-white shadow-lg rounded-lg border border-emerald-500 hover:border-emerald-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex items-center gap-0.5">
        <input
          value={Number(data?.start ?? 0)}
          type="number"
          onChange={onStartChange}
          className="w-14 text-[10px] leading-tight border border-emerald-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-emerald-400 focus:outline-none"
          placeholder="Start"
          aria-label="Start index"
        />
        <span className="text-[10px]">to</span>
        <input
          value={Number(data?.end ?? -1)}
          type="number"
          onChange={onEndChange}
          className="w-14 text-[10px] leading-tight border border-emerald-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-emerald-400 focus:outline-none"
          placeholder="End"
          aria-label="End index"
        />
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type SortNodeData = {
  field?: string;
  order?: string;
  label?: string;
};

export function SortNode({ id, data, onShowOptions }: NodePropsWithOptions<SortNodeData>) {
  const { setNodes } = useReactFlow();

  const onFieldChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, field } } : n
      )
    );
  };

  const onOrderChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const order = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, order } } : n
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

  const nodeInfo = getNodeInfo("sortNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Sort")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-lime-600 to-lime-700 text-white shadow-lg rounded-lg border border-lime-500 hover:border-lime-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="flex items-center gap-0.5">
        <input
          value={String(data?.field ?? "")}
          type="text"
          onChange={onFieldChange}
          className="w-16 text-[10px] leading-tight border border-lime-600 px-1 py-0.5 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-1 focus:ring-lime-400 focus:outline-none"
          placeholder="Field"
          aria-label="Sort field"
        />
        <select
          value={String(data?.order ?? "asc")}
          onChange={onOrderChange}
          className="w-14 text-[10px] leading-tight border border-lime-600 px-1 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-lime-400 focus:outline-none"
        >
          <option value="asc">Asc</option>
          <option value="desc">Desc</option>
        </select>
      </div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type FindNodeData = {
  expression?: string;
  label?: string;
};

export function FindNode({ id, data, onShowOptions }: NodePropsWithOptions<FindNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
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

  const nodeInfo = getNodeInfo("findNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Find")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-sky-600 to-sky-700 text-white shadow-lg rounded-lg border border-sky-500 hover:border-sky-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item.id == 1")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-sky-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-sky-400 focus:outline-none"
        placeholder="item.id == 1"
        aria-label="Find expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type FlatMapNodeData = {
  expression?: string;
  label?: string;
};

export function FlatMapNode({ id, data, onShowOptions }: NodePropsWithOptions<FlatMapNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
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

  const nodeInfo = getNodeInfo("flatMapNode");

  return (
    <NodeWrapper
      title={String(data?.label || "FlatMap")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-indigo-600 to-indigo-700 text-white shadow-lg rounded-lg border border-indigo-500 hover:border-indigo-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item.values")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-indigo-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
        placeholder="item.values"
        aria-label="FlatMap expression"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type GroupByNodeData = {
  key_field?: string;
  label?: string;
};

export function GroupByNode({ id, data, onShowOptions }: NodePropsWithOptions<GroupByNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const key_field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, key_field } } : n
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

  const nodeInfo = getNodeInfo("groupByNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Group By")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-violet-600 to-violet-700 text-white shadow-lg rounded-lg border border-violet-500 hover:border-violet-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.key_field ?? "category")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-violet-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-violet-400 focus:outline-none"
        placeholder="category"
        aria-label="Group by field"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type UniqueNodeData = {
  by_field?: string;
  label?: string;
};

export function UniqueNode({ id, data, onShowOptions }: NodePropsWithOptions<UniqueNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const by_field = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, by_field } } : n
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

  const nodeInfo = getNodeInfo("uniqueNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Unique")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-fuchsia-600 to-fuchsia-700 text-white shadow-lg rounded-lg border border-fuchsia-500 hover:border-fuchsia-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.by_field ?? "")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-fuchsia-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-fuchsia-400 focus:outline-none"
        placeholder="Field (optional)"
        aria-label="Unique by field"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type ChunkNodeData = {
  size?: number;
  label?: string;
};

export function ChunkNode({ id, data, onShowOptions }: NodePropsWithOptions<ChunkNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const size = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, size } } : n
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

  const nodeInfo = getNodeInfo("chunkNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Chunk")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-pink-600 to-pink-700 text-white shadow-lg rounded-lg border border-pink-500 hover:border-pink-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={Number(data?.size ?? 3)}
        type="number"
        onChange={onChange}
        className="w-20 text-xs border border-pink-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-pink-400 focus:outline-none"
        placeholder="Size"
        aria-label="Chunk size"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type ReverseNodeData = {
  label?: string;
};

export function ReverseNode({ id, data, onShowOptions }: NodePropsWithOptions<ReverseNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("reverseNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Reverse")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-rose-600 to-rose-700 text-white shadow-lg rounded-lg border border-rose-500 hover:border-rose-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs py-1">Reverse array order</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type PartitionNodeData = {
  expression?: string;
  label?: string;
};

export function PartitionNode({ id, data, onShowOptions }: NodePropsWithOptions<PartitionNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const expression = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, expression } } : n
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

  const nodeInfo = getNodeInfo("partitionNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Partition")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-orange-600 to-orange-700 text-white shadow-lg rounded-lg border border-orange-500 hover:border-orange-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={String(data?.expression ?? "item > 0")}
        type="text"
        onChange={onChange}
        className="w-32 text-xs border border-orange-600 px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 focus:ring-orange-400 focus:outline-none"
        placeholder="item > 0"
        aria-label="Partition expression"
      />
      <Handle 
        type="source" 
        position={Position.Right} 
        id="true"
        style={{ top: '30%' }}
        className="w-2 h-2 bg-green-500"
        title="Match path"
      />
      <Handle 
        type="source" 
        position={Position.Right} 
        id="false"
        style={{ top: '70%' }}
        className="w-2 h-2 bg-red-500"
        title="No match path"
      />
    </NodeWrapper>
  );
}

type ZipNodeData = {
  label?: string;
};

export function ZipNode({ id, data, onShowOptions }: NodePropsWithOptions<ZipNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("zipNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Zip")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-yellow-600 to-yellow-700 text-white shadow-lg rounded-lg border border-yellow-500 hover:border-yellow-400 transition-all"
    >
      <Handle 
        type="target" 
        position={Position.Left} 
        id="array1"
        style={{ top: '30%' }}
        className="w-2 h-2 bg-blue-400"
        title="First array"
      />
      <Handle 
        type="target" 
        position={Position.Left} 
        id="array2"
        style={{ top: '70%' }}
        className="w-2 h-2 bg-blue-400"
        title="Second array"
      />
      <div className="text-xs py-1">Combine arrays</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type SampleNodeData = {
  count?: number;
  label?: string;
};

export function SampleNode({ id, data, onShowOptions }: NodePropsWithOptions<SampleNodeData>) {
  const { setNodes } = useReactFlow();

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const count = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, count } } : n
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

  const nodeInfo = getNodeInfo("sampleNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Sample")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-blue-600 to-blue-700 text-white shadow-lg rounded-lg border border-blue-500 hover:border-blue-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={Number(data?.count ?? 1)}
        type="number"
        onChange={onChange}
        className="w-20 text-xs border border-blue-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-400 focus:outline-none"
        placeholder="Count"
        aria-label="Sample count"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

type RangeNodeData = {
  start?: number;
  end?: number;
  step?: number;
  label?: string;
};

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

type TransposeNodeData = {
  label?: string;
};

export function TransposeNode({ id, data, onShowOptions }: NodePropsWithOptions<TransposeNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("transposeNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Transpose")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      className="bg-gradient-to-br from-red-600 to-red-700 text-white shadow-lg rounded-lg border border-red-500 hover:border-red-400 transition-all"
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <div className="text-xs py-1">Transpose matrix</div>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}
