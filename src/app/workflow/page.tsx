"use client";
import React, { useCallback, useMemo, useState, useEffect } from "react";
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
  FilterNode,
  BarChartNode,
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
  MapNode,
  ReduceNode,
  SliceNode,
  SortNode,
  FindNode,
  FlatMapNode,
  GroupByNode,
  UniqueNode,
  ChunkNode,
  ReverseNode,
  PartitionNode,
  ZipNode,
  SampleNode,
  RangeNode,
  TransposeNode,
  NodeContextMenu,
  DeleteConfirmDialog,
} from "../../components/nodes";
import { AppNavBar } from "../../components/AppNavBar";
import { WorkflowNavBar } from "../../components/WorkflowNavBar";
import { WorkflowStatusBar } from "../../components/WorkflowStatusBar";
import { NodePalette } from "../../components/NodePalette";
import { JSONPayloadModal } from "../../components/JSONPayloadModal";
import { useRouter } from "next/navigation";

type NodeData = Record<string, unknown>;

// Extended props to include onShowOptions and onOpenInfo
type NodePropsWithOptions = NodeProps<NodeData> & {
  onShowOptions?: (x: number, y: number) => void;
  onOpenInfo?: () => void;
};

// Higher-order component to add context menu and palette close to nodes
const withContextMenu = (
  Component: React.ComponentType<NodePropsWithOptions>, 
  handleContextMenu: (nodeId: string, x: number, y: number) => void,
  closePalette: () => void
) => {
  return (props: NodeProps<NodeData>) => {
    const onShowOptions = (x: number, y: number) => {
      handleContextMenu(props.id, x, y);
    };
    const onOpenInfo = () => {
      closePalette();
    };
    return <Component {...(props as NodePropsWithOptions)} onShowOptions={onShowOptions} onOpenInfo={onOpenInfo} />;
  };
};

// Original three node components - Updated to use NodeWrapper
function NumberNode({ id, data, onShowOptions, onOpenInfo, ...props }: NodePropsWithOptions) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const v = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, value: v } } : n
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

  // Import NodeWrapper at top
  const { NodeWrapper, getNodeInfo } = require("../../components/nodes");
  const nodeInfo = getNodeInfo("numberNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Number")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      onOpenInfo={onOpenInfo}
      
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={typeof data?.value === "number" ? data.value : 0}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

function OperationNode({ id, data, onShowOptions, ...props }: NodePropsWithOptions) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const op = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, op } } : n))
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const { NodeWrapper, getNodeInfo } = require("../../components/nodes");
  const nodeInfo = getNodeInfo("opNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Operation")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <select
        value={typeof data?.op === "string" ? data.op : "add"}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="add">Add</option>
        <option value="subtract">Subtract</option>
        <option value="multiply">Multiply</option>
        <option value="divide">Divide</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

function VizNode({ id, data, onShowOptions, ...props }: NodePropsWithOptions) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const mode = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, mode } } : n))
    );
  };

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const { NodeWrapper, getNodeInfo } = require("../../components/nodes");
  const nodeInfo = getNodeInfo("vizNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Visualization")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <select
        value={typeof data?.mode === "string" ? data.mode : "text"}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="text">Text</option>
        <option value="table">Table</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
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
      {
        type: "barChartNode",
        label: "Bar Chart",
        color: "bg-violet-600",
        defaultData: { 
          orientation: "vertical", 
          bar_color: "#3b82f6", 
          bar_width: "medium",
          show_values: true,
          max_bars: 20
        },
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
        type: "filterNode",
        label: "Filter",
        color: "bg-purple-600",
        defaultData: { condition: "item.age > 0" },
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
  {
    name: "Array Operations",
    nodes: [
      {
        type: "mapNode",
        label: "Map",
        color: "bg-cyan-600",
        defaultData: { expression: "item * 2" },
      },
      {
        type: "reduceNode",
        label: "Reduce",
        color: "bg-teal-600",
        defaultData: { expression: "acc + item", initial_value: "0" },
      },
      {
        type: "sliceNode",
        label: "Slice",
        color: "bg-emerald-600",
        defaultData: { start: 0, end: -1 },
      },
      {
        type: "sortNode",
        label: "Sort",
        color: "bg-lime-600",
        defaultData: { field: "", order: "asc" },
      },
      {
        type: "findNode",
        label: "Find",
        color: "bg-sky-600",
        defaultData: { expression: "item.id == 1" },
      },
      {
        type: "flatMapNode",
        label: "FlatMap",
        color: "bg-indigo-600",
        defaultData: { expression: "item.values" },
      },
      {
        type: "groupByNode",
        label: "Group By",
        color: "bg-violet-600",
        defaultData: { key_field: "category" },
      },
      {
        type: "uniqueNode",
        label: "Unique",
        color: "bg-fuchsia-600",
        defaultData: { by_field: "" },
      },
      {
        type: "chunkNode",
        label: "Chunk",
        color: "bg-pink-600",
        defaultData: { size: 3 },
      },
      {
        type: "reverseNode",
        label: "Reverse",
        color: "bg-rose-600",
        defaultData: {},
      },
      {
        type: "partitionNode",
        label: "Partition",
        color: "bg-orange-600",
        defaultData: { expression: "item > 0" },
      },
      {
        type: "zipNode",
        label: "Zip",
        color: "bg-yellow-600",
        defaultData: {},
      },
      {
        type: "sampleNode",
        label: "Sample",
        color: "bg-blue-600",
        defaultData: { count: 1 },
      },
      {
        type: "rangeNode",
        label: "Range",
        color: "bg-green-600",
        defaultData: { start: 0, end: 10, step: 1 },
      },
      {
        type: "transposeNode",
        label: "Transpose",
        color: "bg-red-600",
        defaultData: {},
      },
    ],
  },
];

