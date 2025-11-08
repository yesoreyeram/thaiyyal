"use client";
import React, {
  ComponentType,
  useCallback,
  useMemo,
  useState,
  useRef,
  useEffect,
} from "react";
import ReactFlow, {
  addEdge,
  ReactFlowProvider,
  Node as RFNode,
  Edge as RFEdge,
  XYPosition,
  Connection,
  NodeProps,
  NodeChange,
  useEdgesState,
  useNodesState,
  useReactFlow,
} from "reactflow";
import "reactflow/dist/style.css";
import * as Nodes from "../../components/nodes";
import {
  HomeIcon,
  FileIcon,
  PlusIcon,
  PlayIcon,
  SaveIcon,
  DownloadIcon,
  UploadIcon,
  SettingsIcon,
  FileJsonIcon,
} from "lucide-react";
import {
  SidebarProvider,
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenuItem,
  SidebarToggle,
  useSidebar,
} from "../../components/ui/sidebar";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
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

type NodePropsWithOptions = Nodes.NodePropsWithOptions;

type NodeComponentMap = Record<string, ComponentType<NodePropsWithOptions>>;

const withContextMenu = (
  Component: ComponentType<Nodes.NodePropsWithOptions>,
  handleContextMenu: (nodeId: string, x: number, y: number) => void,
  handleDelete: (nodeId: string) => void
) => {
  const WrappedComponent = (props: NodeProps<NodeData>) => {
    const onShowOptions = (x: number, y: number) => {
      handleContextMenu(props.id, x, y);
    };
    return (
      <Component
        {...(props as Nodes.NodePropsWithOptions)}
        onShowOptions={onShowOptions}
        onOpenInfo={() => {}}
        onDelete={() => handleDelete(props.id)}
      />
    );
  };
  WrappedComponent.displayName = `withContextMenu(${
    Component.displayName || Component.name || "Component"
  })`;
  return WrappedComponent;
};

