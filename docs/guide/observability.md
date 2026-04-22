# Observability

Jaeger and the OpenTelemetry Collector start with the default stack. No extra step is needed.

Tracing is enabled by default via `OTEL_ENABLED=true` in `.env.example`. To disable it for a run:

```bash
OTEL_ENABLED=false mise run run:api
```


## Viewing traces

Open [http://localhost:16686](http://localhost:16686). After making any HTTP request to `cashback-service-api`, select the service from the dropdown and click **Find Traces**.


## Pipeline

```mermaid
flowchart LR
    SVC["cashback-service-api"]
    OTEL["otel-collector"]
    JAEGER["Jaeger\n:4317 internal\nUI :16686"]

    SVC -->|"OTLP/gRPC :4317"| OTEL
    OTEL -->|"OTLP :4317"| JAEGER

    style SVC    fill:#1e3a5f,color:#bae6fd,stroke:#38bdf8
    style OTEL   fill:#78350f,color:#fed7aa,stroke:#f97316
    style JAEGER fill:#312e81,color:#c7d2fe,stroke:#6366f1
```

The collector receives spans from the service, batches them, and forwards to Jaeger over OTLP.


## Verifying

```bash
docker compose ps jaeger otel-collector
curl http://localhost:16686/api/services
```

After the API starts with tracing enabled, `cashback-service-api` appears in the service list.
