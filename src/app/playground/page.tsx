"use client";
import React, { useState } from "react";
import { AppNavBar } from "../../components/AppNavBar";
import { PlaygroundNavBar } from "../../components/PlaygroundNavBar";
import { HTTPRequestBuilder } from "../../components/playground/HTTPRequestBuilder";
import { PlaygroundResultsPanel } from "../../components/playground/PlaygroundResultsPanel";

export default function PlaygroundPage() {
  const [requestTitle, setRequestTitle] = useState("New Request");
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
          "x-request-id": "mock-id-" + Date.now(),
        },
        data: {
          message: "Mock response data",
          timestamp: new Date().toISOString(),
          request: {
            method,
            url: baseUrl + url,
            headers: headers
              .filter((h) => h.key.trim() !== "")
              .reduce((acc, h) => ({ ...acc, [h.key]: h.value }), {}),
          },
        },
      };

      setResult(mockResponse);
      setIsExecuting(false);
    }, 1500);
  };

  const handleSave = () => {
    // TODO: Implement save functionality
    console.log("Save request");
  };

  const handleExport = () => {
    const exportData = {
      title: requestTitle,
      baseUrl,
      authMethod,
      authUsername,
      authPassword,
      bearerToken,
      timeout: timeoutSeconds,
      method,
      url,
      body,
      headers,
      queryParams,
      paginationEnabled,
      pageSize,
      pageParam,
      parseResponse,
      jsonPath,
    };

    const jsonString = JSON.stringify(exportData, null, 2);
    const blob = new Blob([jsonString], { type: "application/json" });
    const urlObj = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = urlObj;
    link.download = `${requestTitle.replace(/[^a-z0-9]/gi, "_").toLowerCase()}_request.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(urlObj);
  };

  const handleImport = (data: unknown) => {
    try {
      const importedData = data as {
        title?: string;
        baseUrl?: string;
        authMethod?: "none" | "basic" | "bearer";
        authUsername?: string;
        authPassword?: string;
        bearerToken?: string;
        timeout?: string;
        method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
        url?: string;
        body?: string;
        headers?: Array<{ key: string; value: string }>;
        queryParams?: Array<{ key: string; value: string }>;
        paginationEnabled?: boolean;
        pageSize?: string;
        pageParam?: string;
        parseResponse?: boolean;
        jsonPath?: string;
      };

      if (importedData.title) setRequestTitle(importedData.title);
      if (importedData.baseUrl) setBaseUrl(importedData.baseUrl);
      if (importedData.authMethod) setAuthMethod(importedData.authMethod);
      if (importedData.authUsername) setAuthUsername(importedData.authUsername);
      if (importedData.authPassword) setAuthPassword(importedData.authPassword);
      if (importedData.bearerToken) setBearerToken(importedData.bearerToken);
      if (importedData.timeout) setTimeoutSeconds(importedData.timeout);
      if (importedData.method) setMethod(importedData.method);
      if (importedData.url) setUrl(importedData.url);
      if (importedData.body) setBody(importedData.body);
      if (importedData.headers) setHeaders(importedData.headers);
      if (importedData.queryParams) setQueryParams(importedData.queryParams);
      if (importedData.paginationEnabled !== undefined)
        setPaginationEnabled(importedData.paginationEnabled);
      if (importedData.pageSize) setPageSize(importedData.pageSize);
      if (importedData.pageParam) setPageParam(importedData.pageParam);
      if (importedData.parseResponse !== undefined)
        setParseResponse(importedData.parseResponse);
      if (importedData.jsonPath) setJsonPath(importedData.jsonPath);
    } catch {
      alert("Failed to import request configuration");
    }
  };

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      <AppNavBar
        onNewWorkflow={handleNewWorkflow}
        onOpenWorkflow={handleOpenWorkflow}
      />
      <PlaygroundNavBar
        requestTitle={requestTitle}
        onTitleChange={setRequestTitle}
        onRun={handleRun}
        isRunning={isExecuting}
        onSave={handleSave}
        onExport={handleExport}
        onImport={handleImport}
      />

      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="flex-1 overflow-hidden">
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
            onRun={handleRun}
          />
        </div>

        {/* Results panel always visible */}
        <PlaygroundResultsPanel
          isLoading={isExecuting}
          result={result}
          error={error}
          height={resultsPanelHeight}
          onHeightChange={setResultsPanelHeight}
        />
      </div>
    </div>
  );
}
