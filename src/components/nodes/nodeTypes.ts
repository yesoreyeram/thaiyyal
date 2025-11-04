import { NodeProps } from "reactflow";

/**
 * Extended props to include custom handlers for node components
 */
export type NodePropsWithOptions<T = Record<string, unknown>> = NodeProps<T> & {
  onShowOptions?: (x: number, y: number) => void;
  onOpenInfo?: () => void;
};
