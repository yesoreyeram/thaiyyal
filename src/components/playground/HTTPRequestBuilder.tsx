"use client";
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { PlusIcon, XIcon } from "lucide-react";

interface HTTPRequestBuilderProps {
  baseUrl: string;
  onBaseUrlChange: (value: string) => void;
  authMethod: "none" | "basic" | "bearer";
  onAuthMethodChange: (value: "none" | "basic" | "bearer") => void;
  authUsername: string;
  onAuthUsernameChange: (value: string) => void;
  authPassword: string;
  onAuthPasswordChange: (value: string) => void;
  bearerToken: string;
  onBearerTokenChange: (value: string) => void;
  timeout: string;
  onTimeoutChange: (value: string) => void;
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  onMethodChange: (value: "GET" | "POST" | "PUT" | "PATCH" | "DELETE") => void;
  url: string;
  onUrlChange: (value: string) => void;
  body: string;
  onBodyChange: (value: string) => void;
  headers: Array<{ key: string; value: string }>;
  onHeadersChange: (headers: Array<{ key: string; value: string }>) => void;
  queryParams: Array<{ key: string; value: string }>;
  onQueryParamsChange: (
    params: Array<{ key: string; value: string }>
  ) => void;
  paginationEnabled: boolean;
  onPaginationEnabledChange: (enabled: boolean) => void;
  pageSize: string;
  onPageSizeChange: (value: string) => void;
  pageParam: string;
  onPageParamChange: (value: string) => void;
  parseResponse: boolean;
  onParseResponseChange: (enabled: boolean) => void;
  jsonPath: string;
  onJsonPathChange: (value: string) => void;
  onRun: () => void;
}

