"use client";
import React, { useState } from "react";
import { AppNavBar } from "../../components/AppNavBar";
import { PlaygroundNavBar } from "../../components/PlaygroundNavBar";
import { HTTPRequestBuilder } from "../../components/playground/HTTPRequestBuilder";
import { PlaygroundResultsPanel } from "../../components/playground/PlaygroundResultsPanel";

export default function PlaygroundPage() {
  const [isResultsPanelOpen, setIsResultsPanelOpen] = useState(false);
  const [isExecuting, setIsExecuting] = useState(false);
  const [result, setResult] = useState<unknown>(null);
  const [error, setError] = useState<string | null>(null);
  const [resultsPanelHeight, setResultsPanelHeight] = useState(250);

  // HTTP Request Configuration State
  const [baseUrl, setBaseUrl] = useState("");
  const [authMethod, setAuthMethod] = useState<"none" | "basic" | "bearer">(
    "none"
  );
  const [authUsername, setAuthUsername] = useState("");
  const [authPassword, setAuthPassword] = useState("");
  const [bearerToken, setBearerToken] = useState("");
  const [timeoutSeconds, setTimeoutSeconds] = useState("30");

  // HTTP Request State
  const [method, setMethod] = useState<
    "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
  >("GET");
  const [url, setUrl] = useState("");
  const [body, setBody] = useState("");
  const [headers, setHeaders] = useState<Array<{ key: string; value: string }>>(
    []
  );
  const [queryParams, setQueryParams] = useState<
    Array<{ key: string; value: string }>
  >([]);

  // Pagination State
  const [paginationEnabled, setPaginationEnabled] = useState(false);
  const [pageSize, setPageSize] = useState("10");
  const [pageParam, setPageParam] = useState("page");

  // Response Parsing State
  const [parseResponse, setParseResponse] = useState(true);
  const [jsonPath, setJsonPath] = useState("");

  const handleNewWorkflow = () => {
    window.location.href = "/workflow";
  };

  const handleOpenWorkflow = () => {
    window.location.href = "/workflow";
  };

  const handleRun = async () => {
    setIsResultsPanelOpen(true);
    setIsExecuting(true);
    setResult(null);
    setError(null);

    // Mock execution - simulate API call
    setTimeout(() => {
      // Mock successful response
      const mockResponse = {
        status: 200,
        statusText: "OK",
        headers: {
          "content-type": "application/json",
          "content-length": "123",
        },
        data: {
          message: "Mock response data",
          timestamp: new Date().toISOString(),
          request: {
            method,
            url: baseUrl + url,
            headers: headers.reduce(
              (acc, h) => ({ ...acc, [h.key]: h.value }),
              {}
            ),
          },
        },
      };

      setResult(mockResponse);
      setIsExecuting(false);
    }, 1500);
  };

  const handleCloseResultsPanel = () => {
    setIsResultsPanelOpen(false);
    setResult(null);
    setError(null);
  };

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      <AppNavBar
        onNewWorkflow={handleNewWorkflow}
        onOpenWorkflow={handleOpenWorkflow}
      />
      <PlaygroundNavBar onRun={handleRun} isRunning={isExecuting} />

      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="flex-1 overflow-auto">
          <HTTPRequestBuilder
            baseUrl={baseUrl}
            onBaseUrlChange={setBaseUrl}
            authMethod={authMethod}
            onAuthMethodChange={setAuthMethod}
            authUsername={authUsername}
            onAuthUsernameChange={setAuthUsername}
            authPassword={authPassword}
            onAuthPasswordChange={setAuthPassword}
            bearerToken={bearerToken}
            onBearerTokenChange={setBearerToken}
            timeout={timeoutSeconds}
            onTimeoutChange={setTimeoutSeconds}
            method={method}
            onMethodChange={setMethod}
            url={url}
            onUrlChange={setUrl}
            body={body}
            onBodyChange={setBody}
            headers={headers}
            onHeadersChange={setHeaders}
            queryParams={queryParams}
            onQueryParamsChange={setQueryParams}
            paginationEnabled={paginationEnabled}
            onPaginationEnabledChange={setPaginationEnabled}
            pageSize={pageSize}
            onPageSizeChange={setPageSize}
            pageParam={pageParam}
            onPageParamChange={setPageParam}
            parseResponse={parseResponse}
            onParseResponseChange={setParseResponse}
            jsonPath={jsonPath}
            onJsonPathChange={setJsonPath}
          />
        </div>

        <PlaygroundResultsPanel
          isOpen={isResultsPanelOpen}
          isLoading={isExecuting}
          result={result}
          error={error}
          onClose={handleCloseResultsPanel}
          height={resultsPanelHeight}
          onHeightChange={setResultsPanelHeight}
        />
      </div>
    </div>
  );
}
