// Node information and documentation for all node types

export interface NodeInfo {
  description: string;
  inputs?: string[];
  outputs?: string[];
}

export const nodeInfoMap: Record<string, NodeInfo> = {
  numberNode: {
    description: "Input a numeric value that can be used in calculations and operations.",
    inputs: [],
    outputs: ["Number value"],
  },
  textInputNode: {
    description: "Input text that can be processed, transformed, or concatenated with other text.",
    inputs: [],
    outputs: ["Text string"],
  },
  opNode: {
    description: "Perform arithmetic operations on two numeric inputs. Supports addition, subtraction, multiplication, and division.",
    inputs: ["First number", "Second number"],
    outputs: ["Result of operation"],
  },
  textOpNode: {
    description: "Transform text using various operations like uppercase, lowercase, title case, camel case, concatenation, and repetition.",
    inputs: ["Text input(s)"],
    outputs: ["Transformed text"],
  },
  httpNode: {
    description: "Make HTTP GET requests to external APIs and services. Returns the response data.",
    inputs: [],
    outputs: ["HTTP response data"],
  },
  vizNode: {
    description: "Visualize data in different formats. Supports text and table display modes.",
    inputs: ["Data to visualize"],
    outputs: [],
  },
  conditionNode: {
    description: "Branch execution based on a condition. Supports comparison operators: >, <, >=, <=, ==, !=.",
    inputs: ["Value to compare"],
    outputs: ["True branch", "False branch"],
  },
  filterNode: {
    description: "Filter JSON array elements based on expression. Use 'item.field' to access array element properties (e.g., item.age >= 18).",
    inputs: ["Array to filter"],
    outputs: ["Filtered array"],
  },
  forEachNode: {
    description: "Iterate over array elements with a configurable maximum iteration limit.",
    inputs: ["Array or collection"],
    outputs: ["Each array element"],
  },
  whileLoopNode: {
    description: "Loop while a condition is true, with a safety maximum iteration limit.",
    inputs: ["Condition input"],
    outputs: ["Loop output"],
  },
  switchNode: {
    description: "Multi-way branching based on input value. Routes execution to matching case or default path.",
    inputs: ["Value to match"],
    outputs: ["Multiple case outputs"],
  },
  parallelNode: {
    description: "Execute multiple operations concurrently with configurable concurrency limit.",
    inputs: ["Operation inputs"],
    outputs: ["Parallel results"],
  },
  joinNode: {
    description: "Combine multiple inputs using different strategies: all (wait for all), any (first available), or first (first to arrive).",
    inputs: ["Multiple input streams"],
    outputs: ["Joined result"],
  },
  splitNode: {
    description: "Distribute data to multiple downstream paths for parallel processing.",
    inputs: ["Data to split"],
    outputs: ["Multiple output paths"],
  },
  variableNode: {
    description: "Store and retrieve variables in workflow memory. Supports get and set operations.",
    inputs: ["Value to store (for set)"],
    outputs: ["Retrieved value (for get)"],
  },
  extractNode: {
    description: "Extract specific fields from objects or data structures.",
    inputs: ["Object or data structure"],
    outputs: ["Extracted fields"],
  },
  transformNode: {
    description: "Transform data structures between formats: to array, to object, flatten, extract keys, or extract values.",
    inputs: ["Data to transform"],
    outputs: ["Transformed data"],
  },
  accumulatorNode: {
    description: "Accumulate values using operations: sum, product, concatenation, array building, or counting.",
    inputs: ["Values to accumulate"],
    outputs: ["Accumulated result"],
  },
  counterNode: {
    description: "Maintain a counter with operations: increment, decrement, reset, or get current value.",
    inputs: [],
    outputs: ["Counter value"],
  },
  delayNode: {
    description: "Pause workflow execution for a specified duration (e.g., 1s, 100ms, 5m).",
    inputs: ["Data to delay"],
    outputs: ["Delayed data"],
  },
  cacheNode: {
    description: "Cache data with TTL (time-to-live). Supports get, set, and delete operations.",
    inputs: ["Value to cache (for set)"],
    outputs: ["Cached value (for get)"],
  },
  retryNode: {
    description: "Retry failed operations with configurable backoff strategies: exponential, linear, or constant delay.",
    inputs: ["Operation to retry"],
    outputs: ["Operation result or error"],
  },
  tryCatchNode: {
    description: "Handle errors gracefully with fallback logic. Continue workflow execution even when errors occur.",
    inputs: ["Try branch operation"],
    outputs: ["Success or fallback result"],
  },
  timeoutNode: {
    description: "Enforce time limits on operations. Can error or continue with partial results on timeout.",
    inputs: ["Operation with time limit"],
    outputs: ["Result or timeout error"],
  },
  contextVariableNode: {
    description: "Store and retrieve context variables that persist across workflow execution.",
    inputs: ["Context value"],
    outputs: ["Context data"],
  },
  contextConstantNode: {
    description: "Define constant values in the workflow context that cannot be modified during execution.",
    inputs: [],
    outputs: ["Constant value"],
  },
};

export function getNodeInfo(nodeType: string): NodeInfo | undefined {
  return nodeInfoMap[nodeType];
}
