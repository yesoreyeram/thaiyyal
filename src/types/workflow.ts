export interface Workflow {
  id: string;
  title: string;
  data: {
    nodes: unknown[];
    edges: unknown[];
  };
  createdAt: string;
  updatedAt: string;
}
