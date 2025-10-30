"use client";
import React, { useCallback, useMemo, useState, useRef, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
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
  BackgroundVariant,
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
  ContextVariableNode,
  ContextConstantNode,
  DeleteConfirmDialog,
  NodeWrapper,
  getNodeInfo,
} from "../../../components/nodes";
import { AppNavBar } from "../../../components/AppNavBar";
import { WorkflowNavBar } from "../../../components/WorkflowNavBar";
import { JSONPayloadModal } from "../../../components/JSONPayloadModal";
import { OpenWorkflowModal } from "../../../components/OpenWorkflowModal";
import { Toast } from "../../../components/Toast";
import { Workflow } from "../../../types/workflow";

type NodeData = Record<string, unknown>;

type ExtendedNodeProps = NodeProps<NodeData> & {
  onShowOptions?: (x: number, y: number) => void;
};

const BASE_NODE_CLASSES = "bg-gradient-to-br from-gray-700 to-gray-800 text-white shadow-lg rounded-lg border border-gray-600 hover:border-gray-500 transition-all";

function NodeContextMenu({
  x,
  y,
  onClose,
  onRename,
  onDuplicate,
  onDelete,
}: {
  x: number;
  y: number;
  onClose: () => void;
  onRename: () => void;
  onDuplicate: () => void;
  onDelete: () => void;
}) {
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onClose();
      }
    };
    
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };
    
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [onClose]);

  return (
    <div
      ref={menuRef}
      className="fixed bg-gray-800 border border-gray-700 rounded-lg shadow-xl py-1 z-50 min-w-[160px]"
      style={{ left: x, top: y }}
      role="menu"
      aria-label="Node options"
    >
      <button
        onClick={onRename}
        className="w-full px-4 py-2 text-left text-sm text-gray-200 hover:bg-gray-700 transition-colors flex items-center gap-2"
        role="menuitem"
      >
        <span>✏️</span> Rename
      </button>
      <button
        onClick={onDuplicate}
        className="w-full px-4 py-2 text-left text-sm text-gray-200 hover:bg-gray-700 transition-colors flex items-center gap-2"
        role="menuitem"
      >
        <span>📋</span> Duplicate
      </button>
      <div className="border-t border-gray-700 my-1"></div>
      <button
        onClick={onDelete}
        className="w-full px-4 py-2 text-left text-sm text-red-400 hover:bg-gray-700 transition-colors flex items-center gap-2"
        role="menuitem"
      >
        <span>🗑️</span> Delete
      </button>
    </div>
  );
}

function createCompactNode(
  Component: React.ComponentType<ExtendedNodeProps>,
  showMenu: (id: string, x: number, y: number) => void
) {
  return function CompactNodeWrapper(props: NodeProps<NodeData>) {
    const handleContextMenu = (e: React.MouseEvent) => {
      e.preventDefault();
      showMenu(props.id, e.clientX, e.clientY);
    };

    const handleShowOptions = (x: number, y: number) => {
      showMenu(props.id, x, y);
    };

    const extendedProps: ExtendedNodeProps = {
      ...props,
      onShowOptions: handleShowOptions
    };

    return (
      <div onContextMenu={handleContextMenu} className="compact-node-wrapper">
        <Component {...extendedProps} />
      </div>
    );
  };
}

