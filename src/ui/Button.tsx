import React from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md" | "lg";
  isLoading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      variant = "primary",
      size = "md",
      isLoading = false,
      className = "",
      children,
      disabled,
      ...props
    },
    ref
  ) => {
    const baseStyles =
      "rounded font-medium transition-all flex items-center justify-center gap-1.5";

    const variantStyles = {
      primary:
        "bg-blue-600 hover:bg-blue-700 text-white shadow-sm hover:shadow-md",
      secondary:
        "bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white",
      ghost: "bg-transparent hover:bg-gray-800 text-gray-400 hover:text-white",
      danger: "bg-red-600 hover:bg-red-700 text-white shadow-sm hover:shadow-md",
    };

    const sizeStyles = {
      sm: "px-2 py-1 text-xs",
      md: "px-3 py-1.5 text-sm",
      lg: "px-4 py-2 text-base",
    };

    const disabledStyles = "bg-gray-600 cursor-not-allowed opacity-50";

    return (
      <button
        ref={ref}
        className={`${baseStyles} ${
          disabled || isLoading ? disabledStyles : variantStyles[variant]
        } ${sizeStyles[size]} ${className}`}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading && (
          <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
