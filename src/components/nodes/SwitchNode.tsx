import React, { useCallback } from "react";
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

  const moveCaseUp = (caseIndex: number) => {
    if (caseIndex === 0 || cases[caseIndex].is_default) return;
    const updatedCases = [...cases];
    [updatedCases[caseIndex - 1], updatedCases[caseIndex]] = 
      [updatedCases[caseIndex], updatedCases[caseIndex - 1]];
    updateCases(updatedCases);
  };

  const moveCaseDown = (caseIndex: number) => {
    if (caseIndex >= cases.length - 1 || cases[caseIndex].is_default) return;
    // Don't move past default case
    if (cases[caseIndex + 1]?.is_default) return;
    const updatedCases = [...cases];
    [updatedCases[caseIndex], updatedCases[caseIndex + 1]] = 
      [updatedCases[caseIndex + 1], updatedCases[caseIndex]];
    updateCases(updatedCases);
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
          
          const isFirst = caseIndex === 0;
          const isLast = caseIndex === cases.length - 2; // -2 because last is default
          
          return (
            <div
              key={caseIndex}
              className="noDrag relative flex items-center gap-1 p-1 bg-gray-800 border border-gray-600 rounded hover:border-blue-400 transition-colors"
            >
              {/* Reorder buttons */}
              <div className="flex flex-col flex-shrink-0">
                <button
                  onClick={() => moveCaseUp(caseIndex)}
                  disabled={isFirst}
                  className={`text-xs leading-none ${
                    isFirst
                      ? "text-gray-700 cursor-not-allowed"
                      : "text-gray-400 hover:text-white cursor-pointer"
                  }`}
                  title="Move up"
                  aria-label="Move case up"
                >
                  ▲
                </button>
                <button
                  onClick={() => moveCaseDown(caseIndex)}
                  disabled={isLast}
                  className={`text-xs leading-none ${
                    isLast
                      ? "text-gray-700 cursor-not-allowed"
                      : "text-gray-400 hover:text-white cursor-pointer"
                  }`}
                  title="Move down"
                  aria-label="Move case down"
                >
                  ▼
                </button>
              </div>
              
              {/* Expression input */}
              <input
                type="text"
                value={c.when}
                onChange={(e) => updateCaseWhen(caseIndex, e.target.value)}
                draggable={false}
                placeholder="e.g., input > 50"
                className="flex-1 text-xs border-0 border-b border-gray-600 px-1 py-0.5 bg-transparent text-white focus:border-blue-400 focus:outline-none"
              />
              
              {/* Delete button with trash icon */}
              <button
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
          className="noDrag w-full text-xs border border-dashed border-gray-600 px-2 py-1 rounded bg-gray-800 text-gray-400 hover:text-white hover:border-blue-400 transition-colors"
        >
          + Add Case
        </button>
        
        {/* Default case - not draggable, more compact */}
        {defaultCase && (
          <div className="noDrag relative flex items-center gap-1 p-1 bg-gray-700 border-2 border-yellow-600/50 rounded">
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
