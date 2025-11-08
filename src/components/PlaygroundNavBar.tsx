"use client";
import React, { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PlayIcon, SaveIcon, DownloadIcon, UploadIcon, PencilIcon, CheckIcon, XIcon } from "lucide-react";

interface PlaygroundNavBarProps {
  requestTitle: string;
  onTitleChange: (title: string) => void;
  onRun: () => void;
  isRunning: boolean;
  onSave?: () => void;
  onExport?: () => void;
  onImport?: (data: unknown) => void;
}

export function PlaygroundNavBar({
  requestTitle,
  onTitleChange,
  onRun,
  isRunning,
  onSave,
  onExport,
  onImport,
}: PlaygroundNavBarProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(requestTitle);
  const inputRef = useRef<HTMLInputElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  useEffect(() => {
    setEditValue(requestTitle);
  }, [requestTitle]);

  const handleSubmit = () => {
    if (editValue.trim()) {
      onTitleChange(editValue.trim());
    } else {
      setEditValue(requestTitle);
    }
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditValue(requestTitle);
    setIsEditing(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSubmit();
    } else if (e.key === "Escape") {
      handleCancel();
    }
  };

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const json = JSON.parse(event.target?.result as string);
        onImport?.(json);
      } catch {
        alert(
          "Failed to parse JSON file. Please ensure it is a valid JSON file."
        );
      }
    };
    reader.readAsText(file);

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <div className="h-14 border-b bg-card flex items-center justify-between px-6 shadow-sm">
      <div className="flex items-center gap-4">
        {isEditing ? (
          <div className="flex items-center gap-2">
            <Input
              ref={inputRef}
              type="text"
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              onKeyDown={handleKeyDown}
              className="h-9 w-64 font-medium"
              placeholder="Request name"
            />
            <Button
              variant="ghost"
              size="icon"
              onClick={handleSubmit}
              className="h-9 w-9"
            >
              <CheckIcon className="h-4 w-4 text-emerald-600" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleCancel}
              className="h-9 w-9"
            >
              <XIcon className="h-4 w-4 text-destructive" />
            </Button>
          </div>
        ) : (
          <button
            onClick={() => setIsEditing(true)}
            className="flex items-center gap-2 px-3 py-1.5 hover:bg-accent rounded-md transition-colors group"
          >
            <span className="text-base font-semibold">{requestTitle}</span>
            <PencilIcon className="h-3.5 w-3.5 opacity-0 group-hover:opacity-60 transition-opacity" />
          </button>
        )}
        <div className="h-6 w-px bg-border" />
        <span className="text-xs text-muted-foreground">
          Press <kbd className="px-1.5 py-0.5 rounded bg-muted border text-xs">Ctrl+Enter</kbd> to run
        </span>
      </div>

      <div className="flex items-center gap-2">
        {onSave && (
          <Button variant="outline" size="sm" onClick={onSave}>
            <SaveIcon className="h-4 w-4 mr-2" />
            Save
          </Button>
        )}

        {onExport && (
          <Button variant="outline" size="sm" onClick={onExport}>
            <DownloadIcon className="h-4 w-4 mr-2" />
            Export
          </Button>
        )}

        {onImport && (
          <>
            <Button variant="outline" size="sm" onClick={handleImportClick}>
              <UploadIcon className="h-4 w-4 mr-2" />
              Import
            </Button>
            <input
              ref={fileInputRef}
              type="file"
              accept=".json,application/json"
              onChange={handleFileChange}
              className="hidden"
            />
          </>
        )}

        <div className="h-6 w-px bg-border ml-2" />

        <Button
          onClick={onRun}
          disabled={isRunning}
          size="default"
          className="font-semibold"
        >
          {isRunning ? (
            <>
              <span className="h-4 w-4 mr-2 border-2 border-white border-t-transparent rounded-full animate-spin" />
              Running...
            </>
          ) : (
            <>
              <PlayIcon className="h-4 w-4 mr-2 fill-current" />
              Run
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
