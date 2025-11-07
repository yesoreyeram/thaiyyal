"use client";
import React, { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PlayIcon, SaveIcon, DownloadIcon, UploadIcon, Edit2Icon } from "lucide-react";

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

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSubmit();
    } else if (e.key === "Escape") {
      setEditValue(requestTitle);
      setIsEditing(false);
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
    <div className="h-12 border-b flex items-center justify-between px-4">
      <div className="flex items-center gap-3">
        {isEditing ? (
          <Input
            ref={inputRef}
            type="text"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            onBlur={handleSubmit}
            onKeyDown={handleKeyDown}
            className="h-8 min-w-[200px]"
          />
        ) : (
          <button
            onClick={() => setIsEditing(true)}
            className="flex items-center gap-2 px-2 py-1 hover:bg-accent rounded transition-colors group"
          >
            <span className="text-sm font-medium">{requestTitle}</span>
            <Edit2Icon className="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity" />
          </button>
        )}
      </div>

      <div className="flex items-center gap-2">
        {onSave && (
          <Button variant="outline" size="sm" onClick={onSave}>
            <SaveIcon className="h-4 w-4 mr-1" />
            Save
          </Button>
        )}

        {onExport && (
          <Button variant="outline" size="sm" onClick={onExport}>
            <DownloadIcon className="h-4 w-4 mr-1" />
            Export
          </Button>
        )}

        {onImport && (
          <>
            <Button variant="outline" size="sm" onClick={handleImportClick}>
              <UploadIcon className="h-4 w-4 mr-1" />
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

        <Button
          onClick={onRun}
          disabled={isRunning}
          size="sm"
        >
          {isRunning ? (
            <>
              <span className="h-4 w-4 mr-1 border-2 border-white border-t-transparent rounded-full animate-spin" />
              Running...
            </>
          ) : (
            <>
              <PlayIcon className="h-4 w-4 mr-1" />
              Run
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