function Canvas() {
  const router = useRouter();
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [showPayload, setShowPayload] = useState(false);
  const [isPaletteOpen, setIsPaletteOpen] = useState(false);
  const [workflowTitle, setWorkflowTitle] = useState("Untitled Workflow");
  const [contextMenu, setContextMenu] = useState<{ x: number; y: number; nodeId: string } | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ nodeId: string; nodeName: string } | null>(null);
  const { project, getNodes } = useReactFlow();

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds: RFEdge[]) => addEdge(params, eds)),
    [setEdges]
  );

  const handleNodeContextMenu = useCallback((nodeId: string, x: number, y: number) => {
    setContextMenu({ x, y, nodeId });
  }, []);

  const payload = useMemo(
    () => ({
      nodes: nodes.map((n) => ({ id: n.id, type: n.type, data: n.data, position: n.position })),
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
      numberNode: withContextMenu(NumberNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      opNode: withContextMenu(OperationNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      vizNode: withContextMenu(VizNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      barChartNode: withContextMenu(BarChartNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      textInputNode: withContextMenu(TextInputNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      textOpNode: withContextMenu(TextOperationNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      httpNode: withContextMenu(HttpNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      conditionNode: withContextMenu(ConditionNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      filterNode: withContextMenu(FilterNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      forEachNode: withContextMenu(ForEachNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      whileLoopNode: withContextMenu(WhileLoopNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      variableNode: withContextMenu(VariableNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      extractNode: withContextMenu(ExtractNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      transformNode: withContextMenu(TransformNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      accumulatorNode: withContextMenu(AccumulatorNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      counterNode: withContextMenu(CounterNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      switchNode: withContextMenu(SwitchNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      parallelNode: withContextMenu(ParallelNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      joinNode: withContextMenu(JoinNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      splitNode: withContextMenu(SplitNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      delayNode: withContextMenu(DelayNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      cacheNode: withContextMenu(CacheNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      retryNode: withContextMenu(RetryNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      tryCatchNode: withContextMenu(TryCatchNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      timeoutNode: withContextMenu(TimeoutNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      mapNode: withContextMenu(MapNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      reduceNode: withContextMenu(ReduceNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      sliceNode: withContextMenu(SliceNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      sortNode: withContextMenu(SortNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      findNode: withContextMenu(FindNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      flatMapNode: withContextMenu(FlatMapNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      groupByNode: withContextMenu(GroupByNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      uniqueNode: withContextMenu(UniqueNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      chunkNode: withContextMenu(ChunkNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      reverseNode: withContextMenu(ReverseNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      partitionNode: withContextMenu(PartitionNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      zipNode: withContextMenu(ZipNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      sampleNode: withContextMenu(SampleNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      rangeNode: withContextMenu(RangeNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
      transposeNode: withContextMenu(TransposeNode, handleNodeContextMenu, () => setIsPaletteOpen(false)),
    }),
    [handleNodeContextMenu]
  );

  const [nextId, setNextId] = useState(5);

  // Check for overlapping nodes
  const findNonOverlappingPosition = (basePosition: XYPosition): XYPosition => {
    const existingNodes = getNodes();
    const nodeWidth = 150;
    const nodeHeight = 80;
    const padding = 20;
    
    let position = { ...basePosition };
    let attempts = 0;
    const maxAttempts = 50;
    
    while (attempts < maxAttempts) {
      const overlaps = existingNodes.some(node => {
        const dx = Math.abs(node.position.x - position.x);
        const dy = Math.abs(node.position.y - position.y);
        return dx < (nodeWidth + padding) && dy < (nodeHeight + padding);
      });
      
      if (!overlaps) {
        return position;
      }
      
      // Try offset positions
      position = {
        x: basePosition.x + (attempts * 30),
        y: basePosition.y + ((attempts % 5) * 25),
      };
      attempts++;
    }
    
    return position;
  };

  const addNode = (type: string, defaultData: Record<string, unknown>) => {
    const id = String(nextId);
    setNextId((s) => s + 1);
    
    // Get viewport dimensions (accounting for nav bars: 14px app + 12px workflow + 7px status = 33px total)
    const navHeight = 112; // Total height of both navs + status bar
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight - navHeight;
    
    // Get base position at center of visible viewport
    const basePosition: XYPosition = project
      ? project({ x: viewportWidth / 2 - 75, y: viewportHeight / 2 })
      : { x: 400, y: 200 };
    
    // Find non-overlapping position
    const position = findNonOverlappingPosition(basePosition);

    const baseData: NodeData = { ...defaultData, label: `${type} ${id}` };
    const newNode: RFNode<NodeData> = { id, position, data: baseData, type };
    setNodes((nds) => nds.concat(newNode));
  };

  const handleNewWorkflow = () => {
    router.push('/');
  };

  const handleOpenWorkflow = () => {
    // TODO: Open workflow modal
  };

  const handleSave = () => {
    // TODO: Save workflow
    console.log('Save workflow', payload);
  };

  const handleDelete = () => {
    // TODO: Delete workflow
  };

  const handleRun = () => {
    // TODO: Run workflow
    console.log('Run workflow', payload);
  };

  const handleExport = () => {
    const jsonString = JSON.stringify(payload, null, 2);
    const blob = new Blob([jsonString], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    
    // Use workflow title for filename, sanitize it
    const sanitizedTitle = workflowTitle.replace(/[^a-z0-9]/gi, '_').toLowerCase();
    link.download = `${sanitizedTitle}_workflow.json`;
    
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const handleImport = (data: { nodes: unknown[]; edges: unknown[] }) => {
    // Type cast and validate the imported data
    const importedNodes = data.nodes as RFNode<NodeData>[];
    const importedEdges = data.edges as RFEdge[];
    
    // Ensure all nodes have position data, add default position if missing
    const nodesWithPositions = importedNodes.map((node, index) => ({
      ...node,
      position: node.position || { x: 50 + index * 200, y: 50 + index * 100 }
    }));
    
    // Set the imported nodes and edges
    setNodes(nodesWithPositions);
    setEdges(importedEdges);
    
    // Update the next ID to be higher than any existing node ID
    const maxId = importedNodes.reduce((max, node) => {
      const nodeId = parseInt(node.id, 10);
      return isNaN(nodeId) ? max : Math.max(max, nodeId);
    }, 0);
    setNextId(maxId + 1);
  };

  const handleDeleteNode = (nodeId: string) => {
    const node = nodes.find((n) => n.id === nodeId);
    if (node) {
      setDeleteConfirm({ nodeId, nodeName: String(node.data?.label || `Node ${nodeId}`) });
    }
    setContextMenu(null);
  };

  const confirmDelete = () => {
    if (deleteConfirm) {
      setNodes((nds) => nds.filter((n) => n.id !== deleteConfirm.nodeId));
      setEdges((eds) => eds.filter((e) => e.source !== deleteConfirm.nodeId && e.target !== deleteConfirm.nodeId));
    }
    setDeleteConfirm(null);
  };

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      {/* Application Nav Bar */}
      <AppNavBar
        onNewWorkflow={handleNewWorkflow}
        onOpenWorkflow={handleOpenWorkflow}
      />

      {/* Workflow Nav Bar */}
      <WorkflowNavBar
        workflowTitle={workflowTitle}
        onTitleChange={setWorkflowTitle}
        onSave={handleSave}
        onShowJSON={() => setShowPayload(true)}
        onDelete={handleDelete}
        onRun={handleRun}
        onExport={handleExport}
        onImport={handleImport}
      />

      {/* Main Content */}
      <div className="flex-1 relative">
        {/* Add Node Button - Bottom Left */}
        {!isPaletteOpen && (
          <button
            onClick={() => setIsPaletteOpen(true)}
            className="absolute left-4 bottom-4 z-10 bg-gray-800 hover:bg-gray-700 text-white px-3 py-1.5 rounded-lg shadow-lg transition-all border border-gray-700 hover:border-gray-600 text-sm font-medium flex items-center gap-1.5"
            title="Add Node"
          >
            <span className="text-base">+</span>
            <span>Add Node</span>
          </button>
        )}

        {/* Node Palette with Search */}
        <NodePalette
          isOpen={isPaletteOpen}
          onClose={() => setIsPaletteOpen(false)}
          categories={nodeCategories}
          onAddNode={addNode}
        />

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

        {/* JSON Payload Modal */}
        <JSONPayloadModal
          isOpen={showPayload}
          onClose={() => setShowPayload(false)}
          payload={payload}
          workflowTitle={workflowTitle}
        />
      </div>

      {/* Context Menu */}
      {contextMenu && (
        <NodeContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onClose={() => setContextMenu(null)}
          onDelete={() => handleDeleteNode(contextMenu.nodeId)}
        />
      )}

      {/* Delete Confirmation Dialog */}
      {deleteConfirm && (
        <DeleteConfirmDialog
          nodeName={deleteConfirm.nodeName}
          onConfirm={confirmDelete}
          onCancel={() => setDeleteConfirm(null)}
        />
      )}

      {/* Bottom Status Bar */}
      <WorkflowStatusBar
        nodeCount={nodes.length}
        edgeCount={edges.length}
      />
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
