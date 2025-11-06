export const nodeCategories = [
  {
    name: "Input/Output",
    nodes: [
      {
        type: "number",
        label: "Number",
        color: "bg-blue-600",
        defaultData: { value: 0 },
      },
      {
        type: "text_input",
        label: "Text Input",
        color: "bg-green-600",
        defaultData: { text: "" },
      },
      {
        type: "boolean_input",
        label: "Boolean",
        color: "bg-indigo-600",
        defaultData: { boolean_value: false },
      },
      {
        type: "date_input",
        label: "Date",
        color: "bg-cyan-600",
        defaultData: { date_value: "" },
      },
      {
        type: "datetime_input",
        label: "DateTime",
        color: "bg-teal-600",
        defaultData: { datetime_value: "" },
      },
      {
        type: "http",
        label: "HTTP",
        color: "bg-purple-600",
        defaultData: { url: "" },
      },
      {
        type: "visualization",
        label: "Visualization",
        color: "bg-indigo-600",
        defaultData: { mode: "text" },
      },
      {
        type: "visualization",
        label: "Bar Chart",
        color: "bg-violet-600",
        defaultData: {
          orientation: "vertical",
          bar_color: "#3b82f6",
          bar_width: "medium",
          show_values: true,
          max_bars: 20,
        },
      },
      {
        type: "renderer",
        label: "Renderer",
        color: "bg-pink-600",
        defaultData: {},
      },
    ],
  },
  {
    name: "Operations",
    nodes: [
      {
        type: "operation",
        label: "Math Op",
        color: "bg-yellow-600",
        defaultData: { op: "add" },
      },
      {
        type: "text_operation",
        label: "Text Op",
        color: "bg-green-600",
        defaultData: { text_op: "uppercase" },
      },
      {
        type: "expression",
        label: "Expression",
        color: "bg-cyan-600",
        defaultData: { expression: "input * 2" },
      },
      {
        type: "transform",
        label: "Transform",
        color: "bg-teal-600",
        defaultData: { transform_type: "to_array" },
      },
      {
        type: "parse",
        label: "Parse",
        color: "bg-purple-600",
        defaultData: { input_type: "AUTO" },
      },
      {
        type: "extract",
        label: "Extract",
        color: "bg-sky-600",
        defaultData: { field: "" },
      },
    ],
  },
  {
    name: "Control Flow",
    nodes: [
      {
        type: "condition",
        label: "Condition",
        color: "bg-amber-600",
        defaultData: { condition: ">0" },
      },
      {
        type: "filter",
        label: "Filter",
        color: "bg-purple-600",
        defaultData: { condition: "item.age > 0" },
      },
      {
        type: "for_each",
        label: "For Each",
        color: "bg-orange-600",
        defaultData: { max_iterations: 1000 },
      },
      {
        type: "while_loop",
        label: "While Loop",
        color: "bg-red-600",
        defaultData: { condition: ">0", max_iterations: 100 },
      },
      {
        type: "switch",
        label: "Switch",
        color: "bg-pink-600",
        defaultData: { cases: [], default_path: "default" },
      },
    ],
  },
  {
    name: "Parallel & Join",
    nodes: [
      {
        type: "parallel",
        label: "Parallel",
        color: "bg-violet-600",
        defaultData: { max_concurrency: 10 },
      },
      {
        type: "join",
        label: "Join",
        color: "bg-fuchsia-600",
        defaultData: { join_strategy: "all" },
      },
      {
        type: "split",
        label: "Split",
        color: "bg-rose-600",
        defaultData: { paths: ["path1", "path2"] },
      },
    ],
  },
  {
    name: "State & Memory",
    nodes: [
      {
        type: "variable",
        label: "Variable",
        color: "bg-sky-600",
        defaultData: { var_name: "", var_op: "get" },
      },
      {
        type: "accumulator",
        label: "Accumulator",
        color: "bg-blue-500",
        defaultData: { accum_op: "sum" },
      },
      {
        type: "counter",
        label: "Counter",
        color: "bg-indigo-500",
        defaultData: { counter_op: "increment", delta: 1 },
      },
      {
        type: "cache",
        label: "Cache",
        color: "bg-purple-500",
        defaultData: { cache_op: "get", cache_key: "" },
      },
    ],
  },
  {
    name: "Error Handling",
    nodes: [
      {
        type: "retry",
        label: "Retry",
        color: "bg-red-500",
        defaultData: {
          max_attempts: 3,
          backoff_strategy: "exponential",
          initial_delay: "1s",
        },
      },
      {
        type: "try_catch",
        label: "Try-Catch",
        color: "bg-orange-500",
        defaultData: { continue_on_error: true },
      },
      {
        type: "timeout",
        label: "Timeout",
        color: "bg-amber-500",
        defaultData: { timeout: "30s", timeout_action: "error" },
      },
    ],
  },
  {
    name: "Utilities",
    nodes: [
      {
        type: "delay",
        label: "Delay",
        color: "bg-gray-500",
        defaultData: { duration: "1s" },
      },
    ],
  },
  {
    name: "Array Operations",
    nodes: [
      {
        type: "map",
        label: "Map",
        color: "bg-cyan-600",
        defaultData: { expression: "item * 2" },
      },
      {
        type: "reduce",
        label: "Reduce",
        color: "bg-teal-600",
        defaultData: { expression: "acc + item", initial_value: "0" },
      },
      {
        type: "slice",
        label: "Slice",
        color: "bg-emerald-600",
        defaultData: { start: 0, end: -1 },
      },
      {
        type: "sort",
        label: "Sort",
        color: "bg-lime-600",
        defaultData: { field: "", order: "asc" },
      },
      {
        type: "find",
        label: "Find",
        color: "bg-sky-600",
        defaultData: { expression: "item.id == 1" },
      },
      {
        type: "flat_map",
        label: "FlatMap",
        color: "bg-indigo-600",
        defaultData: { expression: "item.values" },
      },
      {
        type: "group_by",
        label: "Group By",
        color: "bg-violet-600",
        defaultData: { key_field: "category" },
      },
      {
        type: "unique",
        label: "Unique",
        color: "bg-fuchsia-600",
        defaultData: { by_field: "" },
      },
      {
        type: "chunk",
        label: "Chunk",
        color: "bg-pink-600",
        defaultData: { size: 3 },
      },
      {
        type: "reverse",
        label: "Reverse",
        color: "bg-rose-600",
        defaultData: {},
      },
      {
        type: "partition",
        label: "Partition",
        color: "bg-orange-600",
        defaultData: { expression: "item > 0" },
      },
      {
        type: "zip",
        label: "Zip",
        color: "bg-yellow-600",
        defaultData: {},
      },
      {
        type: "sample",
        label: "Sample",
        color: "bg-blue-600",
        defaultData: { count: 1 },
      },
      {
        type: "range",
        label: "Range",
        color: "bg-green-600",
        defaultData: { start: 0, end: 10, step: 1 },
      },
      {
        type: "transpose",
        label: "Transpose",
        color: "bg-red-600",
        defaultData: {},
      },
    ],
  },
];
