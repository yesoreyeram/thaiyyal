import React from "react";

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  helperText?: string;
}

export const Checkbox = React.forwardRef<HTMLInputElement, CheckboxProps>(
  ({ label, helperText, className = "", ...props }, ref) => {
    return (
      <div className="flex items-start gap-2">
        <input
          ref={ref}
          type="checkbox"
          className={`w-4 h-4 mt-0.5 text-blue-600 bg-gray-900 border-gray-700 rounded focus:ring-blue-500 focus:ring-2 ${className}`}
          {...props}
        />
        {label && (
          <div className="flex-1">
            <label className="text-sm text-gray-300 cursor-pointer">
              {label}
            </label>
            {helperText && (
              <p className="mt-0.5 text-xs text-gray-500">{helperText}</p>
            )}
          </div>
        )}
      </div>
    );
  }
);

Checkbox.displayName = "Checkbox";
