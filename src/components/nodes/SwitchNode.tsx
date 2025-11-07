import React, { useCallback, useState } from "react";
import { Handle, Position, useReactFlow } from "reactflow";
import { NodePropsWithOptions } from "./nodeTypes";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type SwitchCase = {
  when: string;
  output_path?: string;
  is_default?: boolean;
};

type SwitchNodeData = {
  cases?: SwitchCase[];
  label?: string;
};

export function SwitchNode({
  id,
  data,
  onShowOptions,
}: NodePropsWithOptions<SwitchNodeData>) {
  const { setNodes } = useReactFlow();
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);
  const [dragOverIndex, setDragOverIndex] = useState<number | null>(null);

  const cases = data?.cases || [];
  const nonDefaultCases = cases.filter((c) => !c.is_default);
  const defaultCase = cases.find((c) => c.is_default);

  const updateCases = useCallback(
    (newCases: SwitchCase[]) => {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === id ? { ...n, data: { ...n.data, cases: newCases } } : n
        )
      );
    },
    [id, setNodes]
  );

  const addCase = () => {
    const newCase: SwitchCase = {
      when: "",
      output_path: "case_" + (nonDefaultCases.length + 1),
    };
    const updatedCases = defaultCase
      ? [...nonDefaultCases, newCase, defaultCase]
      : [...nonDefaultCases, newCase];
    updateCases(updatedCases);
  };

  const updateCaseWhen = (caseIndex: number, value: string) => {
    const updatedCases = [...cases];
    // Auto-generate output_path using case index if not set
    const outputPath = updatedCases[caseIndex].output_path || 
                       `case_${caseIndex + 1}`;
    updatedCases[caseIndex] = { 
      ...updatedCases[caseIndex], 
      when: value,
      output_path: outputPath
    };
    updateCases(updatedCases);
  };

  const deleteCase = (caseIndex: number) => {
    const updatedCases = cases.filter((_, i) => i !== caseIndex);
    updateCases(updatedCases);
  };

  // Drag-and-drop handlers
  const handleDragStart = (e: React.DragEvent, caseIndex: number) => {
    if (cases[caseIndex].is_default) return; // Don't allow dragging default case
    setDraggedIndex(caseIndex);
    e.dataTransfer.effectAllowed = "move";
    // Prevent ReactFlow from handling this drag
    e.stopPropagation();
  };

  const handleDragOver = (e: React.DragEvent, caseIndex: number) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (draggedIndex === null || draggedIndex === caseIndex) return;
    if (cases[caseIndex].is_default) return; // Don't allow dropping on default case
    
    setDragOverIndex(caseIndex);
  };

  const handleDragEnd = (e: React.DragEvent) => {
    e.stopPropagation();
    
    if (draggedIndex === null || dragOverIndex === null || draggedIndex === dragOverIndex) {
      setDraggedIndex(null);
      setDragOverIndex(null);
      return;
    }

    // Reorder cases
    const updatedCases = [...cases];
    const [draggedCase] = updatedCases.splice(draggedIndex, 1);
    updatedCases.splice(dragOverIndex, 0, draggedCase);
    
    updateCases(updatedCases);
    setDraggedIndex(null);
    setDragOverIndex(null);
  };

  const handleDragLeave = () => {
    setDragOverIndex(null);
  };

  const nodeInfo = getNodeInfo("switchNode");

  return (
    <NodeWrapper
      id={id}
      title={String(data?.label || "Switch")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />
      
      <div className="w-80 max-h-96 overflow-y-auto space-y-1 p-2">
        {/* Non-default cases */}
        {cases.map((c, caseIndex) => {
          if (c.is_default) return null; // Skip default, render separately below
          
          const isDragging = draggedIndex === caseIndex;
          const isDragOver = dragOverIndex === caseIndex;
          
          return (
            <div
              key={caseIndex}
              draggable={!c.is_default}
              onDragStart={(e) => handleDragStart(e, caseIndex)}
              onDragOver={(e) => handleDragOver(e, caseIndex)}
              onDragEnd={handleDragEnd}
              onDragLeave={handleDragLeave}
              className={`relative flex items-center gap-1 p-1 bg-gray-800 border rounded transition-all ${
                isDragging 
                  ? 'opacity-50 border-blue-500' 
                  : isDragOver 
                  ? 'border-blue-400 border-2' 
                  : 'border-gray-600 hover:border-blue-400'
              }`}
              style={{ cursor: c.is_default ? 'default' : 'grab' }}
            >
              {/* Drag handle icon */}
              <div 
                className="flex-shrink-0 text-gray-500 text-xs cursor-grab active:cursor-grabbing select-none"
                style={{ cursor: 'grab' }}
              >
                ⋮⋮
              </div>
              
              {/* Expression input */}
              <input
                type="text"
                value={c.when}
                onChange={(e) => updateCaseWhen(caseIndex, e.target.value)}
                onMouseDown={(e) => e.stopPropagation()}
                onPointerDown={(e) => e.stopPropagation()}
                placeholder="e.g., input > 50"
                className="flex-1 text-xs border-0 border-b border-gray-600 px-1 py-0.5 bg-transparent text-white focus:border-blue-400 focus:outline-none"
              />
              
              {/* Delete button with trash icon */}
              <button
                onMouseDown={(e) => e.stopPropagation()}
                onPointerDown={(e) => e.stopPropagation()}
                onClick={(e) => {
                  e.stopPropagation();
                  deleteCase(caseIndex);
                }}
                className="flex-shrink-0 text-red-400 hover:text-red-300 text-xs px-1"
                title="Delete case"
                aria-label="Delete case"
              >
                🗑️
              </button>
              
              {/* Dynamic handle for this case - positioned on the right edge of this div */}
              <Handle
                type="source"
                id={c.output_path || `case_${caseIndex}`}
                position={Position.Right}
                style={{ 
                  position: 'absolute',
                  right: '-4px',
                  top: '50%',
                  transform: 'translateY(-50%)',
                }}
                className="!w-2 !h-2 !bg-blue-500"
              />
            </div>
          );
        })}
        
        {/* Add case button */}
        <button
          onClick={addCase}
          className="w-full text-xs border border-dashed border-gray-600 px-2 py-1 rounded bg-gray-800 text-gray-400 hover:text-white hover:border-blue-400 transition-colors"
        >
          + Add Case
        </button>
        
        {/* Default case - not draggable, more compact */}
        {defaultCase && (
          <div className="relative flex items-center gap-1 p-1 bg-gray-700 border-2 border-yellow-600/50 rounded">
            {/* Default label - not editable */}
            <div className="flex-1 text-xs text-yellow-400 font-semibold px-1">
              default
            </div>
            
            {/* Default handle */}
            <Handle
              type="source"
              id={defaultCase.output_path || "default"}
              position={Position.Right}
              style={{ 
                position: 'absolute',
                right: '-4px',
                top: '50%',
                transform: 'translateY(-50%)',
              }}
              className="!w-2 !h-2 !bg-yellow-500"
            />
          </div>
        )}
      </div>
    </NodeWrapper>
  );
}