function Canvas() {
  const [nodes, setNodes, onNodesChangeBase] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
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
  const [isExecuting, setIsExecuting] = useState(false);
  const [executionResult, setExecutionResult] =
    useState<ExecutionResult | null>(null);
  const [executionError, setExecutionError] = useState<string | null>(null);
  const [executionDetails, setExecutionDetails] = useState<string | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const { project, getNodes } = useReactFlow();
  const nextNodeId = useRef(1);
  const onNodesChange = useCallback(
    (changes: NodeChange[]) => onNodesChangeBase(changes),
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

  const handleDeleteNodeDirect = useCallback(
    (nodeId: string) => {
      setNodes((nds) => nds.filter((n) => n.id !== nodeId));
      setEdges((eds) =>
        eds.filter((e) => e.source !== nodeId && e.target !== nodeId)
      );
    },
    [setNodes, setEdges]
  );

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
        sourceHandle: e.sourceHandle,
        targetHandle: e.targetHandle,
      })),
    }),
    [nodes, edges]
  );

  const nodeTypes = useMemo(() => {
    const nodeComponentMap: NodeComponentMap = {
      number: Nodes.NumberNode,
      operation: Nodes.OperationNode,
      visualization: Nodes.VizNode,
      renderer: Nodes.RendererNode,
      text_input: Nodes.TextInputNode,
      boolean_input: Nodes.BooleanInputNode,
      date_input: Nodes.DateInputNode,
      datetime_input: Nodes.DateTimeInputNode,
      text_operation: Nodes.TextOperationNode,
      http: Nodes.HttpNode,
      condition: Nodes.ConditionNode,
      filter: Nodes.FilterNode,
      for_each: Nodes.ForEachNode,
      while_loop: Nodes.WhileLoopNode,
      variable: Nodes.VariableNode,
      extract: Nodes.ExtractNode,
      transform: Nodes.TransformNode,
      parse: Nodes.ParseNode,
      expression: Nodes.ExpressionNode,
      accumulator: Nodes.AccumulatorNode,
      counter: Nodes.CounterNode,
      switch: Nodes.SwitchNode,
      parallel: Nodes.ParallelNode,
      join: Nodes.JoinNode,
      split: Nodes.SplitNode,
      delay: Nodes.DelayNode,
      cache: Nodes.CacheNode,
      retry: Nodes.RetryNode,
      try_catch: Nodes.TryCatchNode,
      timeout: Nodes.TimeoutNode,
      map: Nodes.MapNode,
      reduce: Nodes.ReduceNode,
      slice: Nodes.SliceNode,
      sort: Nodes.SortNode,
      find: Nodes.FindNode,
      flat_map: Nodes.FlatMapNode,
      group_by: Nodes.GroupByNode,
      unique: Nodes.UniqueNode,
      chunk: Nodes.ChunkNode,
      reverse: Nodes.ReverseNode,
      partition: Nodes.PartitionNode,
      zip: Nodes.ZipNode,
      sample: Nodes.SampleNode,
      range: Nodes.RangeNode,
      transpose: Nodes.TransposeNode,
    };

    // Wrap all node components with context menu handlers
    return Object.entries(nodeComponentMap).reduce(
      (acc, [key, Component]) => ({
        ...acc,
        [key]: withContextMenu(
          Component,
          handleNodeContextMenu,
          // () => setIsPaletteOpen(false),
          handleDeleteNodeDirect
        ),
      }),
      {} as Record<string, React.ComponentType<NodeProps<NodeData>>>
    );
  }, [handleNodeContextMenu, handleDeleteNodeDirect]);

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

  const handleOpenWorkflow = () => setShowExamplesModal(true);

  const handleSaveWorkflow = () => {
    // Placeholder for save functionality
    // In the future, this could save to localStorage or backend
    console.log("Saving workflow:", workflowTitle, payload);
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

  const handleRun = async () => {
    // Panel is always visible, just update execution state
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
    // Clear execution results (panel stays visible)
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

  const handleImportClick = () => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = ".json";
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (event) => {
          try {
            const jsonData = JSON.parse(event.target?.result as string);
            handleImport(jsonData);
          } catch (error) {
            console.error("Error parsing JSON:", error);
          }
        };
        reader.readAsText(file);
      }
    };
    input.click();
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

  const { isCollapsed } = useSidebar();
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [titleInput, setTitleInput] = useState(workflowTitle);

  const handleSaveTitle = () => {
    setWorkflowTitle(titleInput);
    setIsEditingTitle(false);
  };

  const handleCancelEditTitle = () => {
    setTitleInput(workflowTitle);
    setIsEditingTitle(false);
  };

  return (
    <div className="h-screen flex bg-gray-950">
      {/* Sidebar Navigation */}
      <Sidebar>
        <SidebarHeader>
          {!isCollapsed && (
            <span className="text-sm font-semibold text-black dark:text-white">
              ⚡ Thaiyyal
            </span>
          )}
          {isCollapsed && (
            <span className="text-sm font-semibold text-black dark:text-white">
              ⚡
            </span>
          )}
        </SidebarHeader>

        <SidebarContent>
          {/* Navigation Group */}
          <SidebarGroup>
            {!isCollapsed && (
              <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            )}
            <SidebarGroupContent>
              <SidebarMenuItem href="/" icon={<HomeIcon className="w-4 h-4" />} label="Home" />
              <SidebarMenuItem href="/playground" icon={<FileJsonIcon className="w-4 h-4" />} label="Playground" />
              <SidebarMenuItem href="/workflow" icon={<FileIcon className="w-4 h-4" />} label="Workflows" active />
            </SidebarGroupContent>
          </SidebarGroup>

          {/* Actions Group */}
          <SidebarGroup>
            {!isCollapsed && <SidebarGroupLabel>Actions</SidebarGroupLabel>}
            <SidebarGroupContent>
              <SidebarMenuItem
                onClick={handleNewWorkflow}
                icon={<PlusIcon className="w-4 h-4" />}
                label="New Workflow"
              />
              <SidebarMenuItem
                onClick={handleOpenWorkflow}
                icon={<UploadIcon className="w-4 h-4" />}
                label="Open Workflow"
              />
              <SidebarMenuItem
                onClick={handleRun}
                icon={<PlayIcon className="w-4 h-4" />}
                label="Run Workflow"
              />
              <SidebarMenuItem
                onClick={() => setShowPayload(true)}
                icon={<FileJsonIcon className="w-4 h-4" />}
                label="Show JSON"
              />
              <SidebarMenuItem
                onClick={handleSaveWorkflow}
                icon={<SaveIcon className="w-4 h-4" />}
                label="Save"
              />
              <SidebarMenuItem
                onClick={handleExport}
                icon={<DownloadIcon className="w-4 h-4" />}
                label="Export"
              />
              <SidebarMenuItem
                onClick={handleImportClick}
                icon={<UploadIcon className="w-4 h-4" />}
                label="Import"
              />
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <SidebarMenuItem
            icon={<SettingsIcon className="w-4 h-4" />}
            label="Settings"
          />
        </SidebarFooter>

        <SidebarToggle />
      </Sidebar>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Workflow Title Bar */}
        <div className="flex items-center justify-between h-12 px-4 border-b border-gray-700 bg-black">
          <div className="flex items-center gap-3">
            {isEditingTitle ? (
              <div className="flex items-center gap-2">
                <Input
                  value={titleInput}
                  onChange={(e) => setTitleInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") handleSaveTitle();
                    if (e.key === "Escape") handleCancelEditTitle();
                  }}
                  className="h-8 w-64 bg-gray-900 border-gray-700 text-white"
                  autoFocus
                />
                <Button
                  onClick={handleSaveTitle}
                  variant="ghost"
                  size="sm"
                  className="h-8 text-white hover:bg-gray-800"
                >
                  ✓
                </Button>
                <Button
                  onClick={handleCancelEditTitle}
                  variant="ghost"
                  size="sm"
                  className="h-8 text-white hover:bg-gray-800"
                >
                  ✕
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <h1 className="text-sm font-medium text-white">{workflowTitle}</h1>
                <Button
                  onClick={() => {
                    setTitleInput(workflowTitle);
                    setIsEditingTitle(true);
                  }}
                  variant="ghost"
                  size="sm"
                  className="h-6 px-2 text-gray-400 hover:text-white hover:bg-gray-800"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-3 h-3"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                    />
                  </svg>
                </Button>
              </div>
            )}
          </div>
        </div>

        {/* Canvas Area with Execution Panel Below (70/30 split) */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Workflow Editor - 70% height */}
          <div className="flex overflow-hidden" style={{ height: "70%" }}>
            {isPaletteOpen && (
              <NodePalette
                isOpen={isPaletteOpen}
                onClose={() => setIsPaletteOpen(false)}
                categories={Nodes.nodeCategories}
                onAddNode={addNode}
              />
            )}
            <div className="flex-1 relative">
              {!isPaletteOpen && (
                <button
                  onClick={() => setIsPaletteOpen(true)}
                  className="absolute left-4 top-4 z-10 bg-gray-800 hover:bg-gray-700 text-white px-3 py-1.5 rounded-lg shadow-lg transition-all border border-gray-700 hover:border-gray-600 text-sm font-medium flex items-center gap-1.5"
                  title="Show Nodes Panel"
                >
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
                </button>
              )}
              <ReactFlow
                nodes={nodes.map((n) => ({ ...n, dragHandle: ".__dragger" }))}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                onDragOver={onDragOver}
                onDrop={onDrop}
                nodeTypes={nodeTypes}
                fitView
                className="bg-gray-950"
              />
            </div>
          </div>

          {/* Execution Panel - 30% height, always visible */}
          <div className="border-t border-gray-700" style={{ height: "30%" }}>
            <ExecutionPanel
              isOpen={true}
              isLoading={isExecuting}
              result={executionResult}
              error={executionError}
              details={executionDetails}
              onCancel={handleCancelExecution}
              onClose={handleCloseExecutionPanel}
              height={0} // Height controlled by parent
              onHeightChange={() => {}} // Not used in fixed layout
            />
            {!isExecutionPanelOpen && !isExecuting && !executionResult && !executionError && (
              <WorkflowStatusBar nodeCount={nodes.length} edgeCount={edges.length} />
            )}
          </div>
        </div>
      </div>
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
        <Nodes.NodeContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onClose={() => setContextMenu(null)}
          onDelete={() => handleDeleteNode(contextMenu.nodeId)}
        />
      )}
      {deleteConfirm && (
        <Nodes.DeleteConfirmDialog
          nodeName={deleteConfirm.nodeName}
          onConfirm={confirmDelete}
          onCancel={() => setDeleteConfirm(null)}
        />
      )}
    </div>
  );
}

export default function Page() {
  return (
    <ReactFlowProvider>
      <SidebarProvider defaultCollapsed={false}>
        <Canvas />
      </SidebarProvider>
    </ReactFlowProvider>
  );
}
