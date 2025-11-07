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
      output_path: "",
    };
    const updatedCases = defaultCase
      ? [...nonDefaultCases, newCase, defaultCase]
      : [...nonDefaultCases, newCase];
    updateCases(updatedCases);
  };

  const updateCase = (index: number, field: keyof SwitchCase, value: string) => {
    const updatedCases = [...cases];
    updatedCases[index] = { ...updatedCases[index], [field]: value };
    updateCases(updatedCases);
  };

  const deleteCase = (index: number) => {
    const updatedCases = cases.filter((_, i) => i !== index);
    updateCases(updatedCases);
  };

  const handleDragStart = (index: number) => {
    setDraggedIndex(index);
  };

  const handleDragOver = (e: React.DragEvent, index: number) => {
    e.preventDefault();
    if (draggedIndex === null || draggedIndex === index) return;
    
    // Don't allow dragging over default case
    if (cases[index]?.is_default) return;

    const updatedCases = [...cases];
    const draggedCase = updatedCases[draggedIndex];
    
    // Remove dragged item
    updatedCases.splice(draggedIndex, 1);
    // Insert at new position
    updatedCases.splice(index, 0, draggedCase);
    
    // Ensure default stays last
    const defaultIdx = updatedCases.findIndex((c) => c.is_default);
    if (defaultIdx !== -1 && defaultIdx !== updatedCases.length - 1) {
      const def = updatedCases.splice(defaultIdx, 1)[0];
      updatedCases.push(def);
    }
    
    updateCases(updatedCases);
    setDraggedIndex(index);
  };

  const handleDragEnd = () => {
    setDraggedIndex(null);
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
      
      <div className="w-80 max-h-96 overflow-y-auto space-y-2 p-2">
        {/* Non-default cases */}
        {nonDefaultCases.map((c, i) => (
          <div
            key={i}
            draggable
            onDragStart={() => handleDragStart(i)}
            onDragOver={(e) => handleDragOver(e, i)}
            onDragEnd={handleDragEnd}
            className={`relative border rounded-lg p-2 bg-gray-800 border-gray-600 hover:border-blue-400 transition-colors cursor-move ${
              draggedIndex === i ? "opacity-50" : ""
            }`}
          >
            {/* Drag handle icon */}
            <div className="absolute left-1 top-1 text-gray-500 text-xs">
              ⋮⋮
            </div>
            
            {/* Case number */}
            <div className="text-xs text-gray-400 mb-1 ml-4">
              Case {i + 1}
            </div>
            
            {/* When expression */}
            <input
              type="text"
              value={c.when}
              onChange={(e) => updateCase(i, "when", e.target.value)}
              placeholder="e.g., input > 50"
              className="w-full text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none mb-1"
            />
            
            {/* Output path */}
            <input
              type="text"
              value={c.output_path || ""}
              onChange={(e) => updateCase(i, "output_path", e.target.value)}
              placeholder="Output path"
              className="w-full text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-1 focus:ring-blue-400 focus:outline-none"
            />
            
            {/* Delete button */}
            <button
              onClick={() => deleteCase(i)}
              className="absolute top-1 right-1 text-red-400 hover:text-red-300 text-sm font-bold"
              title="Delete case"
            >
              ×
            </button>
            
            {/* Dynamic handle for this case */}
            <Handle
              type="source"
              id={c.output_path || `case_${i}`}
              position={Position.Right}
              style={{ top: `${30 + i * 80}px` }}
              className="w-2 h-2 bg-blue-500"
            />
          </div>
        ))}
        
        {/* Add case button */}
        <button
          onClick={addCase}
          className="w-full text-xs border border-dashed border-gray-600 px-2 py-2 rounded bg-gray-800 text-gray-400 hover:text-white hover:border-blue-400 transition-colors"
        >
          + Add Case
        </button>
        
        {/* Default case - not draggable */}
        {defaultCase && (
          <div className="relative border-2 rounded-lg p-2 bg-gray-700 border-gray-500">
            <div className="text-xs text-yellow-400 mb-1 font-semibold">
              🔒 Default Case (always last)
            </div>
            
            {/* When label (not editable) */}
            <div className="text-xs px-2 py-1 rounded bg-gray-800 text-gray-400 mb-1">
              {defaultCase.when || "default"}
            </div>
            
            {/* Output path */}
            <input
              type="text"
              value={defaultCase.output_path || ""}
              onChange={(e) =>
                updateCase(cases.length - 1, "output_path", e.target.value)
              }
              placeholder="Default output path"
              className="w-full text-xs border border-gray-600 px-2 py-1 rounded bg-gray-900 text-white focus:ring-1 focus:ring-yellow-400 focus:outline-none"
            />
            
            {/* Default handle */}
            <Handle
              type="source"
              id={defaultCase.output_path || "default"}
              position={Position.Right}
              style={{ bottom: "10px" }}
              className="w-2 h-2 bg-yellow-500"
            />
          </div>
        )}
        
        {/* Show message if no default case */}
        {!defaultCase && (
          <div className="text-xs text-red-400 p-2 border border-red-400 rounded bg-red-900/20">
            ⚠️ Warning: Default case required
          </div>
        )}
      </div>
    </NodeWrapper>
  );
}
