"use client";
import React, {
  useCallback,
  useMemo,
  useState,
  useRef,
  useEffect,
} from "react";
import ReactFlow, {
  addEdge,
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
  NodeChange,
} from "reactflow";
import "reactflow/dist/style.css";
import {
  TextInputNode,
  BooleanInputNode,
  DateInputNode,
  DateTimeInputNode,
  TextOperationNode,
  HttpNode,
  ConditionNode,
  FilterNode,
  RendererNode,
  ForEachNode,
  WhileLoopNode,
  VariableNode,
  ExtractNode,
  TransformNode,
  ParseNode,
  ExpressionNode,
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
  NodeWrapper,
  getNodeInfo,
  NodePropsWithOptions,
} from "../../components/nodes";
import { AppNavBar } from "../../components/AppNavBar";
import { WorkflowNavBar } from "../../components/WorkflowNavBar";
import { WorkflowStatusBar } from "../../components/WorkflowStatusBar";
import { NodePalette } from "../../components/NodePalette";
import { JSONPayloadModal } from "../../components/JSONPayloadModal";
import { WorkflowExamplesModal } from "../../components/WorkflowExamplesModal";
import {
  ExecutionPanel,
  ExecutionResult,
} from "../../components/ExecutionPanel";
import { WorkflowExample } from "../../data/workflowExamples";

type NodeData = Record<string, unknown>;

// Higher-order component to add context menu and palette close to nodes
const withContextMenu = (
  Component: React.ComponentType<NodePropsWithOptions>,
  handleContextMenu: (nodeId: string, x: number, y: number) => void,
  closePalette: () => void,
  handleDelete: (nodeId: string) => void
) => {
  const WrappedComponent = (props: NodeProps<NodeData>) => {
    const onShowOptions = (x: number, y: number) => {
      handleContextMenu(props.id, x, y);
    };
    const onOpenInfo = () => {
      closePalette();
    };
    const onDelete = () => {
      handleDelete(props.id);
    };
    return (
      <Component
        {...(props as NodePropsWithOptions)}
        onShowOptions={onShowOptions}
        onOpenInfo={onOpenInfo}
        onDelete={onDelete}
      />
    );
  };
  WrappedComponent.displayName = `withContextMenu(${
    Component.displayName || Component.name || "Component"
  })`;
  return WrappedComponent;
};

// Original three node components - Updated to use NodeWrapper
function NumberNode({
  id,
  data,
  onShowOptions,
  onOpenInfo,
}: NodePropsWithOptions) {
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

  const nodeInfo = getNodeInfo("numberNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Number")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
      onOpenInfo={onOpenInfo}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <input
        value={typeof data?.value === "number" ? data.value : 0}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}

function OperationNode({ id, data, onShowOptions }: NodePropsWithOptions) {
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

  const nodeInfo = getNodeInfo("opNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Operation")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
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
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}

