"use client";
import React, { useState, useRef, useEffect } from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

interface PlaygroundResultsPanelProps {
  isLoading: boolean;
  result: unknown;
  error: string | null;
  height: number;
  onHeightChange: (height: number) => void;
}

export function PlaygroundResultsPanel({
  isLoading,
  result,
  error,
  height,
  onHeightChange,
}: PlaygroundResultsPanelProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [activeTab, setActiveTab] = useState<string>("response");
  const dragStartY = useRef(0);
  const dragStartHeight = useRef(0);

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = dragStartY.current - e.clientY;
      const newHeight = Math.max(
        100,
        Math.min(600, dragStartHeight.current + deltaY)
      );
      onHeightChange(newHeight);
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isDragging, onHeightChange]);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    dragStartY.current = e.clientY;
    dragStartHeight.current = height;
  };

  const resultData = result as {
    status?: number;
    statusText?: string;
    headers?: Record<string, string>;
    data?: unknown;
  } | null;

  return (
    <div
      className="border-t flex flex-col bg-background"
      style={{ height: `${height}px` }}
    >
      {/* Resize Handle */}
      <div
        className={`h-1 hover:bg-primary cursor-ns-resize transition-colors ${
          isDragging ? "bg-primary" : "bg-border"
        }`}
        onMouseDown={handleMouseDown}
      />

      {/* Header */}
      <div className="h-10 border-b flex items-center px-4 gap-2">
        <span className="text-sm font-semibold">Response</span>
        {!isLoading && resultData && (
          <>
            <Badge
              variant={
                resultData.status && resultData.status < 400
                  ? "default"
                  : "destructive"
              }
            >
              {resultData.status} {resultData.statusText}
            </Badge>
          </>
        )}
        {isLoading && (
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <div className="w-3 h-3 border-2 border-primary border-t-transparent rounded-full animate-spin" />
            <span>Executing request...</span>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
            <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin mb-4" />
            <p className="text-sm">Executing HTTP request...</p>
            <p className="text-xs text-muted-foreground mt-2">
              This may take a few moments
            </p>
          </div>
        )}

        {!isLoading && error && (
          <div className="p-4">
            <div className="bg-destructive/10 border border-destructive rounded-lg p-4">
              <div className="flex items-start gap-3">
                <div className="w-5 h-5 bg-destructive rounded-full flex items-center justify-center shrink-0 mt-0.5">
                  <span className="text-destructive-foreground text-xs font-bold">
                    !
                  </span>
                </div>
                <div className="flex-1">
                  <h3 className="text-sm font-semibold text-destructive mb-2">
                    Request Error
                  </h3>
                  <pre className="text-xs bg-muted rounded p-3 overflow-x-auto">
                    {error}
                  </pre>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && resultData && (
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
            <TabsList className="w-full justify-start rounded-none border-b bg-transparent px-4">
              <TabsTrigger value="response">Response</TabsTrigger>
              <TabsTrigger value="headers">Headers</TabsTrigger>
              <TabsTrigger value="preview">Preview</TabsTrigger>
              <TabsTrigger value="timing">Timing</TabsTrigger>
            </TabsList>

            <TabsContent value="response" className="p-4 space-y-2 mt-0">
              {(() => {
                try {
                  const stringifiedData = JSON.stringify(
                    resultData.data,
                    null,
                    2
                  );
                  return (
                    <>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-xs text-muted-foreground">
                          Response Body
                        </span>
                        <span className="text-xs text-muted-foreground">
                          Size: {JSON.stringify(resultData.data).length} bytes
                        </span>
                      </div>
                      <pre className="text-xs bg-muted rounded p-3 overflow-x-auto font-mono border">
                        {stringifiedData}
                      </pre>
                    </>
                  );
                } catch (error) {
                  return (
                    <pre className="text-xs text-destructive bg-muted rounded p-3 overflow-x-auto">
                      Error formatting response:{" "}
                      {error instanceof Error ? error.message : "Unknown error"}
                    </pre>
                  );
                }
              })()}
            </TabsContent>

            <TabsContent value="headers" className="p-4 space-y-2 mt-0">
              <div className="text-xs text-muted-foreground mb-2">
                Response Headers
              </div>
              {resultData.headers &&
              Object.keys(resultData.headers).length > 0 ? (
                <div className="space-y-1">
                  {Object.entries(resultData.headers).map(([key, value]) => (
                    <div
                      key={key}
                      className="bg-muted rounded p-2 border flex gap-3 text-xs"
                    >
                      <span className="font-semibold text-primary min-w-[150px]">
                        {key}:
                      </span>
                      <span className="break-all">{value}</span>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-xs text-muted-foreground text-center py-4">
                  No headers available
                </div>
              )}
            </TabsContent>

            <TabsContent value="preview" className="p-4 space-y-2 mt-0">
              <div className="text-xs text-muted-foreground mb-2">Preview</div>
              <div className="bg-muted rounded p-3 border">
                <div className="text-xs">
                  <p>
                    Status: {resultData.status} {resultData.statusText}
                  </p>
                  <p className="mt-2 text-muted-foreground">
                    Preview functionality will render formatted response based on
                    content type
                  </p>
                </div>
              </div>
            </TabsContent>

            <TabsContent value="timing" className="p-4 space-y-2 mt-0">
              <div className="text-xs text-muted-foreground mb-2">Timing</div>
              <div className="bg-muted rounded p-3 border">
                <div className="space-y-2 text-xs">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Total Duration:
                    </span>
                    <span>~1.5s (mock)</span>
                  </div>
                  <Separator />
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">DNS Lookup:</span>
                    <span>-</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      TCP Connection:
                    </span>
                    <span>-</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      TLS Handshake:
                    </span>
                    <span>-</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Time to First Byte:
                    </span>
                    <span>-</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Content Download:
                    </span>
                    <span>-</span>
                  </div>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        )}

        {!isLoading && !resultData && !error && (
          <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
            <div className="text-4xl mb-4">▶️</div>
            <p className="text-sm">
              Click &ldquo;Run&rdquo; to execute the HTTP request
            </p>
            <p className="text-xs text-muted-foreground mt-2">
              Results will appear here
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
