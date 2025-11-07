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
import { PlusIcon, XIcon, InfoIcon } from "lucide-react";

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

const methodColors = {
  GET: "bg-emerald-600 hover:bg-emerald-700",
  POST: "bg-blue-600 hover:bg-blue-700",
  PUT: "bg-amber-600 hover:bg-amber-700",
  PATCH: "bg-purple-600 hover:bg-purple-700",
  DELETE: "bg-red-600 hover:bg-red-700",
};

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
    <div className="h-full flex flex-col bg-background" onKeyDown={handleKeyDown}>
      {/* Method and URL Row */}
      <div className="p-6 border-b bg-card">
        <div className="flex gap-3 items-stretch">
          <Select
            value={props.method}
            onChange={(e) =>
              props.onMethodChange(
                e.target.value as "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
              )
            }
            className={`w-36 font-semibold ${methodColors[props.method]} text-white border-0 focus:ring-2 focus:ring-offset-2`}
          >
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="PATCH">PATCH</option>
            <option value="DELETE">DELETE</option>
          </Select>
          <div className="flex-1 flex flex-col gap-2">
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
              className="text-base h-11"
            />
            {fullUrl && fullUrl !== props.baseUrl + props.url && (
              <div className="flex items-center gap-2 text-xs text-muted-foreground px-1">
                <InfoIcon className="w-3 h-3" />
                <span>Full URL: <span className="font-mono">{fullUrl}</span></span>
              </div>
            )}
          </div>
        </div>
        <div className="mt-3 text-xs text-muted-foreground">
          Press <kbd className="px-1.5 py-0.5 rounded bg-muted border">Ctrl</kbd> + <kbd className="px-1.5 py-0.5 rounded bg-muted border">Enter</kbd> to execute request
        </div>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
        <TabsList className="w-full justify-start rounded-none border-b bg-transparent px-6 py-0 h-auto">
          <TabsTrigger value="query" className="relative data-[state=active]:shadow-none">
            <span>Query Params</span>
            {props.queryParams.length > 0 && (
              <Badge variant="secondary" className="ml-2 h-5 min-w-5 flex items-center justify-center px-1.5">
                {props.queryParams.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="headers" className="relative data-[state=active]:shadow-none">
            <span>Headers</span>
            {props.headers.length > 0 && (
              <Badge variant="secondary" className="ml-2 h-5 min-w-5 flex items-center justify-center px-1.5">
                {props.headers.length}
              </Badge>
            )}
          </TabsTrigger>
          {props.method !== "GET" && props.method !== "DELETE" && (
            <TabsTrigger value="body" className="data-[state=active]:shadow-none">Body</TabsTrigger>
          )}
          <TabsTrigger value="auth" className="data-[state=active]:shadow-none">Auth</TabsTrigger>
          <TabsTrigger value="client" className="data-[state=active]:shadow-none">HTTP Client</TabsTrigger>
          <TabsTrigger value="pagination" className="data-[state=active]:shadow-none">Pagination</TabsTrigger>
          <TabsTrigger value="parsing" className="data-[state=active]:shadow-none">Response Parsing</TabsTrigger>
        </TabsList>

        <div className="flex-1 overflow-auto">
          {/* Query Params Tab */}
          <TabsContent value="query" className="mt-0 p-6 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-semibold">Query Parameters</h3>
                <p className="text-sm text-muted-foreground mt-1">Add URL query parameters to your request</p>
              </div>
              <Button
                variant="default"
                size="sm"
                onClick={() =>
                  props.onQueryParamsChange([
                    ...props.queryParams,
                    { key: "", value: "" },
                  ])
                }
              >
                <PlusIcon className="h-4 w-4 mr-2" />
                Add Parameter
              </Button>
            </div>
            <Separator />
            {props.queryParams.length === 0 ? (
              <div className="text-center py-16">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                  <InfoIcon className="w-8 h-8 text-muted-foreground" />
                </div>
                <p className="text-base font-medium">No query parameters yet</p>
                <p className="text-sm text-muted-foreground mt-2 max-w-md mx-auto">
                  Query parameters are appended to the URL after a question mark (?)
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {props.queryParams.map((param, index) => (
                  <div key={index} className="flex gap-3 items-start p-3 rounded-lg border bg-card hover:bg-accent/5 transition-colors">
                    <div className="flex-1 grid grid-cols-2 gap-3">
                      <div>
                        <Label className="text-xs text-muted-foreground mb-1.5 block">Key</Label>
                        <Input
                          type="text"
                          value={param.key}
                          onChange={(e) => {
                            const newParams = [...props.queryParams];
                            newParams[index] = { key: e.target.value, value: param.value };
                            props.onQueryParamsChange(newParams);
                          }}
                          placeholder="parameter_name"
                          className="font-mono"
                        />
                      </div>
                      <div>
                        <Label className="text-xs text-muted-foreground mb-1.5 block">Value</Label>
                        <Input
                          type="text"
                          value={param.value}
                          onChange={(e) => {
                            const newParams = [...props.queryParams];
                            newParams[index] = { key: param.key, value: e.target.value };
                            props.onQueryParamsChange(newParams);
                          }}
                          placeholder="parameter_value"
                        />
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        props.onQueryParamsChange(
                          props.queryParams.filter((_, i) => i !== index)
                        );
                      }}
                      className="mt-7 hover:bg-destructive/10 hover:text-destructive"
                    >
                      <XIcon className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Headers Tab */}
          <TabsContent value="headers" className="mt-0 p-6 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-semibold">HTTP Headers</h3>
                <p className="text-sm text-muted-foreground mt-1">Configure request headers for authentication, content type, etc.</p>
              </div>
              <Button
                variant="default"
                size="sm"
                onClick={() =>
                  props.onHeadersChange([
                    ...props.headers,
                    { key: "", value: "" },
                  ])
                }
              >
                <PlusIcon className="h-4 w-4 mr-2" />
                Add Header
              </Button>
            </div>
            <Separator />
            {props.headers.length === 0 ? (
              <div className="text-center py-16">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                  <InfoIcon className="w-8 h-8 text-muted-foreground" />
                </div>
                <p className="text-base font-medium">No headers configured</p>
                <p className="text-sm text-muted-foreground mt-2 max-w-md mx-auto">
                  Headers provide additional information about the request or response
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {props.headers.map((header, index) => (
                  <div key={index} className="flex gap-3 items-start p-3 rounded-lg border bg-card hover:bg-accent/5 transition-colors">
                    <div className="flex-1 grid grid-cols-2 gap-3">
                      <div>
                        <Label className="text-xs text-muted-foreground mb-1.5 block">Header Name</Label>
                        <Input
                          type="text"
                          value={header.key}
                          onChange={(e) => {
                            const newHeaders = [...props.headers];
                            newHeaders[index] = { key: e.target.value, value: header.value };
                            props.onHeadersChange(newHeaders);
                          }}
                          placeholder="Content-Type"
                          className="font-mono"
                        />
                      </div>
                      <div>
                        <Label className="text-xs text-muted-foreground mb-1.5 block">Header Value</Label>
                        <Input
                          type="text"
                          value={header.value}
                          onChange={(e) => {
                            const newHeaders = [...props.headers];
                            newHeaders[index] = { key: header.key, value: e.target.value };
                            props.onHeadersChange(newHeaders);
                          }}
                          placeholder="application/json"
                        />
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        props.onHeadersChange(
                          props.headers.filter((_, i) => i !== index)
                        );
                      }}
                      className="mt-7 hover:bg-destructive/10 hover:text-destructive"
                    >
                      <XIcon className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Body Tab */}
          <TabsContent value="body" className="mt-0 p-6 space-y-4">
            <div>
              <h3 className="text-lg font-semibold">Request Body</h3>
              <p className="text-sm text-muted-foreground mt-1">Enter JSON, XML, or other payload data</p>
            </div>
            <Separator />
            <Textarea
              value={props.body}
              onChange={(e) => props.onBodyChange(e.target.value)}
              placeholder={'{\n  "key": "value",\n  "example": "data"\n}'}
              rows={18}
              className="font-mono text-sm resize-none"
            />
          </TabsContent>

          {/* Auth Tab */}
          <TabsContent value="auth" className="mt-0 p-6 space-y-4">
            <div>
              <h3 className="text-lg font-semibold">Authentication</h3>
              <p className="text-sm text-muted-foreground mt-1">Configure authentication credentials for your API</p>
            </div>
            <Separator />
            <div className="space-y-6">
              <div className="space-y-3">
                <Label htmlFor="auth-type" className="text-sm font-medium">Authentication Type</Label>
                <Select
                  id="auth-type"
                  value={props.authMethod}
                  onChange={(e) =>
                    props.onAuthMethodChange(
                      e.target.value as "none" | "basic" | "bearer"
                    )
                  }
                  className="w-full"
                >
                  <option value="none">No Authentication</option>
                  <option value="basic">Basic Auth</option>
                  <option value="bearer">****** Token</option>
                </Select>
              </div>

              {props.authMethod === "basic" && (
                <div className="space-y-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary">
                  <div className="space-y-3">
                    <Label htmlFor="auth-username">Username</Label>
                    <Input
                      id="auth-username"
                      type="text"
                      value={props.authUsername}
                      onChange={(e) => props.onAuthUsernameChange(e.target.value)}
                      placeholder="Enter username"
                    />
                  </div>
                  <div className="space-y-3">
                    <Label htmlFor="auth-password">Password</Label>
                    <Input
                      id="auth-password"
                      type="password"
                      value={props.authPassword}
                      onChange={(e) => props.onAuthPasswordChange(e.target.value)}
                      placeholder="Enter password"
                    />
                  </div>
                </div>
              )}

              {props.authMethod === "bearer" && (
                <div className="p-4 rounded-lg bg-muted/30 border-l-4 border-primary">
                  <div className="space-y-3">
                    <Label htmlFor="bearer-token">****** Token</Label>
                    <Input
                      id="bearer-token"
                      type="password"
                      value={props.bearerToken}
                      onChange={(e) => props.onBearerTokenChange(e.target.value)}
                      placeholder="Enter your access token"
                      className="font-mono"
                    />
                    <p className="text-xs text-muted-foreground">
                      Token will be sent in the Authorization header
                    </p>
                  </div>
                </div>
              )}
            </div>
          </TabsContent>

          {/* HTTP Client Tab */}
          <TabsContent value="client" className="mt-0 p-6 space-y-4">
            <div>
              <h3 className="text-lg font-semibold">HTTP Client Settings</h3>
              <p className="text-sm text-muted-foreground mt-1">Configure connection and network settings</p>
            </div>
            <Separator />
            <div className="space-y-6 max-w-2xl">
              <div className="space-y-3">
                <Label htmlFor="base-url">Base URL</Label>
                <Input
                  id="base-url"
                  type="text"
                  value={props.baseUrl}
                  onChange={(e) => props.onBaseUrlChange(e.target.value)}
                  placeholder="https://api.example.com"
                  className="font-mono"
                />
                <p className="text-xs text-muted-foreground">
                  The base URL will be prepended to all request paths
                </p>
              </div>
              <div className="space-y-3">
                <Label htmlFor="timeout">Request Timeout (seconds)</Label>
                <Input
                  id="timeout"
                  type="number"
                  value={props.timeout}
                  onChange={(e) => props.onTimeoutChange(e.target.value)}
                  placeholder="30"
                  min="1"
                  max="300"
                  className="w-48"
                />
                <p className="text-xs text-muted-foreground">
                  Maximum time to wait for a response
                </p>
              </div>
            </div>
          </TabsContent>

          {/* Pagination Tab */}
          <TabsContent value="pagination" className="mt-0 p-6 space-y-4">
            <div>
              <h3 className="text-lg font-semibold">Pagination Options</h3>
              <p className="text-sm text-muted-foreground mt-1">Configure automatic pagination for list endpoints</p>
            </div>
            <Separator />
            <div className="space-y-6 max-w-2xl">
              <div className="flex items-center gap-3 p-4 rounded-lg border">
                <input
                  type="checkbox"
                  id="pagination-enabled"
                  checked={props.paginationEnabled}
                  onChange={(e) =>
                    props.onPaginationEnabledChange(e.target.checked)
                  }
                  className="h-5 w-5 rounded border-input"
                />
                <div className="flex-1">
                  <Label htmlFor="pagination-enabled" className="text-base font-medium cursor-pointer">
                    Enable Pagination Support
                  </Label>
                  <p className="text-sm text-muted-foreground mt-0.5">
                    Automatically handle paginated API responses
                  </p>
                </div>
              </div>

              {props.paginationEnabled && (
                <div className="space-y-6 p-4 rounded-lg bg-muted/30 border-l-4 border-primary">
                  <div className="space-y-3">
                    <Label htmlFor="page-size">Items Per Page</Label>
                    <Input
                      id="page-size"
                      type="number"
                      value={props.pageSize}
                      onChange={(e) => props.onPageSizeChange(e.target.value)}
                      placeholder="10"
                      min="1"
                      max="1000"
                      className="w-48"
                    />
                  </div>
                  <div className="space-y-3">
                    <Label htmlFor="page-param">Page Parameter Name</Label>
                    <Input
                      id="page-param"
                      type="text"
                      value={props.pageParam}
                      onChange={(e) => props.onPageParamChange(e.target.value)}
                      placeholder="page"
                      className="font-mono"
                    />
                    <p className="text-xs text-muted-foreground">
                      The query parameter name used for pagination
                    </p>
                  </div>
                </div>
              )}
            </div>
          </TabsContent>

          {/* Response Parsing Tab */}
          <TabsContent value="parsing" className="mt-0 p-6 space-y-4">
            <div>
              <h3 className="text-lg font-semibold">Response Parsing</h3>
              <p className="text-sm text-muted-foreground mt-1">Configure how response data should be processed</p>
            </div>
            <Separator />
            <div className="space-y-6 max-w-2xl">
              <div className="flex items-center gap-3 p-4 rounded-lg border">
                <input
                  type="checkbox"
                  id="parse-response"
                  checked={props.parseResponse}
                  onChange={(e) => props.onParseResponseChange(e.target.checked)}
                  className="h-5 w-5 rounded border-input"
                />
                <div className="flex-1">
                  <Label htmlFor="parse-response" className="text-base font-medium cursor-pointer">
                    Parse JSON Response
                  </Label>
                  <p className="text-sm text-muted-foreground mt-0.5">
                    Automatically parse and format JSON responses
                  </p>
                </div>
              </div>

              {props.parseResponse && (
                <div className="p-4 rounded-lg bg-muted/30 border-l-4 border-primary">
                  <div className="space-y-3">
                    <Label htmlFor="json-path">JSON Path (optional)</Label>
                    <Input
                      id="json-path"
                      type="text"
                      value={props.jsonPath}
                      onChange={(e) => props.onJsonPathChange(e.target.value)}
                      placeholder="$.data.items"
                      className="font-mono"
                    />
                    <p className="text-xs text-muted-foreground">
                      Extract specific data from response using JSON path notation (e.g., $.data.users[*].name)
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