function VizNode({ id, data, onShowOptions }: NodePropsWithOptions) {
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

  const nodeInfo = getNodeInfo("vizNode");

  return (
    <NodeWrapper
      title={String(data?.label || "Visualization")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      <select
        value={typeof data?.mode === "string" ? data.mode : "text"}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-1.5 py-0.5 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
      >
        <option value="text">Text</option>
        <option value="table">Table</option>
      </select>
      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}

const initialNodes: RFNode<NodeData>[] = [
  {
    id: "1",
    type: "number",
    data: {
      value: 10,
      label: "Node 1",
    },
    position: {
      x: 64.40527775914373,
      y: -46.56871386685241,
    },
  },
  {
    id: "2",
    type: "number",
    data: {
      value: 7,
      label: "Node 2",
    },
    position: {
      x: 66.00586417682638,
      y: 41.00841584352477,
    },
  },
  {
    id: "3",
    type: "operation",
    data: {
      op: "subtract",
      label: "Node 3 (op)",
    },
    position: {
      x: 289.32942388211575,
      y: -5.379269385139878,
    },
  },
  {
    id: "5",
    type: "expression",
    data: {
      expression: "input * 2",
      label: "expression 5",
    },
    position: {
      x: 544.1813952845607,
      y: -58.416722398767256,
    },
  },
  {
    id: "6",
    type: "renderer",
    data: {
      label: "renderer 6",
      _executionData: 6,
    },
    position: {
      x: 739.4878012390735,
      y: -72.0625848235367,
    },
  },
  {
    id: "7",
    type: "expression",
    data: {
      expression: "input > 0",
      label: "expression 7",
    },
    position: {
      x: 557.8272577093301,
      y: 55.01450900712885,
    },
  },
  {
    id: "8",
    type: "renderer",
    data: {
      label: "renderer 8",
      _executionData: true,
    },
    position: {
      x: 735.5228211229668,
      y: 49.57819290650292,
    },
  },
];

const initialEdges: RFEdge[] = [
  {
    id: "e1-3",
    source: "1",
    target: "3",
    type: "smoothstep",
  },
  {
    id: "e2-3",
    source: "2",
    target: "3",
    type: "smoothstep",
  },
  {
    id: "reactflow__edge-3-5",
    source: "3",
    target: "5",
    type: "smoothstep",
  },
  {
    id: "reactflow__edge-5-6",
    source: "5",
    target: "6",
    type: "smoothstep",
  },
  {
    id: "reactflow__edge-3-7",
    source: "3",
    target: "7",
    type: "smoothstep",
  },
  {
    id: "reactflow__edge-7-8",
    source: "7",
    target: "8",
    type: "smoothstep",
  },
];

// Node categories configuration
const nodeCategories = [
  {
    name: "Input/Output",
    nodes: [
      {
        type: "number",
        label: "Number",
        color: "bg-blue-600",
        defaultData: { value: 0 },
      },
      {
        type: "text_input",
        label: "Text Input",
        color: "bg-green-600",
        defaultData: { text: "" },
      },
      {
        type: "boolean_input",
        label: "Boolean",
        color: "bg-indigo-600",
        defaultData: { value: false },
      },
      {
        type: "date_input",
        label: "Date",
        color: "bg-cyan-600",
        defaultData: { value: "" },
      },
      {
        type: "datetime_input",
        label: "DateTime",
        color: "bg-teal-600",
        defaultData: { value: "" },
      },
      {
        type: "http",
        label: "HTTP",
        color: "bg-purple-600",
        defaultData: { url: "" },
      },
      {
        type: "visualization",
        label: "Visualization",
        color: "bg-indigo-600",
        defaultData: { mode: "text" },
      },
      {
        type: "visualization",
        label: "Bar Chart",
        color: "bg-violet-600",
        defaultData: {
          orientation: "vertical",
          bar_color: "#3b82f6",
          bar_width: "medium",
          show_values: true,
          max_bars: 20,
        },
      },
      {
        type: "renderer",
        label: "Renderer",
        color: "bg-pink-600",
        defaultData: {},
      },
    ],
  },
  {
    name: "Operations",
    nodes: [
      {
        type: "operation",
        label: "Math Op",
        color: "bg-yellow-600",
        defaultData: { op: "add" },
      },
      {
        type: "text_operation",
        label: "Text Op",
        color: "bg-green-600",
        defaultData: { text_op: "uppercase" },
      },
      {
        type: "expression",
        label: "Expression",
        color: "bg-cyan-600",
        defaultData: { expression: "input * 2" },
      },
      {
        type: "transform",
        label: "Transform",
        color: "bg-teal-600",
        defaultData: { transform_type: "to_array" },
      },
      {
        type: "parse",
        label: "Parse",
        color: "bg-purple-600",
        defaultData: { input_type: "AUTO" },
      },
      {
        type: "extract",
        label: "Extract",
        color: "bg-sky-600",
        defaultData: { field: "" },
      },
    ],
  },
  {
    name: "Control Flow",
    nodes: [
      {
        type: "condition",
        label: "Condition",
        color: "bg-amber-600",
        defaultData: { condition: ">0" },
      },
      {
        type: "filter",
        label: "Filter",
        color: "bg-purple-600",
        defaultData: { condition: "item.age > 0" },
      },
      {
        type: "for_each",
        label: "For Each",
        color: "bg-orange-600",
        defaultData: { max_iterations: 1000 },
      },
      {
        type: "while_loop",
        label: "While Loop",
        color: "bg-red-600",
        defaultData: { condition: ">0", max_iterations: 100 },
      },
      {
        type: "switch",
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
        type: "parallel",
        label: "Parallel",
        color: "bg-violet-600",
        defaultData: { max_concurrency: 10 },
      },
      {
        type: "join",
        label: "Join",
        color: "bg-fuchsia-600",
        defaultData: { join_strategy: "all" },
      },
      {
        type: "split",
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
        type: "variable",
        label: "Variable",
        color: "bg-sky-600",
        defaultData: { var_name: "", var_op: "get" },
      },
      {
        type: "accumulator",
        label: "Accumulator",
        color: "bg-blue-500",
        defaultData: { accum_op: "sum" },
      },
      {
        type: "counter",
        label: "Counter",
        color: "bg-indigo-500",
        defaultData: { counter_op: "increment", delta: 1 },
      },
      {
        type: "cache",
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
        type: "retry",
        label: "Retry",
        color: "bg-red-500",
        defaultData: {
          max_attempts: 3,
          backoff_strategy: "exponential",
          initial_delay: "1s",
        },
      },
      {
        type: "try_catch",
        label: "Try-Catch",
        color: "bg-orange-500",
        defaultData: { continue_on_error: true },
      },
      {
        type: "timeout",
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
        type: "delay",
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
        type: "map",
        label: "Map",
        color: "bg-cyan-600",
        defaultData: { expression: "item * 2" },
      },
      {
        type: "reduce",
        label: "Reduce",
        color: "bg-teal-600",
        defaultData: { expression: "acc + item", initial_value: "0" },
      },
      {
        type: "slice",
        label: "Slice",
        color: "bg-emerald-600",
        defaultData: { start: 0, end: -1 },
      },
      {
        type: "sort",
        label: "Sort",
        color: "bg-lime-600",
        defaultData: { field: "", order: "asc" },
      },
      {
        type: "find",
        label: "Find",
        color: "bg-sky-600",
        defaultData: { expression: "item.id == 1" },
      },
      {
        type: "flat_map",
        label: "FlatMap",
        color: "bg-indigo-600",
        defaultData: { expression: "item.values" },
      },
      {
        type: "group_by",
        label: "Group By",
        color: "bg-violet-600",
        defaultData: { key_field: "category" },
      },
      {
        type: "unique",
        label: "Unique",
        color: "bg-fuchsia-600",
        defaultData: { by_field: "" },
      },
      {
        type: "chunk",
        label: "Chunk",
        color: "bg-pink-600",
        defaultData: { size: 3 },
      },
      {
        type: "reverse",
        label: "Reverse",
        color: "bg-rose-600",
        defaultData: {},
      },
      {
        type: "partition",
        label: "Partition",
        color: "bg-orange-600",
        defaultData: { expression: "item > 0" },
      },
      {
        type: "zip",
        label: "Zip",
        color: "bg-yellow-600",
        defaultData: {},
      },
      {
        type: "sample",
        label: "Sample",
        color: "bg-blue-600",
        defaultData: { count: 1 },
      },
      {
        type: "range",
        label: "Range",
        color: "bg-green-600",
        defaultData: { start: 0, end: 10, step: 1 },
      },
      {
        type: "transpose",
        label: "Transpose",
        color: "bg-red-600",
        defaultData: {},
      },
    ],
  },
];

function Canvas() {
  const [nodes, setNodes, onNodesChangeBase] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [showPayload, setShowPayload] = useState(false);
  const [isPaletteOpen, setIsPaletteOpen] = useState(true); // Start with palette open
  const [showExamplesModal, setShowExamplesModal] = useState(false);
  const [workflowTitle, setWorkflowTitle] = useState("Untitled Workflow");
  const [contextMenu, setContextMenu] = useState<{
    x: number;
    y: number;
    nodeId: string;
  } | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{
    nodeId: string;
    nodeName: string;
  } | null>(null);

  // Execution panel state
  const [isExecutionPanelOpen, setIsExecutionPanelOpen] = useState(false);
  const [isExecuting, setIsExecuting] = useState(false);
  const [executionResult, setExecutionResult] =
    useState<ExecutionResult | null>(null);
  const [executionError, setExecutionError] = useState<string | null>(null);
  const [executionDetails, setExecutionDetails] = useState<string | null>(null);
  const [executionPanelHeight, setExecutionPanelHeight] = useState(250);
  const abortControllerRef = useRef<AbortController | null>(null);

  const { project, getNodes } = useReactFlow();

  // Track the next available node ID (never decrements, only increments)
  // Initialize to 9 since initial nodes use IDs 1-8
  const nextNodeId = useRef(9);

  // Wrap onNodesChange - no special handling needed
  const onNodesChange = useCallback(
    (changes: NodeChange[]) => {
      // Apply the changes
      onNodesChangeBase(changes);
    },
    [onNodesChangeBase]
  );

  const onConnect = useCallback(
    (params: Connection) =>
      setEdges((eds: RFEdge[]) =>
        addEdge({ ...params, type: "smoothstep" }, eds)
      ),
    [setEdges]
  );

  const handleNodeContextMenu = useCallback(
    (nodeId: string, x: number, y: number) => {
      setContextMenu({ x, y, nodeId });
    },
    []
  );

  const handleDeleteNodeDirect = useCallback((nodeId: string) => {
    setNodes((nds) => nds.filter((n) => n.id !== nodeId));
    setEdges((eds) =>
      eds.filter(
        (e) => e.source !== nodeId && e.target !== nodeId
      )
    );
  }, [setNodes, setEdges]);

  const payload = useMemo(
    () => ({
      nodes: nodes.map((n) => ({
        id: n.id,
        type: n.type || "",
        data: n.data,
        position: n.position,
      })),
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
      number: withContextMenu(NumberNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      operation: withContextMenu(OperationNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      visualization: withContextMenu(VizNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      renderer: withContextMenu(RendererNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      text_input: withContextMenu(TextInputNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      boolean_input: withContextMenu(BooleanInputNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      date_input: withContextMenu(DateInputNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      datetime_input: withContextMenu(DateTimeInputNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      text_operation: withContextMenu(
        TextOperationNode,
        handleNodeContextMenu,
        () => setIsPaletteOpen(false),
        handleDeleteNodeDirect
      ),
      http: withContextMenu(HttpNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      condition: withContextMenu(ConditionNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      filter: withContextMenu(FilterNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      for_each: withContextMenu(ForEachNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      while_loop: withContextMenu(WhileLoopNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      variable: withContextMenu(VariableNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      extract: withContextMenu(ExtractNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      transform: withContextMenu(TransformNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      parse: withContextMenu(ParseNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      expression: withContextMenu(ExpressionNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      accumulator: withContextMenu(AccumulatorNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      counter: withContextMenu(CounterNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      switch: withContextMenu(SwitchNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      parallel: withContextMenu(ParallelNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      join: withContextMenu(JoinNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      split: withContextMenu(SplitNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      delay: withContextMenu(DelayNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      cache: withContextMenu(CacheNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      retry: withContextMenu(RetryNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      try_catch: withContextMenu(TryCatchNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      timeout: withContextMenu(TimeoutNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      map: withContextMenu(MapNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      reduce: withContextMenu(ReduceNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      slice: withContextMenu(SliceNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      sort: withContextMenu(SortNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      find: withContextMenu(FindNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      flat_map: withContextMenu(FlatMapNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      group_by: withContextMenu(GroupByNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      unique: withContextMenu(UniqueNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      chunk: withContextMenu(ChunkNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      reverse: withContextMenu(ReverseNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      partition: withContextMenu(PartitionNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      zip: withContextMenu(ZipNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      sample: withContextMenu(SampleNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      range: withContextMenu(RangeNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
      transpose: withContextMenu(TransposeNode, handleNodeContextMenu, () =>
        setIsPaletteOpen(false), handleDeleteNodeDirect
      ),
    }),
    [handleNodeContextMenu, handleDeleteNodeDirect]
  );

  // Check for overlapping nodes
  const findNonOverlappingPosition = useCallback(
    (basePosition: XYPosition): XYPosition => {
      const existingNodes = getNodes();
      const nodeWidth = 150;
      const nodeHeight = 80;
      const padding = 20;

      let position = { ...basePosition };
      let attempts = 0;
      const maxAttempts = 50;

      while (attempts < maxAttempts) {
        const overlaps = existingNodes.some((node) => {
          const dx = Math.abs(node.position.x - position.x);
          const dy = Math.abs(node.position.y - position.y);
          return dx < nodeWidth + padding && dy < nodeHeight + padding;
        });

        if (!overlaps) {
          return position;
        }

        // Try offset positions
        position = {
          x: basePosition.x + attempts * 30,
          y: basePosition.y + (attempts % 5) * 25,
        };
        attempts++;
      }

      return position;
    },
    [getNodes]
  );

  const addNode = useCallback(
    (
      type: string,
      defaultData: Record<string, unknown>,
      position?: XYPosition
    ) => {
      // Use incrementing counter for node ID
      const id = String(nextNodeId.current);
      nextNodeId.current += 1;

      let finalPosition: XYPosition;

      if (position) {
        // If position is provided (from drag-and-drop), use it directly
        finalPosition = position;
      } else {
        // Get viewport dimensions (accounting for nav bars: 14px app + 12px workflow + 7px status = 33px total)
        const navHeight = 112; // Total height of both navs + status bar
        const viewportWidth = window.innerWidth;
        const viewportHeight = window.innerHeight - navHeight;

        // Get base position at center of visible viewport
        const basePosition: XYPosition = project
          ? project({ x: viewportWidth / 2 - 75, y: viewportHeight / 2 })
          : { x: 400, y: 200 };

        // Find non-overlapping position
        finalPosition = findNonOverlappingPosition(basePosition);
      }

      const baseData: NodeData = {
        ...defaultData,
        label: `${type} ${id}`,
      };
      const newNode: RFNode<NodeData> = {
        id,
        position: finalPosition,
        data: baseData,
        type,
      };
      setNodes((nds) => nds.concat(newNode));
    },
    [project, setNodes, findNonOverlappingPosition]
  );

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();

      const data = event.dataTransfer.getData("application/reactflow");
      if (!data) return;

      const { type, defaultData } = JSON.parse(data);

      // Get the position where the node was dropped
      const position = project
        ? project({
            x: event.clientX - 75, // offset to center the node
            y: event.clientY - 40,
          })
        : { x: event.clientX, y: event.clientY };

      addNode(type, defaultData, position);
    },
    [project, addNode]
  );

  const handleNewWorkflow = () => {
    // Clear the canvas by resetting nodes and edges
    setNodes([]);
    setEdges([]);
    setWorkflowTitle("Untitled Workflow");
  };

  const handleOpenWorkflow = () => {
    setShowExamplesModal(true);
  };

  const handleSelectExample = (example: WorkflowExample) => {
    // Load the example workflow
    const exampleNodes = example.nodes as RFNode<NodeData>[];
    const exampleEdges = example.edges as RFEdge[];

    // Ensure all nodes have position data
    const nodesWithPositions = exampleNodes.map((node, index) => ({
      ...node,
      position: node.position || { x: 50 + index * 200, y: 50 + index * 100 },
    }));

    // Set the nodes and edges
    setNodes(nodesWithPositions);
    setEdges(exampleEdges.map((edge) => ({ ...edge, type: "smoothstep" })));

    // Update workflow title
    setWorkflowTitle(example.title);

    // Update nextNodeId to be higher than any existing node ID
    const maxId = exampleNodes.reduce((max, node) => {
      const nodeId = parseInt(String(node.id), 10);
      return isNaN(nodeId) ? max : Math.max(max, nodeId);
    }, 0);
    nextNodeId.current = maxId + 1;
  };

  const handleSave = () => {
    // TODO: Save workflow
    console.log("Save workflow", payload);
  };

  const handleDelete = () => {
    // TODO: Delete workflow
  };

  const handleRun = async () => {
    // Open execution panel
    setIsExecutionPanelOpen(true);
    setIsExecuting(true);
    setExecutionResult(null);
    setExecutionError(null);
    setExecutionDetails(null);

    // Create abort controller for cancellation
    abortControllerRef.current = new AbortController();

    try {
      const response = await fetch("/api/v1/workflow/execute", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
        signal: abortControllerRef.current.signal,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || `HTTP error! status: ${response.status}`);
      }

      setExecutionResult(data);
    } catch (error: unknown) {
      if (
        (error as { name: string; message: string; detail: string }).name ===
        "AbortError"
      ) {
        setExecutionError("Execution cancelled by user");
      } else {
        setExecutionError(
          (error as { name: string; message: string; detail: string }).message
        );
        setExecutionDetails(
          (error as { name: string; message: string; detail: string }).detail
        );
      }
    } finally {
      setIsExecuting(false);
      abortControllerRef.current = null;
    }
  };

  const handleCancelExecution = () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
  };

  const handleCloseExecutionPanel = () => {
    setIsExecutionPanelOpen(false);
    setExecutionResult(null);
    setExecutionError(null);
    setExecutionDetails(null);
  };

  const handleExport = () => {
    const jsonString = JSON.stringify(payload, null, 2);
    const blob = new Blob([jsonString], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;

    // Use workflow title for filename, sanitize it
    const sanitizedTitle = workflowTitle
      .replace(/[^a-z0-9]/gi, "_")
      .toLowerCase();
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
      position: node.position || { x: 50 + index * 200, y: 50 + index * 100 },
    }));

    // Set the imported nodes and edges
    setNodes(nodesWithPositions);
    setEdges(importedEdges.map((edge) => ({ ...edge, type: "smoothstep" })));

    // Update nextNodeId to be higher than any existing node ID
    const maxId = importedNodes.reduce((max, node) => {
      const nodeId = parseInt(String(node.id), 10);
      return isNaN(nodeId) ? max : Math.max(max, nodeId);
    }, 0);
    nextNodeId.current = maxId + 1;
  };

  const handleDeleteNode = (nodeId: string) => {
    const node = nodes.find((n) => n.id === nodeId);
    if (node) {
      setDeleteConfirm({
        nodeId,
        nodeName: String(node.data?.label || `Node ${nodeId}`),
      });
    }
    setContextMenu(null);
  };

  const confirmDelete = () => {
    if (deleteConfirm) {
      setNodes((nds) => nds.filter((n) => n.id !== deleteConfirm.nodeId));
      setEdges((eds) =>
        eds.filter(
          (e) =>
            e.source !== deleteConfirm.nodeId &&
            e.target !== deleteConfirm.nodeId
        )
      );
    }
    setDeleteConfirm(null);
  };

  // Update renderer nodes with execution data when execution completes
  useEffect(() => {
    if (executionResult?.success && executionResult.results?.node_results) {
      const nodeResults = executionResult.results.node_results;

      setNodes((nds) =>
        nds.map((node) => {
          // Only update renderer nodes
          if (node.type === "renderer") {
            // Get the execution data for this node
            const executionData = nodeResults[node.id];

            return {
              ...node,
              data: {
                ...node.data,
                _executionData: executionData,
              },
            };
          }
          return node;
        })
      );
    }
  }, [executionResult, setNodes]);

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

      {/* Main Content - Flex container for sidebar and canvas */}
      <div className="flex-1 flex overflow-hidden">
        {/* Node Palette Sidebar */}
        {isPaletteOpen && (
          <NodePalette
            isOpen={isPaletteOpen}
            onClose={() => setIsPaletteOpen(false)}
            categories={nodeCategories}
            onAddNode={addNode}
          />
        )}

        {/* Canvas Container */}
        <div className="flex-1 relative">
          {/* Toggle Sidebar Button - Top Left of Canvas */}
          {isPaletteOpen ? (
            <> </>
          ) : (
            <>
              <button
                onClick={() => setIsPaletteOpen(!isPaletteOpen)}
                className="absolute left-4 top-4 z-10 bg-gray-800 hover:bg-gray-700 text-white px-3 py-1.5 rounded-lg shadow-lg transition-all border border-gray-700 hover:border-gray-600 text-sm font-medium flex items-center gap-1.5"
                title={"Show Nodes Panel"}
              >
                <>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={2}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                    />
                  </svg>
                  <span>Show Nodes</span>
                </>
              </button>
            </>
          )}

          {/* React Flow Canvas */}
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onDragOver={onDragOver}
            onDrop={onDrop}
            nodeTypes={nodeTypes}
            proOptions={{ hideAttribution: true }}
            fitView
            className="bg-gray-950"
          />
        </div>
      </div>

      {/* Modals */}
      <JSONPayloadModal
        isOpen={showPayload}
        onClose={() => setShowPayload(false)}
        payload={payload}
        workflowTitle={workflowTitle}
      />
      <WorkflowExamplesModal
        isOpen={showExamplesModal}
        onClose={() => setShowExamplesModal(false)}
        onSelect={handleSelectExample}
      />

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

      {/* Execution Results Panel */}
      <ExecutionPanel
        isOpen={isExecutionPanelOpen}
        isLoading={isExecuting}
        result={executionResult}
        error={executionError}
        details={executionDetails}
        onCancel={handleCancelExecution}
        onClose={handleCloseExecutionPanel}
        height={executionPanelHeight}
        onHeightChange={setExecutionPanelHeight}
      />

      {/* Bottom Status Bar - Only show when execution panel is closed */}
      {!isExecutionPanelOpen && (
        <WorkflowStatusBar nodeCount={nodes.length} edgeCount={edges.length} />
      )}
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
