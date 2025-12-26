"""
OpenTelemetry tracing configuration for Python services.
"""
import os
from contextvars import ContextVar
from typing import Optional

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
from opentelemetry.sdk.trace import Resource, TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import Tracer

# Context variable for trace_id (for compatibility with existing logger)
trace_id_var: ContextVar[Optional[str]] = ContextVar("trace_id", default=None)

_tracer_provider: Optional[TracerProvider] = None
_tracer: Optional[Tracer] = None


def init_tracing(
    service_name: str,
    service_version: str = "1.0.0",
    environment: str = "development",
    otlp_endpoint: Optional[str] = None,
    sampling_rate: Optional[float] = None,
):
    """
    Initialize OpenTelemetry tracing with OTLP HTTP exporter.
    
    Args:
        service_name: Name of the service
        service_version: Version of the service
        environment: Environment (development, production, etc.)
        otlp_endpoint: OTLP endpoint URL (defaults to http://jaeger:4318)
        sampling_rate: Sampling rate 0.0-1.0 (defaults to 1.0 in dev, 0.1 in prod)
    """
    global _tracer_provider, _tracer
    
    if _tracer_provider is not None:
        # Already initialized
        return
    
    # Determine OTLP endpoint
    if otlp_endpoint is None:
        otlp_endpoint = os.getenv("OTLP_ENDPOINT", "http://jaeger:4318")
    
    # Determine sampling rate
    if sampling_rate is None:
        env_lower = environment.lower()
        if env_lower in ("production", "prod"):
            sampling_rate = 0.1  # 10% in production
        else:
            sampling_rate = 1.0  # 100% in development
    
    # Create resource
    resource = Resource.create({
        "service.name": service_name,
        "service.version": service_version,
        "deployment.environment": environment,
    })
    
    # Create OTLP exporter
    otlp_exporter = OTLPSpanExporter(
        endpoint=otlp_endpoint,
        # Use insecure for development (use TLS in production)
    )
    
    # Create tracer provider with sampling
    _tracer_provider = TracerProvider(
        resource=resource,
    )
    
    # Add span processor
    _tracer_provider.add_span_processor(
        BatchSpanProcessor(otlp_exporter)
    )
    
    # Set global tracer provider
    trace.set_tracer_provider(_tracer_provider)
    
    # Create tracer
    _tracer = trace.get_tracer(service_name)
    
    # Auto-instrument FastAPI and requests
    try:
        FastAPIInstrumentor().instrument()
        RequestsInstrumentor().instrument()
    except Exception as e:
        # Log but don't fail if instrumentation fails
        print(f"Warning: Failed to auto-instrument: {e}")


def get_tracer() -> Tracer:
    """Get the global tracer."""
    global _tracer
    if _tracer is None:
        # Return no-op tracer if not initialized
        return trace.NoOpTracer()
    return _tracer


def get_trace_id() -> Optional[str]:
    """Get trace ID from current span or context variable."""
    # Try to get from current span
    span = trace.get_current_span()
    if span and span.get_span_context().is_valid:
        return format(span.get_span_context().trace_id, "032x")
    
    # Fallback to context variable (for compatibility)
    return trace_id_var.get()


def set_trace_id(trace_id: str):
    """Set trace_id in context variable (for compatibility with existing logger)."""
    trace_id_var.set(trace_id)


def shutdown():
    """Shutdown the tracer provider."""
    global _tracer_provider
    if _tracer_provider:
        _tracer_provider.shutdown()
        _tracer_provider = None
