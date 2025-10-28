"use client";
import React, { useCallback, useMemo, useState } from "react";
import ReactFlow, {
  addEdge,
  Background,
  Controls,
  MiniMap,
  ReactFlowProvider,
  useEdgesState,
  useNodesState,
  Node as RFNode,
  Edge as RFEdge,
  Connection,
} from "reactflow";
import "reactflow/dist/style.css";

type NodeData = { value?: number; op?: string; mode?: string; label?: string };

const initialNodes: RFNode<NodeData>[] = [
  {
    id: "1",
    position: { x: 50, y: 50 },
    data: { value: 10, label: "Node 1" },
    type: "default",
  },
  {
    id: "2",
    position: { x: 50, y: 200 },
    data: { value: 5, label: "Node 2" },
    type: "default",
  },
  {
    id: "3",
    position: { x: 300, y: 120 },
    data: { op: "add", label: "Node 3 (op)" },
    type: "default",
  },
  {
    id: "4",
    position: { x: 600, y: 120 },
    data: { mode: "text", label: "Node 4 (viz)" },
    type: "default",
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

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds: RFEdge[]) => addEdge(params, eds)),
    [setEdges]
  );

  const payload = useMemo(() => {
    return {
      nodes: nodes.map((n: RFNode<NodeData>) => ({ id: n.id, data: n.data })),
      edges: edges.map((e: RFEdge) => ({
        id: e.id,
        source: e.source,
        target: e.target,
      })),
    };
  }, [nodes, edges]);

  // simple UI to edit node values inline using prompt (minimal)
  const editNode = (id: string) => {
    const n = nodes.find((x) => x.id === id);
    if (!n) return;
    if (id === "3") {
      const op = prompt(
        "Enter operation (add, subtract, multiply, divide)",
        String(n.data.op ?? "add")
      );
      if (op)
        setNodes((nds) =>
          nds.map((m) => (m.id === id ? { ...m, data: { ...m.data, op } } : m))
        );
    } else if (id === "4") {
      const mode = prompt(
        "Enter viz mode (text/table)",
        String(n.data.mode ?? "text")
      );
      if (mode)
        setNodes((nds) =>
          nds.map((m) =>
            m.id === id ? { ...m, data: { ...m.data, mode } } : m
          )
        );
    } else {
      const v = prompt("Enter number", String(n.data.value ?? 0));
      if (v !== null)
        setNodes((nds) =>
          nds.map((m) =>
            m.id === id ? { ...m, data: { ...m.data, value: Number(v) } } : m
          )
        );
    }
  };

  return (
    <div className="h-screen flex">
      <div className="w-1/2 border-r flex flex-col">
        <div className="flex-1">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeClick={(_, node) => editNode(node.id)}
            fitView
          >
            <Background />
            <Controls />
          </ReactFlow>
        </div>
        <div className="p-3 border-t bg-white flex items-center justify-between">
          <div className="text-sm text-gray-600">
            Canvas (click nodes to edit)
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
        <pre className="text-white">
          {JSON.stringify(payload || "{}", null, 2)}
        </pre>
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
