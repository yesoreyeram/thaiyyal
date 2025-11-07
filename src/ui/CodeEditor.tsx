import React from "react";

interface CodeEditorProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  language?: "json" | "text" | "javascript";
}

export const CodeEditor = React.forwardRef<
  HTMLTextAreaElement,
  CodeEditorProps
>(({ label, language = "json", className = "", ...props }, ref) => {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm text-gray-400 mb-2">{label}</label>
      )}
      <textarea
        ref={ref}
        className={`w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent ${className}`}
        {...props}
      />
    </div>
  );
});

CodeEditor.displayName = "CodeEditor";
