from fastapi import FastAPI
from opentelemetry import trace
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor
from opentelemetry.instrumentation.httpx import HTTPXClientInstrumentor
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
import structlog

# 1. Configure the Tracer Provider globally
trace.set_tracer_provider(TracerProvider())
tracer_provider = trace.get_tracer_provider()

# 2. Configure OTLP Exporter pointing to the OpenTelemetry Collector DaemonSet
# The collector is running at otel-collector:4317 inside the cluster
otlp_exporter = OTLPSpanExporter(endpoint="http://snisid-otel-collector:4317", insecure=True)
span_processor = BatchSpanProcessor(otlp_exporter)
tracer_provider.add_span_processor(span_processor)

# 3. Setup FastAPI Application
app = FastAPI()

# 4. Auto-instrument FastAPI (creates spans for every HTTP request automatically)
FastAPIInstrumentor.instrument_app(app)

# 5. Auto-instrument database calls (captures SQL queries as spans)
# engine = create_engine(DB_URL)
# SQLAlchemyInstrumentor().instrument(engine=engine)

# 6. Auto-instrument outgoing HTTP requests to other microservices
HTTPXClientInstrumentor().instrument()

# 7. Configure Structlog to inject Trace IDs into application JSON logs
# This allows Grafana to seamlessly jump from a trace to the exact logs for that request
def add_trace_id_to_log(logger, log_method, event_dict):
    span = trace.get_current_span()
    if span.is_recording():
        ctx = span.get_span_context()
        event_dict["trace_id"] = format(ctx.trace_id, "032x")
        event_dict["span_id"] = format(ctx.span_id, "016x")
    return event_dict

structlog.configure(
    processors=[
        add_trace_id_to_log,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer()
    ]
)
