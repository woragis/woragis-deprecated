/**
 * OpenTelemetry tracing configuration for Node.js service.
 */
import { NodeSDK } from '@opentelemetry/sdk-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { AmqplibInstrumentation } from '@opentelemetry/instrumentation-amqplib';
import { trace, context } from '@opentelemetry/api';

let sdk = null;
let traceId = null;

/**
 * Initialize OpenTelemetry tracing
 */
export function initTracing(serviceName, serviceVersion = '1.0.0', environment = 'development') {
  if (sdk) {
    // Already initialized
    return;
  }

  const otlpEndpoint = process.env.OTLP_ENDPOINT || 'http://jaeger:4318';
  
  // Determine sampling rate
  let samplingRate = 1.0; // 100% in development
  if (environment === 'production' || environment === 'prod') {
    samplingRate = 0.1; // 10% in production
  }

  const traceExporter = new OTLPTraceExporter({
    url: `${otlpEndpoint}/v1/traces`,
  });

  sdk = new NodeSDK({
    resource: new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: serviceVersion,
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: environment,
    }),
    traceExporter,
    instrumentations: [
      new HttpInstrumentation(),
      new AmqplibInstrumentation(),
    ],
    // Sampling configuration
    sampler: {
      shouldSample: () => {
        // Simple sampling based on rate
        return Math.random() < samplingRate;
      },
    },
  });

  sdk.start();
}

/**
 * Get current trace ID
 */
export function getTraceId() {
  // Try to get from current span
  const span = trace.getActiveSpan();
  if (span) {
    const spanContext = span.spanContext();
    if (spanContext.isValid) {
      return spanContext.traceId;
    }
  }
  
  // Fallback to stored trace ID (for compatibility)
  return traceId;
}

/**
 * Set trace ID (for compatibility with existing logger)
 */
export function setTraceId(id) {
  traceId = id;
}

/**
 * Shutdown tracing
 */
export function shutdown() {
  if (sdk) {
    sdk.shutdown();
    sdk = null;
  }
}
