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

    try {
      // Build complete URL with base URL, path, and query parameters
      const fullUrl = baseUrl + url;
      const urlObj = new URL(fullUrl);
      
      // Add query parameters
      queryParams
        .filter((param) => param.key.trim() !== "")
        .forEach((param) => {
          urlObj.searchParams.append(param.key, param.value);
        });

      // Build headers object
      const requestHeaders: Record<string, string> = {};
      headers
        .filter((h) => h.key.trim() !== "")
        .forEach((h) => {
          requestHeaders[h.key] = h.value;
        });

      // Add authentication headers
      if (authMethod === "basic" && authUsername && authPassword) {
        const credentials = btoa(`${authUsername}:${authPassword}`);
        requestHeaders["Authorization"] = `Basic ${credentials}`;
      } else if (authMethod === "bearer" && bearerToken) {
        requestHeaders["Authorization"] = `Bearer ${bearerToken}`;
      }

      // Prepare request body
      const requestBody = {
        method,
        url: urlObj.toString(),
        headers: requestHeaders,
        body: ["POST", "PUT", "PATCH"].includes(method) ? body : undefined,
        timeout: parseInt(timeoutSeconds) * 1000,
      };

      // Execute the HTTP request via API
      const response = await fetch("/api/playground/execute", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      const responseData = await response.json();

      if (!response.ok) {
        // Handle error response
        setError(
          responseData.message || "Request failed with status " + response.status
        );
        setIsExecuting(false);
        return;
      }

      // Set successful response
      setResult(responseData);
      setIsExecuting(false);
    } catch (err) {
      // Handle client-side errors
      setError(
        err instanceof Error
          ? err.message
          : "An unexpected error occurred while executing the request"
      );
      setIsExecuting(false);
    }
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
    link.download = `${requestTitle.replace(/[<>:"/\\|?*]/g, "_").toLowerCase()}_request.json`;
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
