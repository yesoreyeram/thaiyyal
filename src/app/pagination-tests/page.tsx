"use client";
import React, { useState } from "react";
import ReactFlow, {
  Background,
  ReactFlowProvider,
  Node as RFNode,
  Edge as RFEdge,
} from "reactflow";
import "reactflow/dist/style.css";
import {
  HttpNode,
  ConditionNode,
  CounterNode,
  AccumulatorNode,
} from "../../components/nodes";
import { NodeProps, Handle, Position, useReactFlow } from "reactflow";

type NodeData = Record<string, any>;

// Simple node components for pagination tests
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
    <div className="p-2 bg-gray-600 text-white shadow rounded border text-xs">
      <Handle type="target" position={Position.Left} />
      <div className="font-medium">Number</div>
      <input
        value={data?.value ?? 0}
        type="number"
        onChange={onChange}
        className="mt-1 w-20 border px-1 py-0.5 rounded text-black text-xs"
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
    <div className="p-2 bg-gray-700 text-white shadow rounded border text-xs">
      <Handle type="target" position={Position.Left} />
      <div className="font-medium">Operation</div>
      <select
        value={data?.op ?? "add"}
        onChange={onChange}
        className="mt-1 w-24 border px-1 py-0.5 rounded text-black text-xs"
      >
        <option value="add">Add</option>
        <option value="multiply">Multiply</option>
      </select>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

function VizNode({ id, data }: NodeProps<NodeData>) {
  return (
    <div className="p-2 bg-gray-800 text-white shadow rounded border text-xs">
      <Handle type="target" position={Position.Left} />
      <div className="font-medium">Viz</div>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

function TextInputNode({ id, data }: NodeProps<NodeData>) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, text } } : n))
    );
  };
  return (
    <div className="p-2 bg-green-600 text-white shadow rounded border text-xs">
      <Handle type="target" position={Position.Left} />
      <div className="font-medium">Text</div>
      <input
        value={data?.text ?? ""}
        onChange={onChange}
        className="mt-1 w-24 border px-1 py-0.5 rounded text-black text-xs"
      />
      <Handle type="source" position={Position.Right} />
    </div>
  );
}

const nodeTypes = {
  numberNode: NumberNode,
  opNode: OperationNode,
  vizNode: VizNode,
  textNode: TextInputNode,
  httpNode: HttpNode,
  conditionNode: ConditionNode,
  counterNode: CounterNode,
  accumulatorNode: AccumulatorNode,
};

// Pagination test scenarios
const paginationScenarios = [
  {
    title: "1. Page-Based Pagination",
    description: "Counter → Condition → HTTP (fetch page 1-5)",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { counter_op: "increment", delta: 1 }, type: "counterNode" },
      { id: "2", position: { x: 250, y: 100 }, data: { condition: "<=5" }, type: "conditionNode" },
      { id: "3", position: { x: 450, y: 100 }, data: { url: "https://api.example.com/items?page=1" }, type: "httpNode" },
      { id: "4", position: { x: 700, y: 100 }, data: {}, type: "vizNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "2" },
      { id: "e2", source: "2", target: "3" },
      { id: "e3", source: "3", target: "4" },
    ],
  },
  {
    title: "2. Offset-Based Pagination",
    description: "Page# → Multiply by PageSize → HTTP (with offset)",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { value: 0 }, type: "numberNode" },
      { id: "2", position: { x: 200, y: 100 }, data: { value: 10 }, type: "numberNode" },
      { id: "3", position: { x: 350, y: 150 }, data: { op: "multiply" }, type: "opNode" },
      { id: "4", position: { x: 550, y: 150 }, data: { url: "https://api.example.com/items?offset=0&limit=10" }, type: "httpNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "3" },
      { id: "e2", source: "2", target: "3" },
      { id: "e3", source: "3", target: "4" },
    ],
  },
  {
    title: "3. Cursor-Based Pagination",
    description: "Text (cursor) → HTTP → Extract next_cursor → Loop",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { text: "" }, type: "textNode" },
      { id: "2", position: { x: 250, y: 100 }, data: { url: "https://api.example.com/items?cursor=" }, type: "httpNode" },
      { id: "3", position: { x: 500, y: 100 }, data: {}, type: "vizNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "2" },
      { id: "e2", source: "2", target: "3" },
    ],
  },
  {
    title: "4. Until-Empty Pattern",
    description: "HTTP → Check response (has_more or count > 0) → Loop",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { url: "https://api.example.com/items?page=1" }, type: "httpNode" },
      { id: "2", position: { x: 350, y: 100 }, data: {}, type: "vizNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "2" },
    ],
  },
  {
    title: "5. With Accumulator",
    description: "HTTP → Accumulator (collect all page results)",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { url: "https://api.example.com/items" }, type: "httpNode" },
      { id: "2", position: { x: 300, y: 100 }, data: { text: "page_data" }, type: "textNode" },
      { id: "3", position: { x: 500, y: 100 }, data: { accum_op: "concat" }, type: "accumulatorNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "2" },
      { id: "e2", source: "2", target: "3" },
    ],
  },
  {
    title: "6. Complete Workflow",
    description: "Counter → Condition → HTTP → Accumulator → Viz",
    nodes: [
      { id: "1", position: { x: 50, y: 100 }, data: { counter_op: "increment" }, type: "counterNode" },
      { id: "2", position: { x: 200, y: 100 }, data: { condition: "<=5" }, type: "conditionNode" },
      { id: "3", position: { x: 350, y: 100 }, data: { url: "https://api.example.com?page=1" }, type: "httpNode" },
      { id: "4", position: { x: 600, y: 100 }, data: { text: "result" }, type: "textNode" },
      { id: "5", position: { x: 750, y: 100 }, data: {}, type: "vizNode" },
    ],
    edges: [
      { id: "e1", source: "1", target: "2" },
      { id: "e2", source: "2", target: "3" },
      { id: "e3", source: "3", target: "4" },
      { id: "e4", source: "4", target: "5" },
    ],
  },
];

