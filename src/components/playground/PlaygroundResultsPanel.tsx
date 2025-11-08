"use client";
import React, { useState } from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { AlertCircleIcon, PlayIcon, ClockIcon, FileJsonIcon } from "lucide-react";

interface PlaygroundResultsPanelProps {
  isLoading: boolean;
  result: unknown;
  error: string | null;
}

export function PlaygroundResultsPanel({
  isLoading,
  result,
  error,
}: PlaygroundResultsPanelProps) {
  const [activeTab, setActiveTab] = useState<string>("response");

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
    <div className="h-full flex flex-col bg-white dark:bg-black overflow-hidden">
      {/* Header */}
      <div className="h-14 border-b border-gray-300 dark:border-gray-700 flex items-center px-6 gap-4">
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
          <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
            <div className="w-4 h-4 border-2 border-black dark:border-white border-t-transparent rounded-full animate-spin" />
            <span>Executing request...</span>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full p-8">
            <div className="w-16 h-16 border-4 border-black dark:border-white border-t-transparent rounded-full animate-spin mb-6" />
            <p className="text-lg font-medium">Executing HTTP request...</p>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
              Please wait while we process your request
            </p>
          </div>
        )}

        {!isLoading && error && (
          <div className="p-8">
            <div className="max-w-3xl mx-auto border border-gray-300 dark:border-gray-700 rounded-lg p-6">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 bg-black dark:bg-white rounded-full flex items-center justify-center shrink-0">
                  <AlertCircleIcon className="w-5 h-5 text-white dark:text-black" />
                </div>
                <div className="flex-1">
                  <h3 className="text-base font-semibold mb-3">
                    Request Failed
                  </h3>
                  <pre className="text-sm bg-gray-100 dark:bg-gray-900 rounded-md p-4 overflow-x-auto border border-gray-300 dark:border-gray-700">
                    {error}
                  </pre>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && resultData && (
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full h-full flex flex-col">
            <TabsList className="w-full justify-start rounded-none border-b border-gray-300 dark:border-gray-700 bg-transparent px-6 h-auto">
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
                          <div className="flex items-center gap-4 text-xs text-gray-600 dark:text-gray-400">
                            <span>Size: <span className="font-mono font-medium">{JSON.stringify(resultData.data).length}</span> bytes</span>
                            <Separator orientation="vertical" className="h-4" />
                            <span>Type: <span className="font-mono font-medium">application/json</span></span>
                          </div>
                        </div>
                        <div className="rounded-lg border border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-950 overflow-hidden">
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
                    <span className="text-xs text-gray-600 dark:text-gray-400">
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
                        className="rounded-lg border border-gray-300 dark:border-gray-700 p-4"
                      >
                        <div className="flex items-start gap-4">
                          <span className="font-semibold font-mono text-sm min-w-[200px]">
                            {key}
                          </span>
                          <span className="text-sm break-all flex-1">{value}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <p className="text-sm text-gray-600 dark:text-gray-400">No headers available</p>
                  </div>
                )}
              </TabsContent>

              <TabsContent value="timing" className="p-6 space-y-4 mt-0">
                <h3 className="text-sm font-semibold">Request Timing</h3>
                <div className="rounded-lg border border-gray-300 dark:border-gray-700 p-6">
                  <div className="space-y-4 text-sm max-w-lg">
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">Total Duration</span>
                      <span className="font-mono font-semibold">~1.5s (mock)</span>
                    </div>
                    <Separator className="bg-gray-300 dark:bg-gray-700" />
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">DNS Lookup</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">TCP Connection</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">TLS Handshake</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">Time to First Byte</span>
                      <span className="font-mono">-</span>
                    </div>
                    <div className="flex justify-between items-center py-2">
                      <span className="text-gray-600 dark:text-gray-400">Content Download</span>
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
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-full border-2 border-gray-300 dark:border-gray-700 mb-6">
              <PlayIcon className="w-10 h-10 fill-current" />
            </div>
            <p className="text-lg font-semibold mb-2">Ready to execute</p>
            <p className="text-sm text-gray-600 dark:text-gray-400 text-center max-w-md">
              Click the <span className="font-semibold">Run</span> button or press{" "}
              <kbd className="px-1.5 py-0.5 rounded bg-gray-200 dark:bg-gray-800 border border-gray-300 dark:border-gray-700 text-xs">Ctrl+Enter</kbd>{" "}
              to send your HTTP request
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
