import { NodeProps, Handle, Position, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type RendererNodeData = {
  label?: string;
  // Execution result data - populated after workflow execution
  _executionData?: unknown;
};

export function RendererNode({
  id,
  data,
  ...props
}: NodeProps<RendererNodeData>) {
  const { setNodes } = useReactFlow();

  const handleTitleChange = (newTitle: string) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, label: newTitle } } : n
      )
    );
  };

  const nodeInfo = getNodeInfo("rendererNode");
  // Type assertion is consistent with other nodes in the codebase
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const onShowOptions = (props as any).onShowOptions;

  // Helper to check if data is empty/null/undefined
  const isEmptyData = (data: unknown): boolean => {
    return !data && data !== 0 && data !== false;
  };

  // Auto-detect the best rendering mode based on data type and structure
  const inferRenderMode = (executionData: unknown): string => {
    if (isEmptyData(executionData)) {
      return "text";
    }

    if (typeof executionData === "boolean") {
      return typeof executionData;
    }

    if (
      typeof executionData === "number" ||
      typeof executionData === "bigint"
    ) {
      return "number";
    }

    // If it's a string, check if it's JSON, CSV, TSV, or XML
    if (typeof executionData === "string") {
      const trimmed = executionData.trim();

      // Check for JSON
      if (
        (trimmed.startsWith("{") && trimmed.endsWith("}")) ||
        (trimmed.startsWith("[") && trimmed.endsWith("]"))
      ) {
        try {
          JSON.parse(trimmed);
          return "json";
        } catch {
          // Not valid JSON, continue checking
        }
      }

      // Check for XML
      if (trimmed.startsWith("<") && trimmed.endsWith(">")) {
        return "xml";
      }

      // Check for CSV/TSV (has multiple lines with separators)
      const lines = trimmed.split("\n");
      if (lines.length > 1) {
        const firstLine = lines[0];
        const commas = (firstLine.match(/,/g) || []).length;
        const tabs = (firstLine.match(/\t/g) || []).length;

        if (tabs > 0 && tabs >= commas) {
          return "tsv";
        } else if (commas > 0) {
          return "csv";
        }
      }

      return "text";
    }

    // If it's an array
    if (Array.isArray(executionData)) {
      if (executionData.length === 0) {
        return "json";
      }

      const firstItem = executionData[0];

      // Check if it's array of objects with 'label' and 'value' for bar chart
      if (typeof firstItem === "object" && firstItem !== null) {
        const obj = firstItem as Record<string, unknown>;
        if (
          ("label" in obj || "name" in obj) &&
          "value" in obj &&
          typeof obj.value === "number"
        ) {
          return "bar_chart";
        }

        // Array of objects - use table
        return "table";
      }

      // Array of primitives - use table with "values" column
      return "table";
    }

    // If it's an object
    if (typeof executionData === "object" && executionData !== null) {
      return "json";
    }

    // Primitives (number, boolean) - use text
    return "text";
  };

  // Render different formats based on auto-detected mode
  const renderData = () => {
    const executionData = data?._executionData;

    if (isEmptyData(executionData)) {
      return (
        <div className="text-xs text-gray-500 italic py-2 px-1">No data</div>
      );
    }

    const mode = inferRenderMode(executionData);

    try {
      switch (mode) {
        case "json": {
          // If data is a string that looks like JSON, parse it first
          let jsonData = executionData;
          if (typeof executionData === "string") {
            try {
              jsonData = JSON.parse(executionData);
            } catch {
              jsonData = executionData;
            }
          }

          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: JSON
              </div>
              <pre className="text-[9px] leading-tight bg-gray-950 p-1.5 rounded border border-gray-700 overflow-auto max-h-48 font-mono">
                {JSON.stringify(jsonData, null, 2)}
              </pre>
            </div>
          );
        }

        case "csv":
        case "tsv": {
          const separator = mode === "csv" ? "," : "\t";
          let csvText = "";

          if (typeof executionData === "string") {
            csvText = executionData;
          } else if (Array.isArray(executionData) && executionData.length > 0) {
            // Get headers from first object
            const firstItem = executionData[0];
            if (typeof firstItem === "object" && firstItem !== null) {
              const headers = Object.keys(firstItem);
              csvText = headers.join(separator) + "\n";

              // Add rows
              executionData.forEach((item) => {
                if (typeof item === "object" && item !== null) {
                  const row = headers.map((h) => {
                    const val = (item as Record<string, unknown>)[h];
                    return val !== undefined && val !== null ? String(val) : "";
                  });
                  csvText += row.join(separator) + "\n";
                }
              });
            }
          } else {
            csvText = String(executionData);
          }

          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: {mode.toUpperCase()}
              </div>
              <pre className="text-[9px] leading-tight bg-gray-950 p-1.5 rounded border border-gray-700 overflow-auto max-h-48 font-mono whitespace-pre">
                {csvText}
              </pre>
            </div>
          );
        }

        case "xml": {
          // Simple XML formatting - convert object to XML or display string
          const toXML = (obj: unknown, rootName = "data"): string => {
            if (obj === null || obj === undefined) {
              return `<${rootName}/>`;
            }
            if (typeof obj !== "object") {
              return `<${rootName}>${String(obj)}</${rootName}>`;
            }
            if (Array.isArray(obj)) {
              return obj.map((item) => toXML(item, "item")).join("\n");
            }
            let xml = `<${rootName}>`;
            Object.entries(obj as Record<string, unknown>).forEach(
              ([key, value]) => {
                xml += `\n  ${toXML(value, key)}`;
              }
            );
            xml += `\n</${rootName}>`;
            return xml;
          };

          const xmlContent =
            typeof executionData === "string"
              ? executionData
              : toXML(executionData);

          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: XML
              </div>
              <pre className="text-[9px] leading-tight bg-gray-950 p-1.5 rounded border border-gray-700 overflow-auto max-h-48 font-mono">
                {xmlContent}
              </pre>
            </div>
          );
        }

        case "table": {
          if (Array.isArray(executionData) && executionData.length > 0) {
            const firstItem = executionData[0];
            
            // Handle array of objects
            if (typeof firstItem === "object" && firstItem !== null) {
              const headers = Object.keys(firstItem);

              return (
                <div>
                  <div className="text-[8px] text-gray-500 mb-1 px-1">
                    Format: Table
                  </div>
                  <div className="text-[9px] overflow-auto max-h-48">
                    <table className="w-full border-collapse">
                      <thead>
                        <tr className="bg-gray-800">
                          {headers.map((h, i) => (
                            <th
                              key={i}
                              className="border border-gray-700 px-1 py-0.5 text-left"
                            >
                              {h}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {executionData.slice(0, 20).map((item, rowIdx) => (
                          <tr key={rowIdx} className="hover:bg-gray-800">
                            {headers.map((h, colIdx) => {
                              const val = (item as Record<string, unknown>)[h];
                              return (
                                <td
                                  key={colIdx}
                                  className="border border-gray-700 px-1 py-0.5"
                                >
                                  {val !== undefined && val !== null
                                    ? String(val)
                                    : ""}
                                </td>
                              );
                            })}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    {executionData.length > 20 && (
                      <div className="text-gray-500 text-center py-1">
                        ... and {executionData.length - 20} more rows
                      </div>
                    )}
                  </div>
                </div>
              );
            }
            
            // Handle array of primitives (numbers, strings, etc.)
            return (
              <div>
                <div className="text-[8px] text-gray-500 mb-1 px-1">
                  Format: Table
                </div>
                <div className="text-[9px] overflow-auto max-h-48">
                  <table className="w-full border-collapse">
                    <thead>
                      <tr className="bg-gray-800">
                        <th className="border border-gray-700 px-1 py-0.5 text-left">
                          values
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {executionData.slice(0, 20).map((item, rowIdx) => (
                        <tr key={rowIdx} className="hover:bg-gray-800">
                          <td className="border border-gray-700 px-1 py-0.5">
                            {item !== undefined && item !== null
                              ? String(item)
                              : ""}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                  {executionData.length > 20 && (
                    <div className="text-gray-500 text-center py-1">
                      ... and {executionData.length - 20} more rows
                    </div>
                  )}
                </div>
              </div>
            );
          }
          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: JSON (fallback)
              </div>
              <pre className="text-[9px] leading-tight bg-gray-950 p-1.5 rounded border border-gray-700 overflow-auto max-h-48">
                {JSON.stringify(executionData, null, 2)}
              </pre>
            </div>
          );
        }

        case "bar_chart": {
          if (Array.isArray(executionData) && executionData.length > 0) {
            // Find max value for scaling
            let maxValue = 0;
            const chartData = executionData.slice(0, 20).map((item, index) => {
              if (typeof item === "object" && item !== null) {
                const obj = item as Record<string, unknown>;
                const label = obj.label
                  ? String(obj.label)
                  : String(obj.name || "");
                const value = typeof obj.value === "number" ? obj.value : 0;
                maxValue = Math.max(maxValue, value);
                return { label, value };
              } else if (typeof item === "number") {
                maxValue = Math.max(maxValue, item);
                return { label: String(index), value: item };
              }
              return { label: "", value: 0 };
            });

            return (
              <div>
                <div className="text-[8px] text-gray-500 mb-1 px-1">
                  Format: Bar Chart
                </div>
                <div className="text-[9px] py-1 px-1 max-h-48 overflow-auto">
                  {chartData.map((item, idx) => (
                    <div key={idx} className="flex items-center gap-1 mb-1">
                      <span
                        className="w-12 truncate text-gray-400"
                        title={item.label}
                      >
                        {item.label}
                      </span>
                      <div className="flex-1 bg-gray-800 rounded h-3 relative">
                        <div
                          className="bg-blue-500 h-full rounded transition-all"
                          style={{
                            width: `${
                              maxValue > 0 ? (item.value / maxValue) * 100 : 0
                            }%`,
                          }}
                        />
                      </div>
                      <span className="w-8 text-right text-gray-300">
                        {item.value}
                      </span>
                    </div>
                  ))}
                  {executionData.length > 20 && (
                    <div className="text-gray-500 text-center py-1">
                      ... and {executionData.length - 20} more items
                    </div>
                  )}
                </div>
              </div>
            );
          }
          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: Text (fallback)
              </div>
              <div className="text-xs text-gray-500 italic py-2 px-1">
                Data must be an array for bar chart
              </div>
            </div>
          );
        }

        case "number": {
          return (
            <div className="bg-gray-950 text-2xl rounded border border-gray-700 overflow-auto text-center">
              {typeof executionData === "number" ||
              typeof executionData === "bigint"
                ? executionData
                : JSON.stringify(executionData, null, 2)}
            </div>
          );
        }

        case "boolean": {
          return typeof executionData === "boolean" ? (
            executionData ? (
              <div className="bg-green-500 text-2xl rounded border border-gray-700 overflow-auto text-center">
                ⬆️
              </div>
            ) : (
              <div className="bg-red-500 text-2xl rounded border border-gray-700 overflow-auto text-center">
                ⬇️
              </div>
            )
          ) : (
            JSON.stringify(executionData, null, 2)
          );
        }

        case "text":
        default:
          return (
            <div>
              <div className="text-[8px] text-gray-500 mb-1 px-1">
                Format: Plain Text
              </div>
              <pre className="text-[9px] leading-tight bg-gray-950 p-1.5 rounded border border-gray-700 overflow-auto max-h-48 whitespace-pre-wrap wrap-break-word">
                {typeof executionData === "string"
                  ? executionData
                  : JSON.stringify(executionData, null, 2)}
              </pre>
            </div>
          );
      }
    } catch (error) {
      return (
        <div className="text-xs text-red-400 py-2 px-1">
          Error rendering data:{" "}
          {error instanceof Error ? error.message : "Unknown error"}
        </div>
      );
    }
  };

  return (
    <NodeWrapper
      title={String(data?.label || "Renderer")}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      onTitleChange={handleTitleChange}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="w-2 h-2 bg-blue-400"
      />

      <div className="flex flex-col gap-1">
        <div className="w-64 border border-gray-600 rounded bg-gray-900 mt-1">
          {renderData()}
        </div>
      </div>

      <Handle
        type="source"
        position={Position.Right}
        className="w-2 h-2 bg-green-400"
      />
    </NodeWrapper>
  );
}
