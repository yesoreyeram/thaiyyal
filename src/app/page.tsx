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

type NodeData = { value?: number; op?: string; mode?: string; label?: string };

// custom node components
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
    <div className="p-2 bg-gray text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Number</div>
      <input
        value={data?.value ?? 0}
        type="number"
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded"
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
    <div className="p-2 bg-gray text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Operation</div>
      <select
        value={data?.op ?? "add"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded"
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
    <div className="p-2 bg-gray text-white shadow rounded border">
      <Handle type="target" position={Position.Left} />
      <div className="text-xs font-medium">Visualization</div>
      <select
        value={data?.mode ?? "text"}
        onChange={onChange}
        className="mt-1 w-32 border px-2 py-1 rounded"
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
      nodes: nodes.map((n) => ({ id: n.id, data: n.data })),
      edges: edges.map((e) => ({
        id: e.id,
        source: e.source,
        target: e.target,
      })),
    }),
    [nodes, edges]
  );

  const nodeTypes = useMemo(
    () => ({ numberNode: NumberNode, opNode: OperationNode, vizNode: VizNode }),
    []
  );

  const [nextId, setNextId] = useState(5);
  const addNode = (type: "numberNode" | "opNode" | "vizNode") => {
    const id = String(nextId);
    setNextId((s) => s + 1);
    const position: XYPosition = project
      ? project({ x: 100, y: 100 })
      : { x: 400 + nextId * 10, y: 120 + (nextId % 3) * 40 };
    const baseData: NodeData =
      type === "numberNode"
        ? { value: 0, label: `Node ${id}` }
        : type === "opNode"
        ? { op: "add", label: `Op ${id}` }
        : { mode: "text", label: `Viz ${id}` };
    const newNode: RFNode<NodeData> = { id, position, data: baseData, type };
    setNodes((nds) => nds.concat(newNode));
  };

  return (
    <div className="h-screen flex">
      <div className="w-1/2 border-r flex flex-col relative">
        <div className="absolute left-2 top-4 z-10 flex flex-col gap-2">
          <button
            onClick={() => addNode("numberNode")}
            className="bg-gray text-white border px-3 py-2 rounded shadow"
          >
            + Number
          </button>
          <button
            onClick={() => addNode("opNode")}
            className="bg-gray text-white border px-3 py-2 rounded shadow"
          >
            + Operation
          </button>
          <button
            onClick={() => addNode("vizNode")}
            className="bg-gray text-white border px-3 py-2 rounded shadow"
          >
            + Viz
          </button>
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
            {/* <Controls /> */}
          </ReactFlow>
        </div>
        <div className="p-3 border-t bg-white flex items-center justify-between">
          <div className="text-sm text-gray-600">
            Canvas (use toolbar to add nodes)
          </div>
          <div>
            <button
              onClick={() => setShow((s) => !s)}
              className="bg-blue-600 text-white px-4 py-2 rounded"
            >
              Show payload
            </button>
          </div>
        </div>
      </div>

      <div className="w-1/2 p-4">
        <h3 className="text-lg font-medium mb-2">Generated JSON Payload</h3>
        <div className="h-[80vh] overflow-auto p-3 bg-white rounded shadow payload-box">
          {show ? (
            <pre className="text-black">{JSON.stringify(payload, null, 2)}</pre>
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
