"use client";
import React, { useState, useRef, useEffect } from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { AlertCircleIcon, PlayIcon, ClockIcon, FileJsonIcon } from "lucide-react";

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
        150,
        Math.min(700, dragStartHeight.current + deltaY)
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

  const getStatusColor = (status?: number) => {
    if (!status) return "bg-muted text-muted-foreground";
    if (status < 200) return "bg-blue-100 text-blue-800 dark:bg-blue-950 dark:text-blue-300";
    if (status < 300) return "bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300";
    if (status < 400) return "bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-300";
    if (status < 500) return "bg-orange-100 text-orange-800 dark:bg-orange-950 dark:text-orange-300";
    return "bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-300";
  };

  return (
    <div
      className="border-t flex flex-col bg-card shadow-lg"
      style={{ height: `${height}px` }}
    >
      {/* Resize Handle */}
      <div
        className={`h-1.5 cursor-ns-resize transition-colors flex items-center justify-center group ${
          isDragging ? "bg-primary" : "bg-border hover:bg-primary/50"
        }`}
        onMouseDown={handleMouseDown}
      >
        <div className="w-12 h-1 rounded-full bg-muted-foreground/20 group-hover:bg-primary/50" />
      </div>

      {/* Header */}
      <div className="h-14 border-b bg-background flex items-center px-6 gap-4">
        <div className="flex items-center gap-3">
          <h2 className="text-base font-semibold">Response</h2>
          {!isLoading && resultData && (
            <Badge
              className={`${getStatusColor(resultData.status)} border-0 font-semibold px-3`}
            >
              {resultData.status} {resultData.statusText}
            </Badge>
          )}
        </div>
        {isLoading && (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />
            <span>Executing request...</span>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto bg-background">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full text-muted-foreground p-8">
            <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mb-6" />
            <p className="text-lg font-medium">Executing HTTP request...</p>
            <p className="text-sm text-muted-foreground mt-2">
              Please wait while we process your request
            </p>
          </div>
        )}

        {!isLoading && error && (
          <div className="p-8">
            <div className="max-w-3xl mx-auto bg-destructive/10 border border-destructive/30 rounded-lg p-6">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 bg-destructive rounded-full flex items-center justify-center shrink-0">
                  <AlertCircleIcon className="w-5 h-5 text-destructive-foreground" />
                </div>
                <div className="flex-1">
                  <h3 className="text-base font-semibold text-destructive mb-3">
                    Request Failed
                  </h3>
                  <pre className="text-sm bg-background rounded-md p-4 overflow-x-auto border">
                    {error}
                  </pre>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && resultData && (
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full h-full flex flex-col">
            <TabsList className="w-full justify-start rounded-none border-b bg-transparent px-6 h-auto">
              <TabsTrigger value="response" className="data-[state=active]:shadow-none gap-2">
                <FileJsonIcon className="w-4 h-4" />
                Response
              </TabsTrigger>
              <TabsTrigger value="headers" className="data-[state=active]:shadow-none gap-2">
                Headers
                {resultData.headers && (
                  <Badge variant="secondary" className="ml-1 h-5">
                    {Object.keys(resultData.headers).length}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger value="preview" className="data-[state=active]:shadow-none">
                Preview
              </TabsTrigger>
              <TabsTrigger value="timing" className="data-[state=active]:shadow-none gap-2">
                <ClockIcon className="w-4 h-4" />
                Timing
              </TabsTrigger>
            </TabsList>

            <div className="flex-1 overflow-auto">
              <TabsContent value="response" className="p-6 space-y-4 mt-0">
                {(() => {
                  try {
                    const stringifiedData = JSON.stringify(
                      resultData.data,
                      null,
                      2
                    );
                    return (
                      <>
                        <div className="flex items-center justify-between">
                          <h3 className="text-sm font-semibold">Response Body</h3>
                          <div className="flex items-center gap-4 text-xs text-muted-foreground">
                            <span>Size: <span className="font-mono font-medium">{JSON.stringify(resultData.data).length}</span> bytes</span>
                            <Separator orientation="vertical" className="h-4" />
                            <span>Type: <span className="font-mono font-medium">application/json</span></span>
                          </div>
                        </div>
                        <div className="rounded-lg border bg-muted/30 overflow-hidden">
                          <pre className="text-sm p-5 overflow-x-auto font-mono leading-relaxed">
                            {stringifiedData}
                          </pre>
                        </div>
                      </>
                    );
                  } catch (error) {
                    return (
                      <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-4">
                        <pre className="text-sm text-destructive">
                          Error formatting response:{" "}
                          {error instanceof Error ? error.message : "Unknown error"}
                        </pre>
                      </div>
                    );
                  }
                })()}
              </TabsContent>

              <TabsContent value="headers" className="p-6 space-y-4 mt-0">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold">Response Headers</h3>
                  {resultData.headers && (
                    <span className="text-xs text-muted-foreground">
                      {Object.keys(resultData.headers).length} headers
                    </span>
                  )}
                </div>
                {resultData.headers &&
                Object.keys(resultData.headers).length > 0 ? (
                  <div className="space-y-2">
                    {Object.entries(resultData.headers).map(([key, value]) => (
                      <div
                        key={key}
                        className="rounded-lg border bg-card p-4 hover:bg-accent/5 transition-colors"
                      >
                        <div className="flex items-start gap-4">
                          <span className="font-semibold text-primary font-mono text-sm min-w-[200px]">
                            {key}
                          </span>
                          <span className="text-sm break-all flex-1">{value}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <p className="text-sm text-muted-foreground">No headers available</p>
                  </div>
                )}
              </TabsContent>

              <TabsContent value="preview" className="p-6 space-y-4 mt-0">
                <h3 className="text-sm font-semibold">Response Preview</h3>
                <div className="rounded-lg border bg-card p-6">
                  <div className="space-y-3 text-sm">
                    <div className="flex items-center gap-3">
                      <span className="text-muted-foreground w-24">Status:</span>
                      <Badge className={getStatusColor(resultData.status)}>
                        {resultData.status} {resultData.statusText}
                      </Badge>
                    </div>
                    <Separator />
                    <p className="text-muted-foreground">
                      Preview functionality will render formatted response based on content type
                      (HTML, images, etc.)
                    </p>
                  </div>
                </div>
              </TabsContent>

              <TabsContent value="timing" className="p-6 space-y-4 mt-0">
                <h3 className="text-sm font-semibold">Request Timing</h3>
                <div className="rounded-lg border bg-card p-6">
                  <div className="space-y-4 text-sm max-w-lg">
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">Total Duration</span>
                      <span className="font-mono font-semibold">~1.5s (mock)</span>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">DNS Lookup</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">TCP Connection</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">TLS Handshake</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">Time to First Byte</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-muted-foreground">Content Download</span>
                      <span className="font-mono">-</span>
                    </div>
                  </div>
                </div>
              </TabsContent>
            </div>
          </Tabs>
        )}

        {!isLoading && !resultData && !error && (
          <div className="flex flex-col items-center justify-center h-full p-8">
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-6">
              <PlayIcon className="w-10 h-10 text-primary fill-current" />
            </div>
            <p className="text-lg font-semibold mb-2">Ready to execute</p>
            <p className="text-sm text-muted-foreground text-center max-w-md">
              Click the <span className="font-semibold">Run</span> button or press{" "}
              <kbd className="px-1.5 py-0.5 rounded bg-muted border text-xs">Ctrl+Enter</kbd>{" "}
              to send your HTTP request
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
