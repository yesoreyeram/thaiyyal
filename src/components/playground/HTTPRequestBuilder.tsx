"use client";
import React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
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
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
      e.preventDefault();
      props.onRun();
    }
  };

  return (
    <div className="h-full overflow-auto bg-white dark:bg-black" onKeyDown={handleKeyDown}>
      <div className="p-6 space-y-4">
        {/* Method */}
        <div className="space-y-1.5">
          <Label htmlFor="method" className="text-sm">Method</Label>
          <Select
            id="method"
            value={props.method}
            onChange={(e) =>
              props.onMethodChange(
                e.target.value as "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
              )
            }
            className="w-full bg-white dark:bg-black border-gray-300 dark:border-gray-700"
          >
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="PATCH">PATCH</option>
            <option value="DELETE">DELETE</option>
          </Select>
        </div>

        {/* URL */}
        <div className="space-y-1.5">
          <Label htmlFor="url" className="text-sm">URL</Label>
          <Input
            id="url"
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
            className="font-mono text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
          />
        </div>

        {/* Query Parameters */}
        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <Label className="text-sm">Query Parameters</Label>
            <Button
              variant="ghost"
              size="sm"
              onClick={() =>
                props.onQueryParamsChange([
                  ...props.queryParams,
                  { key: "", value: "" },
                ])
              }
              className="h-7 text-xs"
            >
              <PlusIcon className="h-3 w-3 mr-1" />
              Add
            </Button>
          </div>
          {props.queryParams.length > 0 && (
            <div className="space-y-2 pl-4 border-l-2 border-gray-200 dark:border-gray-800">
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
                    placeholder="key"
                    className="flex-1 h-8 text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
                  />
                  <Input
                    type="text"
                    value={param.value}
                    onChange={(e) => {
                      const newParams = [...props.queryParams];
                      newParams[index] = { key: param.key, value: e.target.value };
                      props.onQueryParamsChange(newParams);
                    }}
                    placeholder="value"
                    className="flex-1 h-8 text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => {
                      props.onQueryParamsChange(
                        props.queryParams.filter((_, i) => i !== index)
                      );
                    }}
                    className="h-8 w-8"
                  >
                    <XIcon className="h-3 w-3" />
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Headers */}
        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <Label className="text-sm">Headers</Label>
            <Button
              variant="ghost"
              size="sm"
              onClick={() =>
                props.onHeadersChange([
                  ...props.headers,
                  { key: "", value: "" },
                ])
              }
              className="h-7 text-xs"
            >
              <PlusIcon className="h-3 w-3 mr-1" />
              Add
            </Button>
          </div>
          {props.headers.length > 0 && (
            <div className="space-y-2 pl-4 border-l-2 border-gray-200 dark:border-gray-800">
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
                    placeholder="Header-Name"
                    className="flex-1 h-8 text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
                  />
                  <Input
                    type="text"
                    value={header.value}
                    onChange={(e) => {
                      const newHeaders = [...props.headers];
                      newHeaders[index] = { key: header.key, value: e.target.value };
                      props.onHeadersChange(newHeaders);
                    }}
                    placeholder="value"
                    className="flex-1 h-8 text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => {
                      props.onHeadersChange(
                        props.headers.filter((_, i) => i !== index)
                      );
                    }}
                    className="h-8 w-8"
                  >
                    <XIcon className="h-3 w-3" />
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Body (for POST/PUT/PATCH) */}
        {props.method !== "GET" && props.method !== "DELETE" && (
          <div className="space-y-1.5">
            <Label htmlFor="body" className="text-sm">Request Body</Label>
            <Textarea
              id="body"
              value={props.body}
              onChange={(e) => props.onBodyChange(e.target.value)}
              placeholder={'{\n  "key": "value"\n}'}
              rows={6}
              className="font-mono text-xs bg-white dark:bg-black border-gray-300 dark:border-gray-700 resize-none"
            />
          </div>
        )}

        {/* Authentication */}
        <div className="space-y-1.5">
          <Label htmlFor="auth-method" className="text-sm">Authentication</Label>
          <Select
            id="auth-method"
            value={props.authMethod}
            onChange={(e) =>
              props.onAuthMethodChange(
                e.target.value as "none" | "basic" | "bearer"
              )
            }
            className="w-full bg-white dark:bg-black border-gray-300 dark:border-gray-700"
          >
            <option value="none">No Authentication</option>
            <option value="basic">Basic Auth</option>
            <option value="bearer">Bearer Token</option>
          </Select>
        </div>

        {props.authMethod === "basic" && (
          <>
            <div className="space-y-1.5 pl-4">
              <Label htmlFor="auth-username" className="text-sm">Username</Label>
              <Input
                id="auth-username"
                type="text"
                value={props.authUsername}
                onChange={(e) => props.onAuthUsernameChange(e.target.value)}
                placeholder="username"
                className="bg-white dark:bg-black border-gray-300 dark:border-gray-700"
              />
            </div>
            <div className="space-y-1.5 pl-4">
              <Label htmlFor="auth-password" className="text-sm">Password</Label>
              <Input
                id="auth-password"
                type="password"
                value={props.authPassword}
                onChange={(e) => props.onAuthPasswordChange(e.target.value)}
                placeholder="password"
                className="bg-white dark:bg-black border-gray-300 dark:border-gray-700"
              />
            </div>
          </>
        )}

        {props.authMethod === "bearer" && (
          <div className="space-y-1.5 pl-4">
            <Label htmlFor="bearer-token" className="text-sm">Bearer Token</Label>
            <Input
              id="bearer-token"
              type="password"
              value={props.bearerToken}
              onChange={(e) => props.onBearerTokenChange(e.target.value)}
              placeholder="token"
              className="font-mono bg-white dark:bg-black border-gray-300 dark:border-gray-700"
            />
          </div>
        )}

        {/* Base URL */}
        <div className="space-y-1.5">
          <Label htmlFor="base-url" className="text-sm">Base URL</Label>
          <Input
            id="base-url"
            type="text"
            value={props.baseUrl}
            onChange={(e) => props.onBaseUrlChange(e.target.value)}
            placeholder="https://api.example.com"
            className="font-mono text-sm bg-white dark:bg-black border-gray-300 dark:border-gray-700"
          />
        </div>

        {/* Timeout */}
        <div className="space-y-1.5">
          <Label htmlFor="timeout" className="text-sm">Timeout (seconds)</Label>
          <Input
            id="timeout"
            type="number"
            value={props.timeout}
            onChange={(e) => props.onTimeoutChange(e.target.value)}
            placeholder="30"
            min="1"
            max="300"
            className="bg-white dark:bg-black border-gray-300 dark:border-gray-700"
          />
        </div>

        {/* Pagination */}
        <div className="space-y-1.5">
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
            <Label htmlFor="pagination-enabled" className="text-sm cursor-pointer">
              Enable Pagination
            </Label>
          </div>
        </div>

        {props.paginationEnabled && (
          <>
            <div className="space-y-1.5 pl-4">
              <Label htmlFor="page-size" className="text-sm">Page Size</Label>
              <Input
                id="page-size"
                type="number"
                value={props.pageSize}
                onChange={(e) => props.onPageSizeChange(e.target.value)}
                placeholder="10"
                min="1"
                max="1000"
                className="bg-white dark:bg-black border-gray-300 dark:border-gray-700"
              />
            </div>
            <div className="space-y-1.5 pl-4">
              <Label htmlFor="page-param" className="text-sm">Page Parameter Name</Label>
              <Input
                id="page-param"
                type="text"
                value={props.pageParam}
                onChange={(e) => props.onPageParamChange(e.target.value)}
                placeholder="page"
                className="font-mono bg-white dark:bg-black border-gray-300 dark:border-gray-700"
              />
            </div>
          </>
        )}

        {/* Response Parsing */}
        <div className="space-y-1.5">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="parse-response"
              checked={props.parseResponse}
              onChange={(e) => props.onParseResponseChange(e.target.checked)}
              className="h-4 w-4"
            />
            <Label htmlFor="parse-response" className="text-sm cursor-pointer">
              Parse JSON Response
            </Label>
          </div>
        </div>

        {props.parseResponse && (
          <div className="space-y-1.5 pl-4">
            <Label htmlFor="json-path" className="text-sm">JSON Path (optional)</Label>
            <Input
              id="json-path"
              type="text"
              value={props.jsonPath}
              onChange={(e) => props.onJsonPathChange(e.target.value)}
              placeholder="$.data.items"
              className="font-mono bg-white dark:bg-black border-gray-300 dark:border-gray-700"
            />
          </div>
        )}
      </div>
    </div>
  );
}
