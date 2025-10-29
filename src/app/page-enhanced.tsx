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
} from "../components/nodes";

type NodeData = Record<string, any>;

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
    <div className="p-2 bg-gray-600 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Number</div>
      <input
        value={data?.value ?? 0}
        type="number"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black"
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
    <div className="p-2 bg-gray-700 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Operation</div>
      <select
        value={data?.op ?? "add"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black"
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
    <div className="p-2 bg-gray-800 text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Visualization</div>
      <select
        value={data?.mode ?? "text"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded text-black"
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

// Node type configurations
const nodeTypeConfigs = [
  { type: "numberNode", label: "Number", color: "gray", defaultData: { value: 0 } },
  { type: "opNode", label: "Operation", color: "gray", defaultData: { op: "add" } },
  { type: "vizNode", label: "Visualization", color: "gray", defaultData: { mode: "text" } },
  { type: "textInputNode", label: "Text Input", color: "green", defaultData: { text: "" } },
  { type: "textOpNode", label: "Text Operation", color: "green", defaultData: { text_op: "uppercase" } },
  { type: "httpNode", label: "HTTP", color: "purple", defaultData: { url: "" } },
  { type: "conditionNode", label: "Condition", color: "yellow", defaultData: { condition: ">0" } },
  { type: "forEachNode", label: "For Each", color: "yellow", defaultData: { max_iterations: 1000 } },
  { type: "whileLoopNode", label: "While Loop", color: "yellow", defaultData: { condition: ">0", max_iterations: 100 } },
  { type: "variableNode", label: "Variable", color: "blue", defaultData: { var_name: "", var_op: "get" } },
  { type: "extractNode", label: "Extract", color: "blue", defaultData: { field: "" } },
  { type: "transformNode", label: "Transform", color: "blue", defaultData: { transform_type: "to_array" } },
  { type: "accumulatorNode", label: "Accumulator", color: "indigo", defaultData: { accum_op: "sum" } },
  { type: "counterNode", label: "Counter", color: "indigo", defaultData: { counter_op: "increment", delta: 1 } },
  { type: "switchNode", label: "Switch", color: "orange", defaultData: { cases: [], default_path: "default" } },
  { type: "parallelNode", label: "Parallel", color: "orange", defaultData: { max_concurrency: 10 } },
  { type: "joinNode", label: "Join", color: "orange", defaultData: { join_strategy: "all" } },
  { type: "splitNode", label: "Split", color: "pink", defaultData: { paths: ["path1", "path2"] } },
  { type: "delayNode", label: "Delay", color: "pink", defaultData: { duration: "1s" } },
  { type: "cacheNode", label: "Cache", color: "pink", defaultData: { cache_op: "get", cache_key: "" } },
  { type: "retryNode", label: "Retry", color: "red", defaultData: { max_attempts: 3, backoff_strategy: "exponential", initial_delay: "1s" } },
  { type: "tryCatchNode", label: "Try-Catch", color: "red", defaultData: { continue_on_error: true } },
  { type: "timeoutNode", label: "Timeout", color: "red", defaultData: { timeout: "30s", timeout_action: "error" } },
];

function Canvas() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [show, setShow] = useState(false);
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
  
  const addNode = (type: string, defaultData: any) => {
    const id = String(nextId);
    setNextId((s) => s + 1);
    const position: XYPosition = project
      ? project({ x: 100, y: 100 })
      : { x: 400 + nextId * 10, y: 120 + (nextId % 3) * 40 };
    
    const baseData: NodeData = { ...defaultData, label: `${type} ${id}` };
    const newNode: RFNode<NodeData> = { id, position, data: baseData, type };
    setNodes((nds) => nds.concat(newNode));
  };

  return (
    <div className="h-screen flex">
      <div className="w-1/2 border-r flex flex-col relative">
        <div className="absolute left-2 top-4 z-10 flex flex-col gap-1 max-h-[calc(100vh-100px)] overflow-y-auto bg-white/90 p-2 rounded shadow-lg">
          <div className="text-xs font-bold mb-1">Add Nodes:</div>
          {nodeTypeConfigs.map((config) => (
            <button
              key={config.type}
              onClick={() => addNode(config.type, config.defaultData)}
              className={`bg-${config.color}-600 text-white border px-2 py-1 rounded shadow text-xs hover:opacity-80`}
              style={{
                backgroundColor: `var(--${config.color}-600, #6b7280)`,
              }}
            >
              + {config.label}
            </button>
          ))}
        </div>
        <div className="flex-1">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            nodeTypes={nodeTypes}
            fitView
          >
            <Background />
          </ReactFlow>
        </div>
        <div className="p-3 border-t bg-white flex items-center justify-between">
          <div className="text-sm text-gray-600">
            Canvas ({nodes.length} nodes, {edges.length} edges)
          </div>
          <div>
            <button
              onClick={() => setShow((s) => !s)}
              className="bg-blue-600 text-white px-4 py-2 rounded"
            >
              {show ? "Hide" : "Show"} payload
            </button>
          </div>
        </div>
      </div>

      <div className="w-1/2 p-4">
        <h3 className="text-lg font-medium mb-2">Generated JSON Payload</h3>
        <div className="h-[80vh] overflow-auto p-3 bg-white rounded shadow payload-box">
          {show ? (
            <pre className="text-black text-xs">{JSON.stringify(payload, null, 2)}</pre>
          ) : (
            <div className="text-sm text-gray-500">
              Click &quot;Show payload&quot; on the canvas to view generated
              JSON.
            </div>
          )}
        </div>
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
