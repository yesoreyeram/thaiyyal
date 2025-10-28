# Example Workflows & Pseudo-specs

This file contains example workflow descriptions and a small pseudo-spec for the timeseries analytics pipeline referenced in the product brief. Use these as starting templates for building UIs, node SDK tests, and execution engine prototypes.

Example 1 — Observability timeseries pipeline (manual trigger)

Description: Periodically fetch metrics from an HTTP endpoint, extract fields, aggregate into timeseries buckets, store in a timeseries store, and render a chart.

Flow:

- Trigger: Manual Button / Schedule (cron)

- HTTP Node: GET `https://metrics.example.com/app/123/metrics`

- Extract Node: JSONPath or mapping to pick fields: timestamp, metric_name, value, tags

- Transform Node: Aggregate into fixed window buckets (e.g., 1m) and compute average/percentiles

- Timeseries Writer: Push aggregated buckets to ClickHouse/InfluxDB/Prometheus remote-write

- Chart Node: Render timeseries in the UI (or export as PNG/CSV)

Pseudo-spec (JSON-style, simplified)

{
"workflow": {
"id": "wf-timeseries-001",
"nodes": [
{"id": "trigger-1", "type": "manual_trigger"},
{"id": "http-1", "type": "http_get", "config": {"url": "https://metrics.example.com/app/123/metrics", "auth": "{{cred.http}}"}},
{"id": "extract-1", "type": "json_mapper", "config": {"mappings": {"ts": "$.items[*].ts", "value": "$.items[*].val"}}},
{"id": "agg-1", "type": "timeseries_aggregator", "config": {"window": "1m", "aggregates": ["avg","p95"]}},
{"id": "write-1", "type": "timeseries_writer", "config": {"backend": "clickhouse", "table": "metrics_rollup"}},
{"id": "chart-1", "type": "chart_renderer", "config": {"style": "line", "series": ["avg"]}}
],
"edges": [
["trigger-1", "http-1"],
["http-1", "extract-1"],
["extract-1", "agg-1"],
["agg-1", "write-1"],
["agg-1", "chart-1"]
]
}
}

Example 2 — Alert enrichment and routing (webhook trigger)

Description: Receive alert webhook, enrich with asset metadata from CMDB, run an enrichment transform, route to different alerting sinks based on severity.

Flow: Webhook -> Enrich (DB) -> If Severity>=P1 -> PagerDuty -> else -> Slack

Notes on using examples

- Use these examples to build sample data fixtures for the Test-run / Dry-run feature.

- Provide sample payloads near each trigger node in the UI to help users debug mapping and transforms.

## Expanded examples — data stitching, transformations, and visualization

Below are several additional example pipelines that focus on stitching data from multiple APIs, common transformations (joins, filters, data type conversions), and visualizations. Each is written as a human-readable flow and a small pseudo-spec to use for prototyping and tests.

Example A — Stitch metrics across two services (join + aggregate + chart)

Description: Two services expose related metrics about user activity. Stitch them by user_id, normalize types, filter relevant users, compute per-user aggregates and render a top-N timeseries chart.

Flow:

- Trigger: Schedule (every 1m)

- HTTP Fetch A: GET `https://api.service-a.internal/v1/usage?window=1m` (returns items: user_id, ts, clicks)

- HTTP Fetch B: GET `https://api.service-b.internal/v1/engagement?window=1m` (returns items: userId, timestamp, impressions)

- Normalize Node: Map/convert field names and types (e.g., `userId` -> `user_id`, string timestamps -> ISO, numeric strings -> int)

- Join Node: Join A and B results on `user_id` (inner or left-join option)

- Filter Node: Keep only users with clicks+impressions > threshold

- Aggregate Node: Compute per-user rate metrics (e.g., ctr = clicks / impressions) and roll up into 1m series

- Top-N Selector: Select top 10 users by CTR

- Chart Node: Render timeseries chart (multi-series with one line per user) or export CSV for reporting

Pseudo-spec (short)

