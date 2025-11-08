"use client";

import React, { createContext, useContext, useState } from "react";
import { cn } from "@/lib/utils";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface SidebarContextValue {
  isCollapsed: boolean;
  toggleCollapsed: () => void;
}

const SidebarContext = createContext<SidebarContextValue | undefined>(
  undefined
);

export function useSidebar() {
  const context = useContext(SidebarContext);
  if (!context) {
    throw new Error("useSidebar must be used within a SidebarProvider");
  }
  return context;
}

interface SidebarProviderProps {
  children: React.ReactNode;
  defaultCollapsed?: boolean;
}

export function SidebarProvider({
  children,
  defaultCollapsed = false,
}: SidebarProviderProps) {
  const [isCollapsed, setIsCollapsed] = useState(defaultCollapsed);

  const toggleCollapsed = () => {
    setIsCollapsed((prev) => !prev);
  };

  return (
    <SidebarContext.Provider value={{ isCollapsed, toggleCollapsed }}>
      {children}
    </SidebarContext.Provider>
  );
}

interface SidebarProps {
  children: React.ReactNode;
  className?: string;
}

export function Sidebar({ children, className }: SidebarProps) {
  const { isCollapsed } = useSidebar();

  return (
    <aside
      className={cn(
        "relative flex flex-col border-r bg-white dark:bg-black border-gray-300 dark:border-gray-700 transition-all duration-300",
        isCollapsed ? "w-16" : "w-64",
        className
      )}
    >
      {children}
    </aside>
  );
}

interface SidebarHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarHeader({ children, className }: SidebarHeaderProps) {
  return (
    <div
      className={cn(
        "flex items-center h-12 border-b border-gray-300 dark:border-gray-700 px-4",
        className
      )}
    >
      {children}
    </div>
  );
}

interface SidebarContentProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarContent({ children, className }: SidebarContentProps) {
  return (
    <div
      className={cn("flex-1 overflow-auto py-2", className)}
    >
      {children}
    </div>
  );
}

interface SidebarFooterProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarFooter({ children, className }: SidebarFooterProps) {
  return (
    <div
      className={cn(
        "border-t border-gray-300 dark:border-gray-700 p-4",
        className
      )}
    >
      {children}
    </div>
  );
}

interface SidebarGroupProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarGroup({ children, className }: SidebarGroupProps) {
  return <div className={cn("px-2 py-2", className)}>{children}</div>;
}

interface SidebarGroupLabelProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarGroupLabel({
  children,
  className,
}: SidebarGroupLabelProps) {
  const { isCollapsed } = useSidebar();

  if (isCollapsed) return null;

  return (
    <div
      className={cn(
        "px-3 py-2 text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wider",
        className
      )}
    >
      {children}
    </div>
  );
}

interface SidebarGroupContentProps {
  children: React.ReactNode;
  className?: string;
}

export function SidebarGroupContent({
  children,
  className,
}: SidebarGroupContentProps) {
  return <div className={cn("space-y-1", className)}>{children}</div>;
}

interface SidebarMenuItemProps {
  icon: React.ReactNode;
  label: string;
  onClick?: () => void;
  active?: boolean;
  className?: string;
}

export function SidebarMenuItem({
  icon,
  label,
  onClick,
  active = false,
  className,
}: SidebarMenuItemProps) {
  const { isCollapsed } = useSidebar();

  return (
    <button
      onClick={onClick}
      className={cn(
        "w-full flex items-center gap-3 px-3 py-2 rounded-md transition-colors",
        "hover:bg-gray-100 dark:hover:bg-gray-900",
        active &&
          "bg-gray-200 dark:bg-gray-800 text-black dark:text-white font-medium",
        !active && "text-gray-700 dark:text-gray-300",
        isCollapsed && "justify-center px-2",
        className
      )}
      title={isCollapsed ? label : undefined}
    >
      <span className="flex-shrink-0">{icon}</span>
      {!isCollapsed && <span className="text-sm">{label}</span>}
    </button>
  );
}

export function SidebarToggle() {
  const { isCollapsed, toggleCollapsed } = useSidebar();

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={toggleCollapsed}
      className="absolute -right-3 top-12 z-10 h-6 w-6 rounded-full border bg-white dark:bg-black border-gray-300 dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-900"
      aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
    >
      {isCollapsed ? (
        <ChevronRight className="h-4 w-4" />
      ) : (
        <ChevronLeft className="h-4 w-4" />
      )}
    </Button>
  );
}