function NumberNode({ id, data, ...props }: ExtendedNodeProps) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const v = Number(e.target.value);
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, value: v } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("numberNode");
  const onShowOptions = props.onShowOptions;

  return (
    <NodeWrapper
      title={String(data?.label || "Number")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className={BASE_NODE_CLASSES}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <input
        value={Number(data?.value ?? 0)}
        type="number"
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
        aria-label="Number value"
      />
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

function OperationNode({ id, data, ...props }: ExtendedNodeProps) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const op = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, op } } : n))
    );
  };

  const nodeInfo = getNodeInfo("opNode");
  const onShowOptions = props.onShowOptions;

  return (
    <NodeWrapper
      title={String(data?.label || "Operation")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className={BASE_NODE_CLASSES}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <select
        value={String(data?.op ?? "add")}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
        aria-label="Operation type"
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

function VizNode({ id, data, ...props }: ExtendedNodeProps) {
  const { setNodes } = useReactFlow();
  const onChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const mode = e.target.value;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, mode } } : n))
    );
  };

  const nodeInfo = getNodeInfo("vizNode");
  const onShowOptions = props.onShowOptions;

  return (
    <NodeWrapper
      title={String(data?.label || "Visualization")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className={BASE_NODE_CLASSES}
    >
      <Handle type="target" position={Position.Left} className="w-2 h-2 bg-blue-400" />
      <select
        value={String(data?.mode ?? "text")}
        onChange={onChange}
        className="w-24 text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
        aria-label="Visualization mode"
      >
        <option value="text">Text</option>
        <option value="table">Table</option>
      </select>
      <Handle type="source" position={Position.Right} className="w-2 h-2 bg-green-400" />
    </NodeWrapper>
  );
}

const nodeCategories = [
  {
    name: "Input / Output",
    icon: "📥",
    nodes: [
      { type: "numberNode", label: "Number", icon: "🔢", defaultData: { value: 0 } },
      { type: "textInputNode", label: "Text Input", icon: "📝", defaultData: { text: "" } },
      { type: "httpNode", label: "HTTP Request", icon: "🌐", defaultData: { url: "" } },
      { type: "vizNode", label: "Visualization", icon: "📊", defaultData: { mode: "text" } },
    ],
  },
  {
    name: "Operations",
    icon: "⚙️",
    nodes: [
      { type: "opNode", label: "Math Operation", icon: "➕", defaultData: { op: "add" } },
      { type: "textOpNode", label: "Text Operation", icon: "✂️", defaultData: { text_op: "uppercase" } },
      { type: "transformNode", label: "Transform", icon: "🔄", defaultData: { transform_type: "to_array" } },
      { type: "extractNode", label: "Extract", icon: "🔍", defaultData: { field: "" } },
    ],
  },
  {
    name: "Control Flow",
    icon: "🔀",
    nodes: [
      { type: "conditionNode", label: "Condition", icon: "❓", defaultData: { condition: ">0" } },
      { type: "forEachNode", label: "For Each", icon: "🔁", defaultData: { max_iterations: 1000 } },
      { type: "whileLoopNode", label: "While Loop", icon: "🔂", defaultData: { condition: ">0", max_iterations: 100 } },
      { type: "switchNode", label: "Switch", icon: "🔀", defaultData: { cases: [], default_path: "default" } },
    ],
  },
  {
    name: "Parallel & Join",
    icon: "⚡",
    nodes: [
      { type: "parallelNode", label: "Parallel", icon: "⚡", defaultData: { max_concurrency: 10 } },
      { type: "joinNode", label: "Join", icon: "🔗", defaultData: { join_strategy: "all" } },
      { type: "splitNode", label: "Split", icon: "🔱", defaultData: { paths: ["path1", "path2"] } },
    ],
  },
  {
    name: "State & Memory",
    icon: "💾",
    nodes: [
      { type: "variableNode", label: "Variable", icon: "📦", defaultData: { var_name: "", var_op: "get" } },
      { type: "cacheNode", label: "Cache", icon: "💾", defaultData: { cache_op: "get", cache_key: "" } },
      { type: "accumulatorNode", label: "Accumulator", icon: "📈", defaultData: { accum_op: "sum" } },
      { type: "counterNode", label: "Counter", icon: "🔢", defaultData: { counter_op: "increment", delta: 1 } },
    ],
  },
  {
    name: "Error Handling",
    icon: "🛡️",
    nodes: [
      { type: "retryNode", label: "Retry", icon: "🔄", defaultData: { max_attempts: 3, backoff_strategy: "exponential", initial_delay: "1s" } },
      { type: "tryCatchNode", label: "Try-Catch", icon: "🛡️", defaultData: { continue_on_error: true } },
      { type: "timeoutNode", label: "Timeout", icon: "⏱️", defaultData: { timeout: "30s", timeout_action: "error" } },
    ],
  },
  {
    name: "Utilities",
    icon: "🔧",
    nodes: [
      { type: "delayNode", label: "Delay", icon: "⏸️", defaultData: { duration: "1s" } },
    ],
  },
  {
    name: "Context",
    icon: "🎯",
    nodes: [
      { type: "contextVariableNode", label: "Variable", icon: "📦", defaultData: { context_name: "", context_value: "" } },
      { type: "contextConstantNode", label: "Constant", icon: "🔒", defaultData: { context_name: "", context_value: "" } },
    ],
  },
];

