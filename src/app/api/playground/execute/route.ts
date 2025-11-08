import { NextRequest, NextResponse } from "next/server";

/**
 * HTTP Playground Execute API Route
 * 
 * Executes HTTP requests from the playground with proper error handling,
 * timeout management, and security controls.
 */

interface ExecuteRequestBody {
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  url: string;
  headers?: Record<string, string>;
  body?: string;
  timeout?: number;
}

interface ExecuteResponse {
  status: number;
  statusText: string;
  headers: Record<string, string>;
  data: unknown;
  timing: {
    startTime: number;
    endTime: number;
    duration: number;
  };
}

interface ErrorResponse {
  error: string;
  message: string;
  details?: unknown;
}

/**
 * POST /api/playground/execute
 * 
 * Executes an HTTP request with the provided configuration.
 * Returns the response data, headers, status, and timing information.
 */
export async function POST(request: NextRequest): Promise<NextResponse<ExecuteResponse | ErrorResponse>> {
  const startTime = Date.now();

  try {
    // Parse and validate request body
    const body: ExecuteRequestBody = await request.json();
    const { method, url, headers = {}, body: requestBody, timeout = 30000 } = body;

    // Validate required fields
    if (!method || !url) {
      return NextResponse.json(
        {
          error: "Validation Error",
          message: "Method and URL are required",
        },
        { status: 400 }
      );
    }

    // Validate URL format
    let targetUrl: URL;
    try {
      targetUrl = new URL(url);
    } catch {
      return NextResponse.json(
        {
          error: "Invalid URL",
          message: "The provided URL is not valid. Please provide a complete URL including protocol (http:// or https://)",
        },
        { status: 400 }
      );
    }

    // Validate timeout
    const timeoutMs = Math.min(Math.max(timeout, 1000), 60000); // Between 1s and 60s

    // Prepare fetch options
    const fetchOptions: RequestInit = {
      method,
      headers: {
        ...headers,
      },
      signal: AbortSignal.timeout(timeoutMs),
    };

    // Add body for methods that support it
    if (["POST", "PUT", "PATCH"].includes(method) && requestBody) {
      fetchOptions.body = requestBody;
      // Ensure Content-Type is set if not already present
      if (!fetchOptions.headers) {
        fetchOptions.headers = {};
      }
      const headersObj = fetchOptions.headers as Record<string, string>;
      if (!headersObj["Content-Type"] && !headersObj["content-type"]) {
        headersObj["Content-Type"] = "application/json";
      }
    }

    // Execute the HTTP request
    let response: Response;
    try {
      response = await fetch(targetUrl.toString(), fetchOptions);
    } catch (error) {
      const endTime = Date.now();
      
      // Handle specific fetch errors
      if (error instanceof Error) {
        if (error.name === "AbortError" || error.name === "TimeoutError") {
          return NextResponse.json(
            {
              error: "Request Timeout",
              message: `Request timed out after ${timeoutMs}ms`,
              details: { timeout: timeoutMs },
            },
            { status: 408 }
          );
        }
        
        if (error.message.includes("fetch failed")) {
          return NextResponse.json(
            {
              error: "Network Error",
              message: "Failed to connect to the target server. Please check the URL and network connectivity.",
              details: { originalError: error.message },
            },
            { status: 503 }
          );
        }
      }

      return NextResponse.json(
        {
          error: "Request Failed",
          message: error instanceof Error ? error.message : "An unknown error occurred",
          details: { duration: endTime - startTime },
        },
        { status: 500 }
      );
    }

    // Extract response data
    let responseData: unknown;
    const contentType = response.headers.get("content-type") || "";
    
    try {
      if (contentType.includes("application/json")) {
        responseData = await response.json();
      } else if (contentType.includes("text/")) {
        responseData = await response.text();
      } else {
        // For binary or other content types, get text representation
        responseData = await response.text();
      }
    } catch {
      // If parsing fails, return empty response
      responseData = null;
    }

    // Extract response headers
    const responseHeaders: Record<string, string> = {};
    response.headers.forEach((value, key) => {
      responseHeaders[key] = value;
    });

    const endTime = Date.now();

    // Return successful response
    return NextResponse.json<ExecuteResponse>(
      {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
        data: responseData,
        timing: {
          startTime,
          endTime,
          duration: endTime - startTime,
        },
      },
      { status: 200 }
    );
  } catch (error) {
    const endTime = Date.now();
    
    // Handle unexpected errors
    console.error("Unexpected error in playground execute:", error);
    
    return NextResponse.json<ErrorResponse>(
      {
        error: "Internal Server Error",
        message: error instanceof Error ? error.message : "An unexpected error occurred",
        details: { duration: endTime - startTime },
      },
      { status: 500 }
    );
  }
}
