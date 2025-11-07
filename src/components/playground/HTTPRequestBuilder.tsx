"use client";
import React, { useState } from "react";

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
}

interface CollapsibleSectionProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

function CollapsibleSection({
  title,
  children,
  defaultOpen = true,
}: CollapsibleSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  return (
    <div className="border border-gray-800 rounded-lg overflow-hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center justify-between bg-gray-900 hover:bg-gray-800 px-4 py-3 transition-colors"
      >
        <h3 className="text-sm font-semibold text-gray-200">{title}</h3>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={2}
          stroke="currentColor"
          className={`w-4 h-4 text-gray-400 transition-transform ${
            isOpen ? "rotate-180" : ""
          }`}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="m19.5 8.25-7.5 7.5-7.5-7.5"
          />
        </svg>
      </button>
      {isOpen && <div className="p-4 bg-gray-950">{children}</div>}
    </div>
  );
}

export function HTTPRequestBuilder(props: HTTPRequestBuilderProps) {
  const addHeader = () => {
    props.onHeadersChange([...props.headers, { key: "", value: "" }]);
  };

  const updateHeader = (index: number, key: string, value: string) => {
    const newHeaders = [...props.headers];
    newHeaders[index] = { key, value };
    props.onHeadersChange(newHeaders);
  };

  const removeHeader = (index: number) => {
    props.onHeadersChange(props.headers.filter((_, i) => i !== index));
  };

  const addQueryParam = () => {
    props.onQueryParamsChange([...props.queryParams, { key: "", value: "" }]);
  };

  const updateQueryParam = (index: number, key: string, value: string) => {
    const newParams = [...props.queryParams];
    newParams[index] = { key, value };
    props.onQueryParamsChange(newParams);
  };

  const removeQueryParam = (index: number) => {
    props.onQueryParamsChange(props.queryParams.filter((_, i) => i !== index));
  };

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-4">
      {/* HTTP Client Configuration */}
      <CollapsibleSection title="HTTP Client Configuration">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Base URL
            </label>
            <input
              type="text"
              value={props.baseUrl}
              onChange={(e) => props.onBaseUrlChange(e.target.value)}
              placeholder="https://api.example.com"
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Authentication
            </label>
            <select
              value={props.authMethod}
              onChange={(e) =>
                props.onAuthMethodChange(
                  e.target.value as "none" | "basic" | "bearer"
                )
              }
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="none">No Authentication</option>
              <option value="basic">Basic Auth</option>
              <option value="bearer">Bearer Token</option>
            </select>
          </div>

          {props.authMethod === "basic" && (
            <div className="space-y-3 pl-4 border-l-2 border-blue-500">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Username
                </label>
                <input
                  type="text"
                  value={props.authUsername}
                  onChange={(e) => props.onAuthUsernameChange(e.target.value)}
                  placeholder="username"
                  className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Password
                </label>
                <input
                  type="password"
                  value={props.authPassword}
                  onChange={(e) => props.onAuthPasswordChange(e.target.value)}
                  placeholder="password"
                  className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          )}

          {props.authMethod === "bearer" && (
            <div className="pl-4 border-l-2 border-blue-500">
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Bearer Token
              </label>
              <input
                type="password"
                value={props.bearerToken}
                onChange={(e) => props.onBearerTokenChange(e.target.value)}
                placeholder="your-token-here"
                className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Timeout (seconds)
            </label>
            <input
              type="number"
              value={props.timeout}
              onChange={(e) => props.onTimeoutChange(e.target.value)}
              placeholder="30"
              min="1"
              max="300"
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </CollapsibleSection>

      {/* HTTP Request */}
      <CollapsibleSection title="HTTP Request">
        <div className="space-y-4">
          <div className="flex gap-3">
            <div className="w-32">
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Method
              </label>
              <select
                value={props.method}
                onChange={(e) =>
                  props.onMethodChange(
                    e.target.value as
                      | "GET"
                      | "POST"
                      | "PUT"
                      | "PATCH"
                      | "DELETE"
                  )
                }
                className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
                <option value="PATCH">PATCH</option>
                <option value="DELETE">DELETE</option>
              </select>
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-300 mb-2">
                URL Path
              </label>
              <input
                type="text"
                value={props.url}
                onChange={(e) => props.onUrlChange(e.target.value)}
                placeholder="/api/users"
                className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* Query Parameters */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-300">
                Query Parameters
              </label>
              <button
                onClick={addQueryParam}
                className="text-xs text-blue-400 hover:text-blue-300"
              >
                + Add Parameter
              </button>
            </div>
            {props.queryParams.length === 0 ? (
              <div className="text-xs text-gray-500 py-2">
                No query parameters added
              </div>
            ) : (
              <div className="space-y-2">
                {props.queryParams.map((param, index) => (
                  <div key={index} className="flex gap-2">
                    <input
                      type="text"
                      value={param.key}
                      onChange={(e) =>
                        updateQueryParam(index, e.target.value, param.value)
                      }
                      placeholder="Key"
                      className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <input
                      type="text"
                      value={param.value}
                      onChange={(e) =>
                        updateQueryParam(index, param.key, e.target.value)
                      }
                      placeholder="Value"
                      className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button
                      onClick={() => removeQueryParam(index)}
                      className="px-3 py-2 text-red-400 hover:text-red-300 hover:bg-gray-800 rounded"
                      title="Remove parameter"
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Headers */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-300">
                Headers
              </label>
              <button
                onClick={addHeader}
                className="text-xs text-blue-400 hover:text-blue-300"
              >
                + Add Header
              </button>
            </div>
            {props.headers.length === 0 ? (
              <div className="text-xs text-gray-500 py-2">
                No headers added
              </div>
            ) : (
              <div className="space-y-2">
                {props.headers.map((header, index) => (
                  <div key={index} className="flex gap-2">
                    <input
                      type="text"
                      value={header.key}
                      onChange={(e) =>
                        updateHeader(index, e.target.value, header.value)
                      }
                      placeholder="Header name"
                      className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <input
                      type="text"
                      value={header.value}
                      onChange={(e) =>
                        updateHeader(index, header.key, e.target.value)
                      }
                      placeholder="Header value"
                      className="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <button
                      onClick={() => removeHeader(index)}
                      className="px-3 py-2 text-red-400 hover:text-red-300 hover:bg-gray-800 rounded"
                      title="Remove header"
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Body */}
          {(props.method === "POST" ||
            props.method === "PUT" ||
            props.method === "PATCH") && (
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Request Body (JSON)
              </label>
              <textarea
                value={props.body}
                onChange={(e) => props.onBodyChange(e.target.value)}
                placeholder='{"key": "value"}'
                rows={8}
                className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          )}
        </div>
      </CollapsibleSection>

      {/* Pagination Options */}
      <CollapsibleSection title="Pagination Options" defaultOpen={false}>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="pagination-enabled"
              checked={props.paginationEnabled}
              onChange={(e) => props.onPaginationEnabledChange(e.target.checked)}
              className="w-4 h-4 text-blue-600 bg-gray-900 border-gray-700 rounded focus:ring-blue-500 focus:ring-2"
            />
            <label
              htmlFor="pagination-enabled"
              className="text-sm font-medium text-gray-300"
            >
              Enable Pagination
            </label>
          </div>

          {props.paginationEnabled && (
            <div className="space-y-3 pl-4 border-l-2 border-blue-500">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Page Size
                </label>
                <input
                  type="number"
                  value={props.pageSize}
                  onChange={(e) => props.onPageSizeChange(e.target.value)}
                  placeholder="10"
                  min="1"
                  max="1000"
                  className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Page Parameter Name
                </label>
                <input
                  type="text"
                  value={props.pageParam}
                  onChange={(e) => props.onPageParamChange(e.target.value)}
                  placeholder="page"
                  className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          )}
        </div>
      </CollapsibleSection>

      {/* Response Parsing Options */}
      <CollapsibleSection title="Response Parsing Options" defaultOpen={false}>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="parse-response"
              checked={props.parseResponse}
              onChange={(e) => props.onParseResponseChange(e.target.checked)}
              className="w-4 h-4 text-blue-600 bg-gray-900 border-gray-700 rounded focus:ring-blue-500 focus:ring-2"
            />
            <label
              htmlFor="parse-response"
              className="text-sm font-medium text-gray-300"
            >
              Parse JSON Response
            </label>
          </div>

          {props.parseResponse && (
            <div className="pl-4 border-l-2 border-blue-500">
              <label className="block text-sm font-medium text-gray-300 mb-2">
                JSON Path (optional)
              </label>
              <input
                type="text"
                value={props.jsonPath}
                onChange={(e) => props.onJsonPathChange(e.target.value)}
                placeholder="$.data.items"
                className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <p className="mt-1 text-xs text-gray-500">
                Extract specific data from response using JSON path notation
              </p>
            </div>
          )}
        </div>
      </CollapsibleSection>
    </div>
  );
}
