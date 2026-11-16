package telemetry

import (
	"context"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type OTelTraceHandler struct {
	WorkerID      string
	workerIDBytes []byte
	tp            *sdktrace.TracerProvider
	tracer        trace.Tracer
	exporter      *tracetest.InMemoryExporter
	SamplingRate  float64
}

// NewOTelTraceHandler creates Tier 3 (100% Tracing) or Tier 4 (Sampled Tracing).
// samplingRatio = 1.0 means 100% AlwaysOn.
// samplingRatio = 0.01 means 1% Adaptive Ratio Sampling.
func NewOTelTraceHandler(workerID string, samplingRatio float64) *OTelTraceHandler {
	exporter := tracetest.NewInMemoryExporter()

	var sampler sdktrace.Sampler
	if samplingRatio >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if samplingRatio <= 0.0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingRatio))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithSyncer(exporter),
	)

	tracer := tp.Tracer("telemetry-tax-ingestor")

	return &OTelTraceHandler{
		WorkerID:      workerID,
		workerIDBytes: []byte(workerID),
		tp:            tp,
		tracer:        tracer,
		exporter:      exporter,
		SamplingRate:  samplingRatio,
	}
}

func (h *OTelTraceHandler) Close() error {
	return h.tp.Shutdown(context.Background())
}

// HandleFastHTTP represents Tier 3 / Tier 4: OpenTelemetry Distributed Tracing.
func (h *OTelTraceHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBodyString("Method not allowed")
		return
	}

	start := time.Now()
	bodyBytes := ctx.PostBody()

	// In-place byte extraction
	eventIDBytes, _, _, err := jsonparser.Get(bodyBytes, "event_id")
	if err != nil || len(eventIDBytes) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
		ctx.SetBodyString("event_id is required")
		return
	}

	tenantIDBytes, _, _, err := jsonparser.Get(bodyBytes, "tenant_id")
	if err != nil || len(tenantIDBytes) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
		ctx.SetBodyString("tenant_id is required")
		return
	}

	timestamp, err := jsonparser.GetInt(bodyBytes, "timestamp")
	if err != nil || timestamp <= 0 {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
		ctx.SetBodyString("timestamp must be greater than zero")
		return
	}

	// 1. OpenTelemetry Tracing Span Lifecycle (The Tracing Tax)
	// Tracing creates Go context, span objects, and attribute key-value allocations
	reqCtx := context.Background()
	_, span := h.tracer.Start(reqCtx, "events.ingest",
		trace.WithAttributes(
			attribute.String("event_id", string(eventIDBytes)),
			attribute.String("tenant_id", string(tenantIDBytes)),
			attribute.Int64("timestamp", timestamp),
			attribute.Int("http.status_code", 200),
			attribute.String("worker.node", h.WorkerID),
		),
	)
	span.End()

	// Response formatting
	buf := bytebufferpool.Get()
	buf.Reset()
	defer bytebufferpool.Put(buf)

	var numBuf [32]byte
	buf.WriteString(`{"tier":"otel_tracing","sampling_rate":`)
	buf.Write(strconv.AppendFloat(numBuf[:0], h.SamplingRate, 'f', 2, 64))
	buf.WriteString(`,"status":"SUCCESS","event_id":"`)
	buf.Write(eventIDBytes)
	buf.WriteString(`","tenant_id":"`)
	buf.Write(tenantIDBytes)
	buf.WriteString(`","ingest_ns":`)
	buf.Write(strconv.AppendInt(numBuf[:0], start.UnixNano(), 10))
	buf.WriteString(`,"worker_node":"`)
	buf.Write(h.workerIDBytes)
	buf.WriteString("\"}\n")

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(buf.Bytes())
}