```json
{
  "workflow": {
    "id": "wf-stitch-001",
    "nodes": [
      {
        "id": "trigger-1",
        "type": "schedule",
        "config": { "cron": "*/1 * * * *" }
      },
      {
        "id": "fetch-a",
        "type": "http_get",
        "config": {
          "url": "https://api.service-a.internal/v1/usage?window=1m",
          "auth": "{{cred.a}}"
        }
      },
      {
        "id": "fetch-b",
        "type": "http_get",
        "config": {
          "url": "https://api.service-b.internal/v1/engagement?window=1m",
          "auth": "{{cred.b}}"
        }
      },
      {
        "id": "norm",
        "type": "json_mapper",
        "config": {
          "mappings": {
            "user_id": "$.userId || $.user_id",
            "ts": "$.timestamp || $.ts",
            "clicks": "$.clicks"
          },
          "conversions": { "clicks": "int", "ts": "date" }
        }
      },
      {
        "id": "join",
        "type": "join",
        "config": {
          "left": "fetch-a",
          "right": "fetch-b",
          "on": "user_id",
          "strategy": "inner"
        }
      },
      {
        "id": "filter",
        "type": "filter",
        "config": { "expr": "payload.clicks + payload.impressions > 10" }
      },
      {
        "id": "agg",
        "type": "timeseries_aggregator",
        "config": { "window": "1m", "aggregates": ["avg", "sum"] }
      },
      { "id": "topn", "type": "top_n", "config": { "n": 10, "metric": "ctr" } },
      {
        "id": "chart",
        "type": "chart_renderer",
        "config": { "style": "line", "series": "topn" }
      }
    ],
    "edges": [
      ["trigger-1", "fetch-a"],
      ["trigger-1", "fetch-b"],
      ["fetch-a", "norm"],
      ["fetch-b", "norm"],
      ["norm", "join"],
      ["join", "filter"],
      ["filter", "agg"],
      ["agg", "topn"],
      ["topn", "chart"]
    ]
  }
}
```

Example B — Correlate logs and metrics for an incident (join on trace_id + heatmap)

Description: An alert webhook arrives with a trace_id. Query logs and metrics APIs for that trace_id/time window, join records, compute heatmap of error frequency by service, and visualize.

Flow:

- Trigger: Webhook with payload {trace_id, ts}

- Query Logs Node: Query log store for `trace_id` and time window

- Query Metrics Node: Query metrics API for same time window and service set

- Enrich Node: Attach metadata from CMDB (service owner, environment) via a DB/API call

- Join Node: Join logs and metrics on `service` + `time_bucket`

- Aggregate Node: Count errors per service and bucket

- Heatmap Node: Render heatmap of error counts across services and time buckets

Pseudo-spec (short)

```json
{
  "workflow": {
    "id": "wf-correlate-001",
    "nodes": [
      { "id": "trigger-webhook", "type": "webhook" },
      {
        "id": "logs",
        "type": "log_query",
        "config": { "query": "trace_id:{{payload.trace_id}}", "window": "5m" }
      },
      {
        "id": "metrics",
        "type": "http_get",
        "config": {
          "url": "https://metrics.internal/query?trace={{payload.trace_id}}"
        }
      },
      {
        "id": "enrich",
        "type": "db_lookup",
        "config": { "table": "cmdb_services", "key": "service" }
      },
      {
        "id": "join",
        "type": "join",
        "config": { "on": ["service", "time_bucket"] }
      },
      {
        "id": "agg",
        "type": "aggregator",
        "config": {
          "groupBy": ["service", "time_bucket"],
          "metrics": { "errors": "count" }
        }
      },
      { "id": "heatmap", "type": "heatmap_renderer" }
    ],
    "edges": [
      ["trigger-webhook", "logs"],
      ["trigger-webhook", "metrics"],
      ["logs", "enrich"],
      ["metrics", "enrich"],
      ["enrich", "join"],
      ["join", "agg"],
      ["agg", "heatmap"]
    ]
  }
}
```

Example C — Cross-API time alignment and correlation (resample + scatter plot)

Description: Two APIs provide related metrics sampled at different intervals. Resample both to a common resolution, interpolate/match missing points, compute correlation coefficient, and render scatter/line charts.