export function HTTPRequestBuilder(props: HTTPRequestBuilderProps) {
  const [activeTab, setActiveTab] = useState<string>("query");

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
      e.preventDefault();
      props.onRun();
    }
  };

  const fullUrl =
    props.baseUrl +
    props.url +
    (props.queryParams.length > 0 &&
    props.queryParams.some((p) => p.key.trim() !== "")
      ? "?" +
        props.queryParams
          .filter((p) => p.key.trim() !== "")
          .map(
            (p) =>
              `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`
          )
          .join("&")
      : "");

  return (
    <div className="h-full flex flex-col" onKeyDown={handleKeyDown}>
      {/* Method and URL Row */}
      <div className="p-4 border-b">
        <div className="flex gap-3 items-center">
          <Select
            value={props.method}
            onChange={(e) =>
              props.onMethodChange(
                e.target.value as "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
              )
            }
            className="w-32"
          >
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="PATCH">PATCH</option>
            <option value="DELETE">DELETE</option>
          </Select>
          <div className="flex-1 flex flex-col gap-1">
            <Input
              type="text"
              value={props.baseUrl + props.url}
              onChange={(e) => {
                const value = e.target.value;
                try {
                  if (
                    value.startsWith("http://") ||
                    value.startsWith("https://")
                  ) {
                    const urlObj = new URL(value);
                    props.onBaseUrlChange(urlObj.origin);
                    props.onUrlChange(urlObj.pathname + urlObj.search);
                  } else {
                    props.onUrlChange(value);
                  }
                } catch {
                  props.onUrlChange(value);
                }
              }}
              placeholder="https://api.example.com/endpoint"
            />
            {fullUrl && fullUrl !== props.baseUrl + props.url && (
              <div className="text-xs text-muted-foreground px-1">
                Full URL: {fullUrl}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
        <TabsList className="w-full justify-start rounded-none border-b bg-transparent px-4">
          <TabsTrigger value="query" className="relative">
            Query Params
            {props.queryParams.length > 0 && (
              <Badge variant="secondary" className="ml-1.5 text-xs">
                {props.queryParams.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="headers" className="relative">
            Headers
            {props.headers.length > 0 && (
              <Badge variant="secondary" className="ml-1.5 text-xs">
                {props.headers.length}
              </Badge>
            )}
          </TabsTrigger>
          {props.method !== "GET" && props.method !== "DELETE" && (
            <TabsTrigger value="body">Body</TabsTrigger>
          )}
          <TabsTrigger value="auth">Auth</TabsTrigger>
          <TabsTrigger value="client">HTTP Client</TabsTrigger>
          <TabsTrigger value="pagination">Pagination</TabsTrigger>
          <TabsTrigger value="parsing">Response Parsing</TabsTrigger>
        </TabsList>

        <div className="flex-1 overflow-auto p-4">
          {/* Query Params Tab */}
          <TabsContent value="query" className="mt-0 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-medium">Query Parameters</h3>
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  props.onQueryParamsChange([
                    ...props.queryParams,
                    { key: "", value: "" },
                  ])
                }
              >
                <PlusIcon className="h-4 w-4 mr-1" />
                Add
              </Button>
            </div>
            {props.queryParams.length === 0 ? (
              <div className="text-sm text-muted-foreground text-center py-8">
                No query parameters. Click &quot;Add&quot; to add one.
              </div>
            ) : (
              <div className="space-y-2">
                {props.queryParams.map((param, index) => (
                  <div key={index} className="flex gap-2 items-center">
                    <Input
                      type="text"
                      value={param.key}
                      onChange={(e) => {
                        const newParams = [...props.queryParams];
                        newParams[index] = { key: e.target.value, value: param.value };
                        props.onQueryParamsChange(newParams);
                      }}
                      placeholder="Key"
                      className="flex-1"
                    />
                    <Input
                      type="text"
                      value={param.value}
                      onChange={(e) => {
                        const newParams = [...props.queryParams];
                        newParams[index] = { key: param.key, value: e.target.value };
                        props.onQueryParamsChange(newParams);
                      }}
                      placeholder="Value"
                      className="flex-1"
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        props.onQueryParamsChange(
                          props.queryParams.filter((_, i) => i !== index)
                        );
                      }}
                    >
                      <XIcon className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Headers Tab */}
          <TabsContent value="headers" className="mt-0 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-medium">Headers</h3>
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  props.onHeadersChange([
                    ...props.headers,
                    { key: "", value: "" },
                  ])
                }
              >
                <PlusIcon className="h-4 w-4 mr-1" />
                Add
              </Button>
            </div>
            {props.headers.length === 0 ? (
              <div className="text-sm text-muted-foreground text-center py-8">
                No headers. Click &quot;Add&quot; to add one.
              </div>
            ) : (
              <div className="space-y-2">
                {props.headers.map((header, index) => (
                  <div key={index} className="flex gap-2 items-center">
                    <Input
                      type="text"
                      value={header.key}
                      onChange={(e) => {
                        const newHeaders = [...props.headers];
                        newHeaders[index] = { key: e.target.value, value: header.value };
                        props.onHeadersChange(newHeaders);
                      }}
                      placeholder="Header name"
                      className="flex-1"
                    />
                    <Input
                      type="text"
                      value={header.value}
                      onChange={(e) => {
                        const newHeaders = [...props.headers];
                        newHeaders[index] = { key: header.key, value: e.target.value };
                        props.onHeadersChange(newHeaders);
                      }}
                      placeholder="Header value"
                      className="flex-1"
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        props.onHeadersChange(
                          props.headers.filter((_, i) => i !== index)
                        );
                      }}
                    >
                      <XIcon className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Body Tab */}
          <TabsContent value="body" className="mt-0 space-y-4">
            <h3 className="text-sm font-medium">Request Body</h3>
            <Textarea
              value={props.body}
              onChange={(e) => props.onBodyChange(e.target.value)}
              placeholder='{"key": "value"}'
              rows={15}
              className="font-mono text-sm"
            />
          </TabsContent>

          {/* Auth Tab */}
          <TabsContent value="auth" className="mt-0 space-y-4">
            <h3 className="text-sm font-medium">Authentication</h3>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Auth Type</Label>
                <Select
                  value={props.authMethod}
                  onChange={(e) =>
                    props.onAuthMethodChange(
                      e.target.value as "none" | "basic" | "bearer"
                    )
                  }
                >
                  <option value="none">No Authentication</option>
                  <option value="basic">Basic Auth</option>
                  <option value="bearer">Bearer Token</option>
                </Select>
              </div>

              {props.authMethod === "basic" && (
                <div className="space-y-4 pl-4 border-l-2 border-primary">
                  <div className="space-y-2">
                    <Label>Username</Label>
                    <Input
                      type="text"
                      value={props.authUsername}
                      onChange={(e) => props.onAuthUsernameChange(e.target.value)}
                      placeholder="username"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Password</Label>
                    <Input
                      type="password"
                      value={props.authPassword}
                      onChange={(e) => props.onAuthPasswordChange(e.target.value)}
                      placeholder="password"
                    />
                  </div>
                </div>
              )}

              {props.authMethod === "bearer" && (
                <div className="pl-4 border-l-2 border-primary">
                  <div className="space-y-2">
                    <Label>Bearer Token</Label>
                    <Input
                      type="password"
                      value={props.bearerToken}
                      onChange={(e) => props.onBearerTokenChange(e.target.value)}
                      placeholder="your-token-here"
                    />
                  </div>
                </div>
              )}
            </div>
          </TabsContent>

          {/* HTTP Client Tab */}
          <TabsContent value="client" className="mt-0 space-y-4">
            <h3 className="text-sm font-medium">HTTP Client Settings</h3>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>Base URL</Label>
                <Input
                  type="text"
                  value={props.baseUrl}
                  onChange={(e) => props.onBaseUrlChange(e.target.value)}
                  placeholder="https://api.example.com"
                />
              </div>
              <div className="space-y-2">
                <Label>Timeout (seconds)</Label>
                <Input
                  type="number"
                  value={props.timeout}
                  onChange={(e) => props.onTimeoutChange(e.target.value)}
                  placeholder="30"
                  min="1"
                  max="300"
                />
              </div>
            </div>
          </TabsContent>

          {/* Pagination Tab */}
          <TabsContent value="pagination" className="mt-0 space-y-4">
            <h3 className="text-sm font-medium">Pagination Options</h3>
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="pagination-enabled"
                  checked={props.paginationEnabled}
                  onChange={(e) =>
                    props.onPaginationEnabledChange(e.target.checked)
                  }
                  className="h-4 w-4"
                />
                <Label htmlFor="pagination-enabled">Enable Pagination</Label>
              </div>

              {props.paginationEnabled && (
                <div className="space-y-4 pl-4 border-l-2 border-primary">
                  <div className="space-y-2">
                    <Label>Page Size</Label>
                    <Input
                      type="number"
                      value={props.pageSize}
                      onChange={(e) => props.onPageSizeChange(e.target.value)}
                      placeholder="10"
                      min="1"
                      max="1000"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Page Parameter Name</Label>
                    <Input
                      type="text"
                      value={props.pageParam}
                      onChange={(e) => props.onPageParamChange(e.target.value)}
                      placeholder="page"
                    />
                  </div>
                </div>
              )}
            </div>
          </TabsContent>

          {/* Response Parsing Tab */}
          <TabsContent value="parsing" className="mt-0 space-y-4">
            <h3 className="text-sm font-medium">Response Parsing</h3>
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="parse-response"
                  checked={props.parseResponse}
                  onChange={(e) => props.onParseResponseChange(e.target.checked)}
                  className="h-4 w-4"
                />
                <Label htmlFor="parse-response">Parse JSON Response</Label>
              </div>

              {props.parseResponse && (
                <div className="pl-4 border-l-2 border-primary">
                  <div className="space-y-2">
                    <Label>JSON Path (optional)</Label>
                    <Input
                      type="text"
                      value={props.jsonPath}
                      onChange={(e) => props.onJsonPathChange(e.target.value)}
                      placeholder="$.data.items"
                      className="font-mono"
                    />
                    <p className="text-xs text-muted-foreground">
                      Extract specific data from response using JSON path notation
                    </p>
                  </div>
                </div>
              )}
            </div>
          </TabsContent>
        </div>
      </Tabs>
    </div>
  );
}
