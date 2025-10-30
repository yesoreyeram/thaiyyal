export interface Workflow {
  id: string;
  title: string;
  data: {
    nodes: any[];
    edges: any[];
  };
  createdAt: string;
  updatedAt: string;
}