Flow:

- Trigger: Manual or schedule

- Fetch M1: HTTP GET API A -> timeseries A (irregular sample)

- Fetch M2: HTTP GET API B -> timeseries B (regular sample)

- Resample Node: Resample both series to 1m buckets (interpolate/forward-fill)

- Join by timestamp: Merge aligned buckets

- Compute Node: Compute correlation and moving-window covariance

- Chart Node: Scatter plot (M1 vs M2) and time series overlay

Pseudo-spec (short)

```json
{
  "workflow": {
    "id": "wf-resample-001",
    "nodes": [
      { "id": "t", "type": "manual_trigger" },
      {
        "id": "m1",
        "type": "http_get",
        "config": { "url": "https://api-a/metrics" }
      },
      {
        "id": "m2",
        "type": "http_get",
        "config": { "url": "https://api-b/metrics" }
      },
      {
        "id": "resample",
        "type": "resample",
        "config": { "window": "1m", "method": "linear" }
      },
      { "id": "join", "type": "join", "config": { "on": "timestamp" } },
      { "id": "corr", "type": "compute", "config": { "expr": "pearson(x,y)" } },
      { "id": "scatter", "type": "scatter_renderer" }
    ],
    "edges": [
      ["t", "m1"],
      ["t", "m2"],
      ["m1", "resample"],
      ["m2", "resample"],
      ["resample", "join"],
      ["join", "corr"],
      ["corr", "scatter"]
    ]
  }
}
```

Example D — Table report with type conversions, joins and CSV export

Description: Produce a CSV report joining customer records from an API with usage stats from a DB, perform type conversions and computing percentiles, then export as CSV and render a preview table.

Flow:

- Trigger: Schedule (daily)

- Fetch Customers: HTTP GET `https://api.company/internal/customers`

- Fetch Usage: DB Query -> returns rows {customer_id, usage_value_str}

- Convert Types: Convert `usage_value_str` -> int, normalize currencies

- Join: Join customers and usage on `customer_id`

- Compute Stats: Compute percentiles (p50, p90) and flags for anomalies

- Table Renderer: Render preview table in UI

- CSV Export: Write CSV to S3 and attach to run artifacts

Pseudo-spec (short)

```json
{
  "workflow": {
    "id": "wf-report-001",
    "nodes": [
      { "id": "sched", "type": "schedule", "config": { "cron": "0 1 * * *" } },
      {
        "id": "cust",
        "type": "http_get",
        "config": {
          "url": "https://api.company/internal/customers",
          "auth": "{{cred.api}}"
        }
      },
      {
        "id": "usage",
        "type": "db_query",
        "config": {
          "sql": "SELECT customer_id, usage_value FROM usage_table WHERE date = '{{today}}'"
        }
      },
      {
        "id": "conv",
        "type": "transform",
        "config": { "conversions": { "usage_value": "int" } }
      },
      { "id": "join", "type": "join", "config": { "on": "customer_id" } },
      {
        "id": "stats",
        "type": "compute",
        "config": {
          "exprs": {
            "p50": "percentile(usage_value,50)",
            "p90": "percentile(usage_value,90)"
          }
        }
      },
      { "id": "table", "type": "table_renderer" },
      {
        "id": "export",
        "type": "s3_write",
        "config": { "path": "s3://reports/{{workflow.id}}/{{run.id}}.csv" }
      }
    ],
    "edges": [
      ["sched", "cust"],
      ["sched", "usage"],
      ["cust", "join"],
      ["usage", "join"],
      ["join", "conv"],
      ["conv", "stats"],
      ["stats", "table"],
      ["stats", "export"]
    ]
  }
}
```

Execution & testing tips

- Provide representative sample payloads for each trigger to allow quick dry-run of subgraphs.
- For join operations, include small sample datasets (2–5 rows) for each input to verify join keys and behavior (left/inner/outer).
- For resampling/interpolation, include edge cases with missing data and verify your chosen fill strategy (null, zero, forward-fill).
- For transforms that accept user code, recommend adding a linter and runtime guard (time limit, memory limit) and show example unit tests.
