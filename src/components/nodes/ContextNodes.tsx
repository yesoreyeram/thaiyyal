import { NodeProps, useReactFlow } from "reactflow";
import React from "react";
import { NodeWrapper } from "./NodeWrapper";
import { getNodeInfo } from "./nodeInfo";

type ContextNodeData = {
  context_name?: string;
  context_value?: string | number;
  label?: string;
};

type ContextNodeConfig = {
  type: "variable" | "constant";
  icon: string;
  color: {
    gradient: string;
    border: string;
    borderHover: string;
    text: string;
    inputBorder: string;
    focusRing: string;
    templateText: string;
  };
};

const configs: Record<"variable" | "constant", ContextNodeConfig> = {
  variable: {
    type: "variable",
    icon: "📦",
    color: {
      gradient: "from-purple-700 to-purple-800",
      border: "border-purple-400",
      borderHover: "hover:border-purple-300",
      text: "text-purple-200",
      inputBorder: "border-purple-600",
      focusRing: "focus:ring-purple-400",
      templateText: "text-purple-300",
    },
  },
  constant: {
    type: "constant",
    icon: "🔒",
    color: {
      gradient: "from-amber-700 to-amber-800",
      border: "border-amber-400",
      borderHover: "hover:border-amber-300",
      text: "text-amber-200",
      inputBorder: "border-amber-600",
      focusRing: "focus:ring-amber-400",
      templateText: "text-amber-300",
    },
  },
};

export function BaseContextNode({ 
  id, 
  data,
  config,
  ...props
}: NodeProps<ContextNodeData> & { config: ContextNodeConfig }) {
  const { setNodes } = useReactFlow();
  
  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const context_name = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, context_name } } : n
      )
    );
  };

  const onValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const context_value = e.target.value;
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, context_value } } : n
      )
    );
  };

  const inputClassName = `w-full text-xs border ${config.color.inputBorder} px-2 py-1 rounded bg-gray-900 text-white placeholder-gray-500 focus:ring-2 ${config.color.focusRing} focus:outline-none`;

  const label = data?.label || (config.type === "variable" ? "Variable" : "Constant");
  const templatePrefix = config.type === "variable" ? "variable" : "const";
  const contextName = data?.context_name || "...";

  const nodeInfo = getNodeInfo(config.type === "variable" ? "contextVariableNode" : "contextConstantNode");
  const onShowOptions = (props as any).onShowOptions;

  return (
    <NodeWrapper
      title={label}
      nodeInfo={nodeInfo}
      onShowOptions={onShowOptions}
      className={`bg-gradient-to-br ${config.color.gradient} text-white shadow-lg rounded-lg border-2 ${config.color.border} ${config.color.borderHover} transition-all`}
    >
      <div className="space-y-1">
        <input
          value={String(data?.context_name ?? "")}
          type="text"
          onChange={onNameChange}
          className={inputClassName}
          placeholder="Name"
          aria-label={`${config.type === "variable" ? "Variable" : "Constant"} name`}
        />
        <input
          value={String(data?.context_value ?? "")}
          type="text"
          onChange={onValueChange}
          className={inputClassName}
          placeholder="Value"
          aria-label={`${config.type === "variable" ? "Variable" : "Constant"} value`}
        />
      </div>
      <div className={`mt-1 text-xs ${config.color.templateText} font-mono`}>
        {`{{ ${templatePrefix}.${contextName} }}`}
      </div>
    </NodeWrapper>
  );
}

export function ContextVariableNode(props: NodeProps<ContextNodeData>) {
  return <BaseContextNode {...props} config={configs.variable} />;
}

export function ContextConstantNode(props: NodeProps<ContextNodeData>) {
  return <BaseContextNode {...props} config={configs.constant} />;
}
