# Instrumentation recipes

## OpenTelemetry setup (any language)
Environment-driven configuration works across SDKs:
```bash
OTEL_SERVICE_NAME=api
OTEL_RESOURCE_ATTRIBUTES=service.version=1.8.2,deployment.environment=prod
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
OTEL_TRACES_SAMPLER=parentbased_traceidratio
OTEL_TRACES_SAMPLER_ARG=0.1
OTEL_PROPAGATORS=tracecontext,baggage
```
| Stack | Auto-instrumentation |
|---|---|
| Java | `-javaagent:opentelemetry-javaagent.jar` |
| Node.js | `node --require @opentelemetry/auto-instrumentations-node/register app.js` |
| Python | `opentelemetry-bootstrap -a install` then `opentelemetry-instrument python app.py` |
| .NET | `OpenTelemetry.AutoInstrumentation` or `AddOpenTelemetry().WithTracing(...)` |
| Go | manual SDK + `otelhttp`, `otelgrpc`, `otelsql` wrappers |
| Kubernetes | OTel Operator `Instrumentation` CR + pod annotation |

- `parentbased_*` respects the caller's sampling decision so traces are not broken mid-flight.
- Baggage is propagated to every downstream; never put PII or secrets in it.

## Manual span pattern (pseudo-code)
```
span = tracer.start_span("checkout.place_order", attributes={"order.items": n})
try:
    ...
except Exception as e:
    span.record_exception(e); span.set_status(ERROR); raise
finally:
    span.end()
```
- Span names are low-cardinality (`GET /users/{id}`, not `/users/42`). IDs go in attributes, not names.

## Messaging propagation
- Producer: `propagator.inject(context, message.headers)`. Consumer: `ctx = propagator.extract(message.headers)`; start a consumer span with that parent (or a link for batches).
- Record queue time: enqueue timestamp header -> histogram `messaging_queue_time_seconds`.

## Log correlation
- Add `trace_id` and `span_id` from the active span to every log record (OTel log bridge / MDC in Java, `structlog` processor in Python, `pino` mixin in Node).
- Example line:
```json
{"timestamp":"2025-01-01T12:00:00.123Z","level":"error","service":"api","version":"1.8.2","environment":"prod","trace_id":"4bf92f3577b34da6a3ce929d0e0e4736","span_id":"00f067aa0ba902b7","request_id":"r-123","route":"/orders/{id}","status":502,"duration_ms":1840,"error.type":"UpstreamTimeout","message":"payment provider timeout"}
```

## Metrics to add per service
| Metric | Type | Labels |
|---|---|---|
| `http_server_request_duration_seconds` | histogram | `method`, `route`, `status_class` or `code` |
| `http_server_requests_total` | counter (or histogram `_count`) | same |
| `<dep>_client_request_duration_seconds` | histogram | `dependency`, `operation`, `outcome` |
| `db_pool_connections{state}` | gauge | `pool`, `state=idle|used` |
| `queue_messages_processed_total` | counter | `queue`, `outcome` |
| `queue_oldest_message_age_seconds` | gauge | `queue` |
| `build_info` | gauge = 1 | `version`, `commit` (for deploy correlation) |
- Prefer OTel semantic-convention names if the stack is OTel-native (`http.server.request.duration`); do not mix both schemes in one service.

## Collector pipeline (OTel Collector)
```yaml
processors:
  memory_limiter: { check_interval: 1s, limit_percentage: 80, spike_limit_percentage: 25 }
  batch: {}
  attributes/redact:
    actions:
      - { key: http.request.header.authorization, action: delete }
      - { key: user.email, action: hash }
  tail_sampling:
    decision_wait: 10s
    policies:
      - { name: errors, type: status_code, status_code: { status_codes: [ERROR] } }
      - { name: slow, type: latency, latency: { threshold_ms: 1000 } }
      - { name: rest, type: probabilistic, probabilistic: { sampling_percentage: 5 } }
service:
  pipelines:
    traces: { receivers: [otlp], processors: [memory_limiter, attributes/redact, tail_sampling, batch], exporters: [otlp] }
```
- `memory_limiter` first, `batch` last. Tail sampling needs all spans of a trace on the same collector (use a load-balancing exporter by trace ID).

## Health endpoints
- `/healthz` (liveness): process is alive, no dependency calls. `/ready` (readiness): can serve now (warm caches, critical config loaded). Do not expose internals publicly.

## Frontend / RUM
- Collect Core Web Vitals (LCP, INP, CLS) via `web-vitals` library to a RUM backend, tag with route template and release version. Propagate `traceparent` on API calls from the browser only to your own origins (CORS `Access-Control-Allow-Headers: traceparent`).

## Verification
- Local: run the collector with the `debug` exporter and check spans have parents across services.
- Confirm a log line can be found from a trace ID and vice versa in the backend.
- Check new series count after deploy: `count({__name__=~"<prefix>.*"})`.
