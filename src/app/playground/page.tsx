"use client";
import React, { useState } from "react";
import {
  HomeIcon,
  FlaskConicalIcon,
  FileIcon,
  PlusIcon,
  PlayIcon,
  SaveIcon,
  DownloadIcon,
  UploadIcon,
  SettingsIcon,
} from "lucide-react";
import {
  SidebarProvider,
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenuItem,
  SidebarToggle,
  useSidebar,
} from "../../components/ui/sidebar";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { HTTPRequestBuilder } from "../../components/playground/HTTPRequestBuilder";
import { PlaygroundResultsPanel } from "../../components/playground/PlaygroundResultsPanel";

function PlaygroundContent() {
  const [requestTitle, setRequestTitle] = useState("New Request");
  const [isExecuting, setIsExecuting] = useState(false);
  const [result, setResult] = useState<unknown>(null);
  const [error, setError] = useState<string | null>(null);
  const [resultsPanelHeight, setResultsPanelHeight] = useState(250);
  const { isCollapsed } = useSidebar();

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

  const handleHome = () => {
    window.location.href = "/";
  };

  const handleNewWorkflow = () => {
    window.location.href = "/workflow";
  };

  const handleOpenWorkflow = () => {
    window.location.href = "/workflow";
  };

  const fileInputRef = React.useRef<HTMLInputElement>(null);

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

      // Prepare request body for Go backend API
      const requestBody = {
        method,
        url: urlObj.toString(),
        headers: requestHeaders,
        body: ["POST", "PUT", "PATCH"].includes(method) ? body : undefined,
        timeout: parseInt(timeoutSeconds), // timeout in seconds
      };

      // Execute the HTTP request via Go backend API
      const response = await fetch("/api/v1/playground/execute", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      const responseData = await response.json();

      // Check if there's an error in the response
      if (responseData.error) {
        setError(responseData.error);
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

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const json = JSON.parse(event.target?.result as string);
        handleImport(json);
      } catch {
        alert(
          "Failed to parse JSON file. Please ensure it is a valid JSON file."
        );
      }
    };
    reader.readAsText(file);

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <div className="h-screen flex bg-white dark:bg-black">
      {/* Collapsible Sidebar */}
      <Sidebar>
        <SidebarToggle />
        <SidebarHeader>
          <div className="flex items-center gap-2">
            <span className="text-xl">⚡</span>
            {!isCollapsed && (
              <span className="font-bold text-black dark:text-white">
                Thaiyyal
              </span>
            )}
          </div>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenuItem
                icon={<HomeIcon className="h-5 w-5" />}
                label="Home"
                onClick={handleHome}
              />
              <SidebarMenuItem
                icon={<FlaskConicalIcon className="h-5 w-5" />}
                label="Playground"
                onClick={() => {}}
                active={true}
              />
              <SidebarMenuItem
                icon={<FileIcon className="h-5 w-5" />}
                label="Workflows"
                onClick={handleOpenWorkflow}
              />
            </SidebarGroupContent>
          </SidebarGroup>

          <SidebarGroup>
            <SidebarGroupLabel>Actions</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenuItem
                icon={<PlusIcon className="h-5 w-5" />}
                label="New Workflow"
                onClick={handleNewWorkflow}
              />
              <SidebarMenuItem
                icon={<PlayIcon className="h-5 w-5" />}
                label="Run Request"
                onClick={handleRun}
              />
              <SidebarMenuItem
                icon={<SaveIcon className="h-5 w-5" />}
                label="Save Request"
                onClick={handleSave}
              />
              <SidebarMenuItem
                icon={<DownloadIcon className="h-5 w-5" />}
                label="Export"
                onClick={handleExport}
              />
              <SidebarMenuItem
                icon={<UploadIcon className="h-5 w-5" />}
                label="Import"
                onClick={handleImportClick}
              />
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <SidebarMenuItem
            icon={<SettingsIcon className="h-5 w-5" />}
            label="Settings"
            onClick={() => {}}
          />
        </SidebarFooter>
      </Sidebar>

      {/* Hidden file input for import */}
      <input
        ref={fileInputRef}
        type="file"
        accept=".json,application/json"
        onChange={handleFileChange}
        className="hidden"
      />

      {/* Main content area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top bar with request title and run button */}
        <div className="h-14 border-b border-gray-300 dark:border-gray-700 bg-white dark:bg-black flex items-center justify-between px-6">
          <div className="flex items-center gap-4">
            <Input
              type="text"
              value={requestTitle}
              onChange={(e) => setRequestTitle(e.target.value)}
              className="h-9 w-64 font-medium bg-white dark:bg-black"
              placeholder="Request name"
            />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              Press <kbd className="px-1.5 py-0.5 rounded bg-gray-100 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-xs">Ctrl+Enter</kbd> to run
            </span>
          </div>

          <Button
            onClick={handleRun}
            disabled={isExecuting}
            size="default"
            className="font-semibold bg-black dark:bg-white text-white dark:text-black hover:bg-gray-800 dark:hover:bg-gray-200"
          >
            {isExecuting ? (
              <>
                <span className="h-4 w-4 mr-2 border-2 border-current border-t-transparent rounded-full animate-spin" />
                Running...
              </>
            ) : (
              <>
                <PlayIcon className="h-4 w-4 mr-2 fill-current" />
                Run
              </>
            )}
          </Button>
        </div>

        <div className="flex-1 flex overflow-hidden">
          {/* Left half - Form */}
          <div className="w-1/2 overflow-auto border-r border-gray-300 dark:border-gray-700">
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

          {/* Right half - Results panel */}
          <div className="w-1/2 flex flex-col">
            <PlaygroundResultsPanel
              isLoading={isExecuting}
              result={result}
              error={error}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default function PlaygroundPage() {
  return (
    <SidebarProvider defaultCollapsed={false}>
      <PlaygroundContent />
    </SidebarProvider>
  );
}