function WorkflowEditor() {
  const params = useParams();
  const router = useRouter();
  const workflowId = params?.id as string;
  
  const [nodes, setNodes, onNodesChange] = useNodesState<NodeData>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [expandedCategories, setExpandedCategories] = useState<Record<string, boolean>>({
    "Input / Output": true,
  });
  const [contextMenu, setContextMenu] = useState<{
    nodeId: string;
    x: number;
    y: number;
  } | null>(null);
  const [renamingNode, setRenamingNode] = useState<string | null>(null);
  const [renameValue, setRenameValue] = useState("");
  const [deletingNode, setDeletingNode] = useState<string | null>(null);
  
  const [workflowTitle, setWorkflowTitle] = useState("Untitled Workflow");
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);
  const [showJSONModal, setShowJSONModal] = useState(false);
  const [showOpenModal, setShowOpenModal] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" | "info" | "warning" } | null>(null);

  const { project } = useReactFlow();

  // Load workflow on mount
  useEffect(() => {
    if (workflowId) {
      const saved = localStorage.getItem("thaiyyal_workflows");
      if (saved) {
        const workflows: Workflow[] = JSON.parse(saved);
        const workflow = workflows.find((w) => w.id === workflowId);
        if (workflow) {
          setNodes(workflow.data.nodes as RFNode<NodeData>[]);
          setEdges(workflow.data.edges as RFEdge[]);
          // Use a separate effect or initial state
          if (workflowTitle === "Untitled Workflow") {
            setWorkflowTitle(workflow.title);
          }
          if (hasUnsavedChanges) {
            setHasUnsavedChanges(false);
          }
        }
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [workflowId, setNodes, setEdges]);

  const onConnect = useCallback(
    (params: Connection) => {
      setEdges((eds: RFEdge[]) => addEdge(params, eds));
      setHasUnsavedChanges(true);
    },
    [setEdges]
  );

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

  const showContextMenu = useCallback((nodeId: string, x: number, y: number) => {
    setContextMenu({ nodeId, x, y });
  }, []);

  const nodeTypes = useMemo(
    () => ({
      numberNode: createCompactNode(NumberNode, showContextMenu),
      opNode: createCompactNode(OperationNode, showContextMenu),
      vizNode: createCompactNode(VizNode, showContextMenu),
      textInputNode: createCompactNode(TextInputNode, showContextMenu),
      textOpNode: createCompactNode(TextOperationNode, showContextMenu),
      httpNode: createCompactNode(HttpNode, showContextMenu),
      conditionNode: createCompactNode(ConditionNode, showContextMenu),
      forEachNode: createCompactNode(ForEachNode, showContextMenu),
      whileLoopNode: createCompactNode(WhileLoopNode, showContextMenu),
      variableNode: createCompactNode(VariableNode, showContextMenu),
      extractNode: createCompactNode(ExtractNode, showContextMenu),
      transformNode: createCompactNode(TransformNode, showContextMenu),
      accumulatorNode: createCompactNode(AccumulatorNode, showContextMenu),
      counterNode: createCompactNode(CounterNode, showContextMenu),
      switchNode: createCompactNode(SwitchNode, showContextMenu),
      parallelNode: createCompactNode(ParallelNode, showContextMenu),
      joinNode: createCompactNode(JoinNode, showContextMenu),
      splitNode: createCompactNode(SplitNode, showContextMenu),
      delayNode: createCompactNode(DelayNode, showContextMenu),
      cacheNode: createCompactNode(CacheNode, showContextMenu),
      retryNode: createCompactNode(RetryNode, showContextMenu),
      tryCatchNode: createCompactNode(TryCatchNode, showContextMenu),
      timeoutNode: createCompactNode(TimeoutNode, showContextMenu),
      contextVariableNode: createCompactNode(ContextVariableNode, showContextMenu),
      contextConstantNode: createCompactNode(ContextConstantNode, showContextMenu),
    }),
    [showContextMenu]
  );

  const [nextId, setNextId] = useState(1);
  
  const addNode = useCallback((type: string, defaultData: Record<string, unknown>, label: string) => {
    const id = String(nextId);
    setNextId((s) => s + 1);
    const position: XYPosition = project
      ? project({ x: 200, y: 100 })
      : { x: 400 + nextId * 10, y: 120 + (nextId % 3) * 40 };
    
    const baseData: NodeData = { ...defaultData, label };
    const newNode: RFNode<NodeData> = { id, position, data: baseData, type };
    setNodes((nds) => nds.concat(newNode));
    setHasUnsavedChanges(true);
  }, [nextId, project, setNodes]);

  const handleRename = useCallback(() => {
    if (!contextMenu) return;
    const node = nodes.find((n) => n.id === contextMenu.nodeId);
    if (node) {
      setRenamingNode(contextMenu.nodeId);
      setRenameValue(String(node.data?.label || ""));
    }
    setContextMenu(null);
  }, [contextMenu, nodes]);

  const handleDuplicate = useCallback(() => {
    if (!contextMenu) return;
    const node = nodes.find((n) => n.id === contextMenu.nodeId);
    if (node) {
      const id = String(nextId);
      setNextId((s) => s + 1);
      const newNode: RFNode<NodeData> = {
        ...node,
        id,
        position: { x: node.position.x + 50, y: node.position.y + 50 },
        data: { ...node.data, label: `${node.data?.label || node.type} (copy)` },
      };
      setNodes((nds) => nds.concat(newNode));
      setHasUnsavedChanges(true);
    }
    setContextMenu(null);
  }, [contextMenu, nodes, nextId, setNodes]);

  const handleDelete = useCallback(() => {
    if (!contextMenu) return;
    setDeletingNode(contextMenu.nodeId);
    setContextMenu(null);
  }, [contextMenu]);

  const confirmDelete = useCallback(() => {
    if (!deletingNode) return;
    setNodes((nds) => nds.filter((n) => n.id !== deletingNode));
    setEdges((eds) => eds.filter((e) => e.source !== deletingNode && e.target !== deletingNode));
    setDeletingNode(null);
    setHasUnsavedChanges(true);
  }, [deletingNode, setNodes, setEdges]);

  const cancelDelete = useCallback(() => {
    setDeletingNode(null);
  }, []);

  const handleRenameSubmit = useCallback(() => {
    if (renamingNode) {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === renamingNode ? { ...n, data: { ...n.data, label: renameValue } } : n
        )
      );
      setRenamingNode(null);
      setRenameValue("");
      setHasUnsavedChanges(true);
    }
  }, [renamingNode, renameValue, setNodes]);

  const toggleCategory = useCallback((categoryName: string) => {
    setExpandedCategories((prev) => ({
      ...prev,
      [categoryName]: !prev[categoryName],
    }));
  }, []);

  const filteredCategories = useMemo(() => {
    if (!searchTerm) return nodeCategories;
    
    const lowercaseSearch = searchTerm.toLowerCase();
    return nodeCategories
      .map((category) => ({
        ...category,
        nodes: category.nodes.filter((node) =>
          node.label.toLowerCase().includes(lowercaseSearch)
        ),
      }))
      .filter((category) => category.nodes.length > 0);
  }, [searchTerm]);

  const handleNewWorkflow = useCallback(() => {
    if (hasUnsavedChanges) {
      if (!confirm("You have unsaved changes. Create new workflow anyway?")) {
        return;
      }
    }
    const newWorkflowId = `workflow-${Date.now()}`;
    router.push(`/workflow/${newWorkflowId}`);
  }, [hasUnsavedChanges, router]);

  const handleSaveWorkflow = useCallback(() => {
    const workflow: Workflow = {
      id: workflowId,
      title: workflowTitle,
      data: { nodes, edges },
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    const saved = localStorage.getItem("thaiyyal_workflows");
    const workflows: Workflow[] = saved ? JSON.parse(saved) : [];
    
    const existingIndex = workflows.findIndex((w) => w.id === workflowId);
    if (existingIndex >= 0) {
      workflow.createdAt = workflows[existingIndex].createdAt;
      workflows[existingIndex] = workflow;
    } else {
      workflows.push(workflow);
    }
    
    localStorage.setItem("thaiyyal_workflows", JSON.stringify(workflows));
    setHasUnsavedChanges(false);
    setToast({ message: "Workflow saved successfully", type: "success" });
  }, [workflowId, workflowTitle, nodes, edges]);

  const handleOpenWorkflow = useCallback((workflow: Workflow) => {
    router.push(`/workflow/${workflow.id}`);
  }, [router]);

  const handleDeleteWorkflow = useCallback(() => {
    if (!confirm(`Are you sure you want to delete "${workflowTitle}"?`)) {
      return;
    }
    
    const saved = localStorage.getItem("thaiyyal_workflows");
    if (saved) {
      const workflows: Workflow[] = JSON.parse(saved);
      const updated = workflows.filter((w) => w.id !== workflowId);
      localStorage.setItem("thaiyyal_workflows", JSON.stringify(updated));
    }
    
    setToast({ message: "Workflow deleted", type: "success" });
    router.push("/");
  }, [workflowId, workflowTitle, router]);

  const handleRunWorkflow = useCallback(() => {
    setToast({ message: "Workflow execution started! (Mock)", type: "info" });
  }, []);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault();
        handleSaveWorkflow();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleSaveWorkflow]);

  useEffect(() => {
    if (nodes.length > 0 || edges.length > 0) {
      const timer = setTimeout(() => {
        // Mark as changed after initial load
      }, 100);
      return () => clearTimeout(timer);
    }
  }, [nodes, edges]);

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      <AppNavBar
        onNewWorkflow={handleNewWorkflow}
        onOpenWorkflow={() => setShowOpenModal(true)}
      />

      <WorkflowNavBar
        workflowTitle={workflowTitle}
        onTitleChange={(title) => {
          setWorkflowTitle(title);
          setHasUnsavedChanges(true);
        }}
        onSave={handleSaveWorkflow}
        onShowJSON={() => setShowJSONModal(true)}
        onDelete={handleDeleteWorkflow}
        onRun={handleRunWorkflow}
        hasUnsavedChanges={hasUnsavedChanges}
      />

      <div className="flex-1 flex overflow-hidden relative">
        <div 
          className="absolute left-4 top-4 z-10 w-64 bg-gray-900/95 backdrop-blur-sm border border-gray-700 rounded-xl shadow-2xl overflow-hidden"
          role="toolbar"
          aria-label="Node palette"
        >
          <div className="p-3 border-b border-gray-700">
            <input
              type="text"
              placeholder="🔍 Search nodes..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full px-3 py-2 text-sm bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              aria-label="Search nodes"
            />
          </div>

          <div className="max-h-[calc(100vh-250px)] overflow-y-auto custom-scrollbar">
            {filteredCategories.map((category) => (
              <div key={category.name} className="border-b border-gray-800 last:border-b-0">
                <button
                  onClick={() => toggleCategory(category.name)}
                  className="w-full px-3 py-2 flex items-center justify-between text-left text-sm font-semibold text-gray-300 hover:bg-gray-800 transition-colors"
                  aria-expanded={expandedCategories[category.name]}
                  aria-controls={`category-${category.name}`}
                >
                  <span className="flex items-center gap-2">
                    <span>{category.icon}</span>
                    <span>{category.name}</span>
                  </span>
                  <span className="text-gray-500 text-xs">
                    {expandedCategories[category.name] ? "▼" : "▶"}
                  </span>
                </button>
                
                {expandedCategories[category.name] && (
                  <div id={`category-${category.name}`} className="bg-gray-900/50">
                    {category.nodes.map((node) => (
                      <button
                        key={node.type}
                        onClick={() => addNode(node.type, node.defaultData, node.label)}
                        className="w-full px-4 py-2 text-left text-sm text-gray-300 hover:bg-gray-800 transition-colors flex items-center gap-2 group"
                        aria-label={`Add ${node.label} node`}
                      >
                        <span className="text-base">{node.icon}</span>
                        <span className="flex-1">{node.label}</span>
                        <span className="text-xs text-gray-600 group-hover:text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity">+</span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
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
            className="bg-gray-950"
            defaultEdgeOptions={{
              type: 'smoothstep',
              animated: true,
              style: { stroke: '#6b7280', strokeWidth: 2 },
            }}
          >
            <Background 
              variant={BackgroundVariant.Dots} 
              gap={16} 
              size={1} 
              color="#374151"
            />
          </ReactFlow>

          <div className="absolute bottom-0 left-0 right-0 px-4 py-3 border-t border-gray-800 bg-gray-900/95 backdrop-blur-sm flex items-center justify-between">
            <div className="text-sm text-gray-400 flex items-center gap-4">
              <span className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                {nodes.length} nodes
              </span>
              <span className="flex items-center gap-2">
                <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                {edges.length} connections
              </span>
            </div>
          </div>
        </div>
      </div>

      <JSONPayloadModal
        isOpen={showJSONModal}
        onClose={() => setShowJSONModal(false)}
        payload={payload}
      />

      <OpenWorkflowModal
        isOpen={showOpenModal}
        onClose={() => setShowOpenModal(false)}
        onSelect={handleOpenWorkflow}
      />

      {contextMenu && (
        <NodeContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onClose={() => setContextMenu(null)}
          onRename={handleRename}
          onDuplicate={handleDuplicate}
          onDelete={handleDelete}
        />
      )}

      {renamingNode && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 w-96 shadow-2xl">
            <h3 className="text-lg font-semibold text-gray-200 mb-4">Rename Node</h3>
            <input
              type="text"
              value={renameValue}
              onChange={(e) => setRenameValue(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleRenameSubmit();
                if (e.key === "Escape") setRenamingNode(null);
              }}
              className="w-full px-3 py-2 bg-gray-900 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              autoFocus
              aria-label="Node name"
            />
            <div className="flex gap-2 mt-4">
              <button
                onClick={handleRenameSubmit}
                className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium"
              >
                Save
              </button>
              <button
                onClick={() => setRenamingNode(null)}
                className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors font-medium"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {deletingNode && (
        <DeleteConfirmDialog
          nodeName={String(nodes.find((n) => n.id === deletingNode)?.data?.label || deletingNode)}
          onConfirm={confirmDelete}
          onCancel={cancelDelete}
        />
      )}

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          isVisible={!!toast}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}

export default function WorkflowPage() {
  return (
    <ReactFlowProvider>
      <WorkflowEditor />
    </ReactFlowProvider>
  );
}