function PaginationTestScenario({ scenario, index }: { scenario: any; index: number }) {
  const [showPayload, setShowPayload] = useState(false);

  const payload = {
    nodes: scenario.nodes.map((n: any) => ({
      id: n.id,
      type: n.type,
      data: n.data,
    })),
    edges: scenario.edges,
  };

  return (
    <div className="border rounded-lg p-4 bg-white shadow-sm">
      <h3 className="text-lg font-semibold mb-2">{scenario.title}</h3>
      <p className="text-sm text-gray-600 mb-3">{scenario.description}</p>

      <div style={{ height: "250px", border: "1px solid #ddd", borderRadius: "4px" }}>
        <ReactFlowProvider>
          <ReactFlow
            nodes={scenario.nodes as RFNode<NodeData>[]}
            edges={scenario.edges as RFEdge[]}
            nodeTypes={nodeTypes}
            fitView
            zoomOnScroll={false}
            panOnDrag={false}
            nodesDraggable={false}
            elementsSelectable={false}
          >
            <Background />
          </ReactFlow>
        </ReactFlowProvider>
      </div>

      <button
        onClick={() => setShowPayload(!showPayload)}
        className="mt-3 bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700"
      >
        {showPayload ? "Hide" : "Show"} JSON Payload
      </button>

      {showPayload && (
        <pre className="mt-3 p-3 bg-gray-100 rounded text-xs overflow-auto max-h-64">
          {JSON.stringify(payload, null, 2)}
        </pre>
      )}
    </div>
  );
}

export default function PaginationTestsPage() {
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-2">HTTP Pagination Test Scenarios</h1>
        <p className="text-gray-600 mb-6">
          Demonstrating composable pagination patterns using existing workflow nodes
        </p>

        <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mb-8">
          <h2 className="text-lg font-semibold mb-2">Composable Approach</h2>
          <p className="text-sm mb-2">
            These scenarios demonstrate how to achieve HTTP pagination using building blocks:
          </p>
          <ul className="text-sm space-y-1 ml-4">
            <li>• <strong>Counter</strong> - Track page numbers</li>
            <li>• <strong>Condition</strong> - Check page limits</li>
            <li>• <strong>HTTP</strong> - Make page requests</li>
            <li>• <strong>Accumulator</strong> - Collect results</li>
            <li>• <strong>Operation</strong> - Calculate offsets</li>
          </ul>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {paginationScenarios.map((scenario, index) => (
            <PaginationTestScenario key={index} scenario={scenario} index={index} />
          ))}
        </div>

        <div className="bg-green-50 border-l-4 border-green-500 p-4">
          <h2 className="text-lg font-semibold mb-2">✅ Test Coverage</h2>
          <p className="text-sm mb-2">
            {paginationScenarios.length} visual test scenarios covering:
          </p>
          <ul className="text-sm space-y-1 ml-4">
            <li>• Page-based pagination (most common)</li>
            <li>• Offset-based pagination (legacy APIs)</li>
            <li>• Cursor-based pagination (modern APIs)</li>
            <li>• Until-empty pattern (unknown page count)</li>
            <li>• Result accumulation across pages</li>
            <li>• Complete end-to-end workflow</li>
          </ul>
          <p className="text-sm mt-3">
            <strong>Backend Tests:</strong> 8 comprehensive test cases with mock HTTP servers
          </p>
        </div>
      </div>
    </div>
  );
}
