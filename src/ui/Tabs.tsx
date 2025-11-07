import React from "react";

interface Tab {
  id: string;
  label: string;
  count?: number;
  hidden?: boolean;
}

interface TabsProps {
  tabs: Tab[];
  activeTab: string;
  onChange: (tabId: string) => void;
  className?: string;
}

export function Tabs({ tabs, activeTab, onChange, className = "" }: TabsProps) {
  return (
    <div
      className={`flex items-center gap-1 border-b border-gray-800 bg-gray-950 ${className}`}
    >
      {tabs.map(
        (tab) =>
          !tab.hidden && (
            <button
              key={tab.id}
              onClick={() => onChange(tab.id)}
              className={`px-3 py-2 text-sm rounded-t transition-colors ${
                activeTab === tab.id
                  ? "bg-gray-900 text-white border-t border-l border-r border-gray-700"
                  : "text-gray-400 hover:text-white hover:bg-gray-900/50"
              }`}
            >
              {tab.label}
              {tab.count !== undefined && tab.count > 0 && (
                <span className="ml-1.5 text-xs bg-gray-700 px-1.5 py-0.5 rounded">
                  {tab.count}
                </span>
              )}
            </button>
          )
      )}
    </div>
  );
}
