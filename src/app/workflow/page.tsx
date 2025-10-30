"use client";
import React, { useCallback, useMemo, useState } from "react";
import ReactFlow, {
  addEdge,
  Background,
  ReactFlowProvider,
  useEdgesState,
  useNodesState,
  useReactFlow,
  XYPosition,
  Node as RFNode,
  Edge as RFEdge,
  Connection,
  NodeProps,
  Handle,
  Position,
} from "reactflow";
import "reactflow/dist/style.css";
import {
  TextInputNode,
  TextOperationNode,
  HttpNode,
  ConditionNode,
  ForEachNode,
  WhileLoopNode,
  VariableNode,
  ExtractNode,
  TransformNode,
  AccumulatorNode,
  CounterNode,
  SwitchNode,
  ParallelNode,
  JoinNode,
  SplitNode,
  DelayNode,
  CacheNode,
  RetryNode,
  TryCatchNode,
  TimeoutNode,
} from "../../components/nodes";

type NodeData = Record<string, unknown>;

// Original three node components
function NumberNode({ id, data }: NodeProps<NodeData>) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const v = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, value: v } } : n
      )
    );
  };
  return (
    <div className="p-2 bg-gray-700 text-white shadow rounded border border-gray-600">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Number</div>
      <input
        value={typeof data?.value === "number" ? data.value : 0}
        type="number"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded bg-gray-800 text-white border-gray-600"
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

function OperationNode({ id, data }: NodeProps<NodeData>) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const op = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, op } } : n))
    );
  };
  return (
    <div className="p-2 bg-gray-700 text-white shadow rounded border border-gray-600">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Operation</div>
      <select
        value={typeof data?.op === "string" ? data.op : "add"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded bg-gray-800 text-white border-gray-600"
      >
        <option value="add">Add</option>
        <option value="subtract">Subtract</option>
        <option value="multiply">Multiply</option>
        <option value="divide">Divide</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

function VizNode({ id, data }: NodeProps<NodeData>) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const mode = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, mode } } : n))
    );
  };
  return (
    <div className="p-2 bg-gray-700 text-white shadow rounded border border-gray-600">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Visualization</div>
      <select
        value={typeof data?.mode === "string" ? data.mode : "text"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded bg-gray-800 text-white border-gray-600"
      >
        <option value="text">Text</option>
        <option value="table">Table</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

const initialNodes: RFNode<NodeData>[] = [
  {
    id: "1",
    position: { x: 50, y: 50 },
    data: { value: 10, label: "Node 1" },
    type: "numberNode",
  },
  {
    id: "2",
    position: { x: 50, y: 200 },
    data: { value: 5, label: "Node 2" },
    type: "numberNode",
  },
  {
    id: "3",
    position: { x: 300, y: 120 },
    data: { op: "add", label: "Node 3 (op)" },
    type: "opNode",
  },
  {
    id: "4",
    position: { x: 600, y: 120 },
    data: { mode: "text", label: "Node 4 (viz)" },
    type: "vizNode",
  },
];

const initialEdges: RFEdge[] = [
  { id: "e1-3", source: "1", target: "3" },
  { id: "e2-3", source: "2", target: "3" },
  { id: "e3-4", source: "3", target: "4" },
];

// Node categories configuration
const nodeCategories = [
  {
    name: "Input/Output",
    nodes: [
      {
        type: "numberNode",
        label: "Number",
        color: "bg-blue-600",
        defaultData: { value: 0 },
      },
      {
        type: "textInputNode",
        label: "Text Input",
        color: "bg-green-600",
        defaultData: { text: "" },
      },
      {
        type: "httpNode",
        label: "HTTP",
        color: "bg-purple-600",
        defaultData: { url: "" },
      },
      {
        type: "vizNode",
        label: "Visualization",
        color: "bg-indigo-600",
        defaultData: { mode: "text" },
      },
    ],
  },
  {
    name: "Operations",
    nodes: [
      {
        type: "opNode",
        label: "Math Op",
        color: "bg-yellow-600",
        defaultData: { op: "add" },
      },
      {
        type: "textOpNode",
        label: "Text Op",
        color: "bg-green-600",
        defaultData: { text_op: "uppercase" },
      },
      {
        type: "transformNode",
        label: "Transform",
        color: "bg-cyan-600",
        defaultData: { transform_type: "to_array" },
      },
      {
        type: "extractNode",
        label: "Extract",
        color: "bg-teal-600",
        defaultData: { field: "" },
      },
    ],
  },
  {
    name: "Control Flow",
    nodes: [
      {
        type: "conditionNode",
        label: "Condition",
        color: "bg-amber-600",
        defaultData: { condition: ">0" },
      },
      {
        type: "forEachNode",
        label: "For Each",
        color: "bg-orange-600",
        defaultData: { max_iterations: 1000 },
      },
      {
        type: "whileLoopNode",
        label: "While Loop",
        color: "bg-red-600",
        defaultData: { condition: ">0", max_iterations: 100 },
      },
      {
        type: "switchNode",
        label: "Switch",
        color: "bg-pink-600",
        defaultData: { cases: [], default_path: "default" },
      },
    ],
  },
  {
    name: "Parallel & Join",
    nodes: [
      {
        type: "parallelNode",
        label: "Parallel",
        color: "bg-violet-600",
        defaultData: { max_concurrency: 10 },
      },
      {
        type: "joinNode",
        label: "Join",
        color: "bg-fuchsia-600",
        defaultData: { join_strategy: "all" },
      },
      {
        type: "splitNode",
        label: "Split",
        color: "bg-rose-600",
        defaultData: { paths: ["path1", "path2"] },
      },
    ],
  },
  {
    name: "State & Memory",
    nodes: [
      {
        type: "variableNode",
        label: "Variable",
        color: "bg-sky-600",
        defaultData: { var_name: "", var_op: "get" },
      },
      {
        type: "accumulatorNode",
        label: "Accumulator",
        color: "bg-blue-500",
        defaultData: { accum_op: "sum" },
      },
      {
        type: "counterNode",
        label: "Counter",
        color: "bg-indigo-500",
        defaultData: { counter_op: "increment", delta: 1 },
      },
      {
        type: "cacheNode",
        label: "Cache",
        color: "bg-purple-500",
        defaultData: { cache_op: "get", cache_key: "" },
      },
    ],
  },
  {
    name: "Error Handling",
    nodes: [
      {
        type: "retryNode",
        label: "Retry",
        color: "bg-red-500",
        defaultData: {
          max_attempts: 3,
          backoff_strategy: "exponential",
          initial_delay: "1s",
        },
      },
      {
        type: "tryCatchNode",
        label: "Try-Catch",
        color: "bg-orange-500",
        defaultData: { continue_on_error: true },
      },
      {
        type: "timeoutNode",
        label: "Timeout",
        color: "bg-amber-500",
        defaultData: { timeout: "30s", timeout_action: "error" },
      },
    ],
  },
  {
    name: "Utilities",
    nodes: [
      {
        type: "delayNode",
        label: "Delay",
        color: "bg-gray-500",
        defaultData: { duration: "1s" },
      },
    ],
  },
];

function Canvas() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [showPayload, setShowPayload] = useState(false);
  const [isPaletteOpen, setIsPaletteOpen] = useState(false);
  const { project } = useReactFlow();

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds: RFEdge[]) => addEdge(params, eds)),
    [setEdges]
  );

  const payload = useMemo(
    () => ({
      nodes: nodes.map((n) => ({ id: n.id, type: n.type, data: n.data })),
      edges: edges.map((e) => ({
        id: e.id,
        source: e.source,
        target: e.target,
      })),
    }),
    [nodes, edges]
  );

  const nodeTypes = useMemo(
    () => ({
      numberNode: NumberNode,
      opNode: OperationNode,
      vizNode: VizNode,
      textInputNode: TextInputNode,
      textOpNode: TextOperationNode,
      httpNode: HttpNode,
      conditionNode: ConditionNode,
      forEachNode: ForEachNode,
      whileLoopNode: WhileLoopNode,
      variableNode: VariableNode,
      extractNode: ExtractNode,
      transformNode: TransformNode,
      accumulatorNode: AccumulatorNode,
      counterNode: CounterNode,
      switchNode: SwitchNode,
      parallelNode: ParallelNode,
      joinNode: JoinNode,
      splitNode: SplitNode,
      delayNode: DelayNode,
      cacheNode: CacheNode,
      retryNode: RetryNode,
      tryCatchNode: TryCatchNode,
      timeoutNode: TimeoutNode,
    }),
    []
  );

  const [nextId, setNextId] = useState(5);

  const addNode = (type: string, defaultData: Record<string, unknown>) => {
    const id = String(nextId);
    setNextId((s) => s + 1);
    const position: XYPosition = project
      ? project({ x: 100, y: 100 })
      : { x: 400 + nextId * 10, y: 120 + (nextId % 3) * 40 };

    const baseData: NodeData = { ...defaultData, label: `${type} ${id}` };
    const newNode: RFNode<NodeData> = { id, position, data: baseData, type };
    setNodes((nds) => nds.concat(newNode));
    setIsPaletteOpen(false); // Close palette after adding a node
  };

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      {/* Top Bar */}
      <div className="h-14 bg-gray-900 border-b border-gray-700 flex items-center justify-between px-4">
        <div className="flex items-center gap-4">
          <h1 className="text-xl font-bold text-white">Workflow Builder</h1>
          <div className="text-sm text-gray-400">
            {nodes.length} nodes, {edges.length} connections
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowPayload((s) => !s)}
            className="bg-gray-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm transition-colors"
          >
            {showPayload ? "Hide" : "View"} JSON Payload
          </button>
          <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm transition-colors">
            ▶︎
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 relative">
        {/* Add Node Button (when palette is closed) */}
        {!isPaletteOpen && (
          <button
            onClick={() => setIsPaletteOpen(true)}
            className="absolute left-4 top-4 z-10 bg-blue-600 hover:bg-blue-700 text-white p-3 rounded-full shadow-lg transition-all hover:scale-110"
            title="Add Node"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={2}
              stroke="currentColor"
              className="w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 4.5v15m7.5-7.5h-15"
              />
            </svg>
          </button>
        )}

        {/* Collapsible Floating Node Palette */}
        {isPaletteOpen && (
          <div className="absolute left-4 top-4 z-10 bg-gray-900 border border-gray-700 rounded-lg shadow-2xl max-h-[calc(100vh-120px)] overflow-y-auto w-64">
            <div className="sticky top-0 bg-gray-900 border-b border-gray-700 p-3 flex items-center justify-between">
              <div className="text-sm font-bold text-white">Add Nodes</div>
              <button
                onClick={() => setIsPaletteOpen(false)}
                className="text-gray-400 hover:text-white transition-colors"
                title="Close"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className="w-5 h-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>

            {nodeCategories.map((category) => (
              <div
                key={category.name}
                className="p-3 border-b border-gray-800 last:border-b-0"
              >
                <div className="text-xs font-semibold text-gray-400 mb-2 uppercase tracking-wide">
                  {category.name}
                </div>
                <div className="flex flex-col gap-1">
                  {category.nodes.map((config) => (
                    <button
                      key={config.type}
                      onClick={() => addNode(config.type, config.defaultData)}
                      className={`${config.color} hover:opacity-80 text-white px-3 py-2 rounded text-sm transition-all text-left`}
                    >
                      + {config.label}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* React Flow Canvas */}
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          fitView
          className="bg-gray-950"
        >
          <Background className="bg-gray-950" />
        </ReactFlow>

        {/* JSON Payload Panel */}
        {showPayload && (
          <div className="absolute right-4 top-4 bottom-4 w-96 bg-gray-900 border border-gray-700 rounded-lg shadow-2xl overflow-hidden flex flex-col">
            <div className="bg-gray-800 border-b border-gray-700 p-3 flex items-center justify-between">
              <div className="text-sm font-bold text-white">JSON Payload</div>
              <button
                onClick={() => setShowPayload(false)}
                className="text-gray-400 hover:text-white transition-colors"
                title="Close"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className="w-5 h-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
            <div className="flex-1 overflow-auto p-4">
              <pre className="text-gray-300 text-xs">
                {JSON.stringify(payload, null, 2)}
              </pre>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function Page() {
  return (
    <ReactFlowProvider>
      <Canvas />
    </ReactFlowProvider>
  );
}
