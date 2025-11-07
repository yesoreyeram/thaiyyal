import React, { useState, useCallback } from "react";
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

  const updateCase = (caseIndex: number, field: keyof SwitchCase, value: string) => {
    const updatedCases = [...cases];
    updatedCases[caseIndex] = { ...updatedCases[caseIndex], [field]: value };
    updateCases(updatedCases);
  };

  const deleteCase = (caseIndex: number) => {
    const updatedCases = cases.filter((_, i) => i !== caseIndex);
    updateCases(updatedCases);
  };

  const handleDragStart = (e: React.DragEvent, caseIndex: number) => {
    e.stopPropagation();
    setDraggedIndex(caseIndex);
  };

  const handleDragOver = (e: React.DragEvent, caseIndex: number) => {
    e.preventDefault();
    e.stopPropagation();
    if (draggedIndex === null || draggedIndex === caseIndex) return;
    
    // Don't allow dragging over default case
    if (cases[caseIndex]?.is_default) return;

    const updatedCases = [...cases];
    const draggedCase = updatedCases[draggedIndex];
    
    // Remove dragged item
    updatedCases.splice(draggedIndex, 1);
    // Insert at new position
    updatedCases.splice(caseIndex, 0, draggedCase);
    
    // Ensure default stays last
    const defaultIdx = updatedCases.findIndex((c) => c.is_default);
    if (defaultIdx !== -1 && defaultIdx !== updatedCases.length - 1) {
      const def = updatedCases.splice(defaultIdx, 1)[0];
      updatedCases.push(def);
    }
    
    updateCases(updatedCases);
    setDraggedIndex(caseIndex);
  };

  const handleDragEnd = (e: React.DragEvent) => {
    e.stopPropagation();
    setDraggedIndex(null);
  };

  const nodeInfo = getNodeInfo("switchNode");

  // Calculate handle positions (52px header + 8px padding + accumulated heights)
  const headerHeight = 52;
  const padding = 8;
  const caseHeight = 45; // approximate height per case row
  const buttonHeight = 28;
  let currentY = headerHeight + padding;

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
          
          const handleY = currentY + 22; // Center of the case row
          currentY += caseHeight;
          
          return (
            <div
              key={caseIndex}
              className={`relative flex flex-col gap-1 p-1 bg-gray-800 border border-gray-600 rounded hover:border-blue-400 transition-colors ${
                draggedIndex === caseIndex ? "opacity-50" : ""
              }`}
            >
              <div className="flex items-center gap-1">
                {/* Drag handle - only draggable part */}
                <div
                  draggable
                  onDragStart={(e) => handleDragStart(e, caseIndex)}
                  onDragOver={(e) => handleDragOver(e, caseIndex)}
                  onDragEnd={handleDragEnd}
                  className="flex-shrink-0 cursor-move text-gray-500 hover:text-gray-300 px-1 text-xs"
                  title="Drag to reorder"
                >
                  ⋮⋮
                </div>
                
                {/* Expression input */}
                <input
                  type="text"
                  value={c.when}
                  onChange={(e) => updateCase(caseIndex, "when", e.target.value)}
                  placeholder="e.g., input > 50"
                  className="flex-1 text-xs border-0 border-b border-gray-600 px-1 py-0.5 bg-transparent text-white focus:border-blue-400 focus:outline-none"
                  onClick={(e) => e.stopPropagation()}
                />
                
                {/* Delete button */}
                <button
                  onClick={() => deleteCase(caseIndex)}
                  className="flex-shrink-0 text-red-400 hover:text-red-300 text-sm font-bold px-1"
                  title="Delete case"
                >
                  ×
                </button>
              </div>
              
              {/* Output path - second row, smaller */}
              <input
                type="text"
                value={c.output_path || ""}
                onChange={(e) => updateCase(caseIndex, "output_path", e.target.value)}
                placeholder="output path"
                className="text-xs border-0 border-b border-gray-600 px-1 py-0.5 ml-5 bg-transparent text-gray-400 focus:border-blue-400 focus:outline-none focus:text-white"
                onClick={(e) => e.stopPropagation()}
              />
              
              {/* Dynamic handle for this case - positioned relative to NodeWrapper */}
              <Handle
                type="source"
                id={c.output_path || `case_${caseIndex}`}
                position={Position.Right}
                style={{ 
                  top: `${handleY}px`,
                  right: 0,
                  zIndex: 10
                }}
                className="w-2 h-2 bg-blue-500"
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
        {defaultCase && (() => {
          currentY += buttonHeight + 4; // button height + gap
          const defaultHandleY = currentY + 12; // Center of default case
          
          return (
            <div className="relative flex items-center gap-1 p-1 bg-gray-700 border-2 border-yellow-600/50 rounded">
              {/* Default label - not editable */}
              <div className="flex-shrink-0 text-xs text-yellow-400 font-semibold px-1">
                default
              </div>
              
              {/* Default handle */}
              <Handle
                type="source"
                id={defaultCase.output_path || "default"}
                position={Position.Right}
                style={{ 
                  top: `${defaultHandleY}px`,
                  right: 0,
                  zIndex: 10
                }}
                className="w-2 h-2 bg-yellow-500"
              />
            </div>
          );
        })()}
      </div>
    </NodeWrapper>
  );
}
