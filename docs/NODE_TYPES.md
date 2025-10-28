
# Node Types Catalogue

This file lists suggested node types for a workflow builder focused on observability and analytics, with short descriptions and execution considerations.

## Triggers

- Manual / Button: Start by user click. Useful for manual debugging and ad-hoc runs.

- Schedule: Cron or interval-based scheduling for periodic jobs.

- Webhook / Event: Receive external HTTP events or messages.

- Message Queue (Kafka/SQS): Start when a message arrives.

## Ingest / Connector Nodes

- HTTP Request (GET/POST/PUT): Generic REST client with retries, timeouts, auth.

- Database Query (Postgres/ClickHouse/BigQuery): Run queries and stream results.

- Metrics Scraper (Prometheus): Pull metrics endpoint and produce samples.

- Log/Tracing Connector (Loki/Elastic): Query logs / traces.

## Transform / Mapper Nodes

- JSONPath / JQ / Mapper: Field extraction and mapping using declarative mapping.

- Script (JS/Python): Flexible transform with sandboxing; return structured payload.

- SQL Transform: Run SQL-like transformations on a batch payload.

## Analytics / Aggregation Nodes

- Timeseries Aggregator: Bucket and aggregate values over time windows.

- Statistical: Compute percentiles, moving averages, anomaly scores.

- ML/Model Score: Call a model endpoint or a lightweight inference runtime.

## Control Flow Nodes

- If / Switch: Conditional routing based on payload or metadata.

- Parallel / Map: Fan-out execution over array elements and optionally join results.

- Delay / Wait: Pause for a specified time or until a condition.

- Loop / Iterator: Controlled loops with max-iterations and escape conditions.

## Output / Sink Nodes

- HTTP Webhook / Post: Call external APIs (alerting, webhook destinations).

- Alerting (PagerDuty/Slack/Email): Send notifications.

- Timeseries Writer: Push to Prometheus remote-write, InfluxDB, ClickHouse.

- Storage (S3/Blob): Persist artifacts, charts, or CSVs.

## Stateful & Orchestration Nodes

- Checkpoint / Persist: Save state to durable store.

- Saga / Compensate: Manage distributed transactions and compensating actions.

## Utility Nodes

- Logger: Structured logs with levels, optionally persisted to log store.

- Metrics: Emit custom metrics about workflow state.

- No-op / Comment: For documentation or temporarily disabling branches.

## Execution notes

- Each node should declare its input/output ports and optionally output schema to enable safer wiring on the canvas.

- Nodes that perform side-effects must accept idempotency keys and provide retry/backoff configuration.

- Transform nodes should be sandboxed and have resource/time limits in enterprise deployments.
