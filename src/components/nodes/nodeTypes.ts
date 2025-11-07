import { NodeProps } from "reactflow";

/**
 * Extended props to include custom handlers for node components
 */
export type NodePropsWithOptions<T = Record<string, unknown>> = NodeProps<T> & {
  data: T & { label?: string };
  onShowOptions?: (x: number, y: number) => void;
  onOpenInfo?: () => void;
  onDelete?: () => void;
};
