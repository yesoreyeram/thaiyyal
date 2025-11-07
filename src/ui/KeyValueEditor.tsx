import React from "react";

interface KeyValuePair {
  key: string;
  value: string;
}

interface KeyValueEditorProps {
  items: KeyValuePair[];
  onChange: (items: KeyValuePair[]) => void;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  emptyMessage?: string;
}

export function KeyValueEditor({
  items,
  onChange,
  keyPlaceholder = "Key",
  valuePlaceholder = "Value",
  emptyMessage = "No items. Click + to add one.",
}: KeyValueEditorProps) {
  const addItem = () => {
    onChange([...items, { key: "", value: "" }]);
  };

  const updateItem = (index: number, key: string, value: string) => {
    const newItems = [...items];
    newItems[index] = { key, value };
    onChange(newItems);
  };

  const removeItem = (index: number) => {
    onChange(items.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <button
          onClick={addItem}
          className="text-xs text-blue-400 hover:text-blue-300 px-2 py-1 rounded hover:bg-gray-800"
        >
          + Add
        </button>
      </div>
      {items.length === 0 ? (
        <div className="text-sm text-gray-500 text-center py-8">
          {emptyMessage}
        </div>
      ) : (
        <div className="space-y-2">
          {items.map((item, index) => (
            <div key={index} className="flex gap-2 items-center">
              <input
                type="text"
                value={item.key}
                onChange={(e) => updateItem(index, e.target.value, item.value)}
                placeholder={keyPlaceholder}
                className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <input
                type="text"
                value={item.value}
                onChange={(e) => updateItem(index, item.key, e.target.value)}
                placeholder={valuePlaceholder}
                className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <button
                onClick={() => removeItem(index)}
                className="px-3 py-2 text-red-400 hover:text-red-300 hover:bg-gray-800 rounded"
                title="Remove"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
