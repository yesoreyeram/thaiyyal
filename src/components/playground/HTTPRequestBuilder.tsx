"use client";
import React, { useState } from "react";
import {
  Input,
  Select,
  Tabs,
  Checkbox,
  KeyValueEditor,
  CodeEditor,
} from "../../ui";

interface HTTPRequestBuilderProps {
  // HTTP Client Config
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

  // HTTP Request
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

  // Pagination
  paginationEnabled: boolean;
  onPaginationEnabledChange: (enabled: boolean) => void;
  pageSize: string;
  onPageSizeChange: (value: string) => void;
  pageParam: string;
  onPageParamChange: (value: string) => void;

  // Response Parsing
  parseResponse: boolean;
  onParseResponseChange: (enabled: boolean) => void;
  jsonPath: string;
  onJsonPathChange: (value: string) => void;

  // Run handler
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

  const tabs = [
    { id: "query", label: "Query Params", count: props.queryParams.length },
    { id: "headers", label: "Headers", count: props.headers.length },
    {
      id: "body",
      label: "Body",
      count: 0,
      hidden: props.method === "GET" || props.method === "DELETE",
    },
    { id: "auth", label: "Auth", count: 0 },
    { id: "client", label: "HTTP Client", count: 0 },
    { id: "pagination", label: "Pagination", count: 0 },
    { id: "parsing", label: "Response Parsing", count: 0 },
  ];

  return (
    <div className="h-full flex flex-col" onKeyDown={handleKeyDown}>
      {/* Method and URL Row */}
      <div className="p-4 border-b border-gray-800">
        <div className="flex gap-3 items-center">
          <Select
            value={props.method}
            onChange={(e) =>
              props.onMethodChange(
                e.target.value as "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
              )
            }
            className="!w-auto px-4 font-medium"
            options={[
              { value: "GET", label: "GET" },
              { value: "POST", label: "POST" },
              { value: "PUT", label: "PUT" },
              { value: "PATCH", label: "PATCH" },
              { value: "DELETE", label: "DELETE" },
            ]}
          />
          <div className="flex-1 flex flex-col gap-1">
            <Input
              type="text"
              value={props.baseUrl + props.url}
              onChange={(e) => {
                const value = e.target.value;
                // Try to split into base URL and path
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
                  // If URL parsing fails, just set as URL
                  props.onUrlChange(value);
                }
              }}
              placeholder="https://api.example.com/endpoint"
            />
            {fullUrl && fullUrl !== props.baseUrl + props.url && (
              <div className="text-xs text-gray-500 px-1">
                Full URL: {fullUrl}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} className="px-4 pt-2" />

      {/* Tab Content */}
      <div className="flex-1 overflow-auto p-4 bg-gray-950">
        {/* Query Params Tab */}
        {activeTab === "query" && (
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-300">
              Query Parameters
            </h3>
            <KeyValueEditor
              items={props.queryParams}
              onChange={props.onQueryParamsChange}
              keyPlaceholder="Parameter name"
              valuePlaceholder="Parameter value"
              emptyMessage='No query parameters. Click "+ Add" to add one.'
            />
          </div>
        )}

        {/* Headers Tab */}
        {activeTab === "headers" && (
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-300">Headers</h3>
            <KeyValueEditor
              items={props.headers}
              onChange={props.onHeadersChange}
              keyPlaceholder="Header name"
              valuePlaceholder="Header value"
              emptyMessage='No headers. Click "+ Add" to add one.'
            />
          </div>
        )}

        {/* Body Tab */}
        {activeTab === "body" && (
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-300">Request Body</h3>
            <CodeEditor
              value={props.body}
              onChange={(e) => props.onBodyChange(e.target.value)}
              placeholder='{"key": "value"}'
              rows={15}
              language="json"
            />
          </div>
        )}

        {/* Auth Tab */}
        {activeTab === "auth" && (
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-300">Authentication</h3>
            <Select
              label="Auth Type"
              value={props.authMethod}
              onChange={(e) =>
                props.onAuthMethodChange(
                  e.target.value as "none" | "basic" | "bearer"
                )
              }
              options={[
                { value: "none", label: "No Authentication" },
                { value: "basic", label: "Basic Auth" },
                { value: "bearer", label: "Bearer Token" },
              ]}
            />

            {props.authMethod === "basic" && (
              <div className="space-y-3 pl-4 border-l-2 border-blue-500">
                <Input
                  label="Username"
                  type="text"
                  value={props.authUsername}
                  onChange={(e) => props.onAuthUsernameChange(e.target.value)}
                  placeholder="username"
                />
                <Input
                  label="Password"
                  type="password"
                  value={props.authPassword}
                  onChange={(e) => props.onAuthPasswordChange(e.target.value)}
                  placeholder="password"
                />
              </div>
            )}

            {props.authMethod === "bearer" && (
              <div className="pl-4 border-l-2 border-blue-500">
                <Input
                  label="Bearer Token"
                  type="password"
                  value={props.bearerToken}
                  onChange={(e) => props.onBearerTokenChange(e.target.value)}
                  placeholder="your-token-here"
                />
              </div>
            )}
          </div>
        )}

        {/* HTTP Client Tab */}
        {activeTab === "client" && (
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-300">
              HTTP Client Settings
            </h3>
            <Input
              label="Base URL"
              type="text"
              value={props.baseUrl}
              onChange={(e) => props.onBaseUrlChange(e.target.value)}
              placeholder="https://api.example.com"
            />
            <Input
              label="Timeout (seconds)"
              type="number"
              value={props.timeout}
              onChange={(e) => props.onTimeoutChange(e.target.value)}
              placeholder="30"
              min="1"
              max="300"
            />
          </div>
        )}

        {/* Pagination Tab */}
        {activeTab === "pagination" && (
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-300">
              Pagination Options
            </h3>
            <Checkbox
              id="pagination-enabled"
              checked={props.paginationEnabled}
              onChange={(e) =>
                props.onPaginationEnabledChange(e.target.checked)
              }
              label="Enable Pagination"
            />

            {props.paginationEnabled && (
              <div className="space-y-3 pl-4 border-l-2 border-blue-500">
                <Input
                  label="Page Size"
                  type="number"
                  value={props.pageSize}
                  onChange={(e) => props.onPageSizeChange(e.target.value)}
                  placeholder="10"
                  min="1"
                  max="1000"
                />
                <Input
                  label="Page Parameter Name"
                  type="text"
                  value={props.pageParam}
                  onChange={(e) => props.onPageParamChange(e.target.value)}
                  placeholder="page"
                />
              </div>
            )}
          </div>
        )}

        {/* Response Parsing Tab */}
        {activeTab === "parsing" && (
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-300">
              Response Parsing
            </h3>
            <Checkbox
              id="parse-response"
              checked={props.parseResponse}
              onChange={(e) => props.onParseResponseChange(e.target.checked)}
              label="Parse JSON Response"
            />

            {props.parseResponse && (
              <div className="pl-4 border-l-2 border-blue-500">
                <Input
                  label="JSON Path (optional)"
                  type="text"
                  value={props.jsonPath}
                  onChange={(e) => props.onJsonPathChange(e.target.value)}
                  placeholder="$.data.items"
                  className="font-mono"
                  helperText="Extract specific data from response using JSON path notation"
                />
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
