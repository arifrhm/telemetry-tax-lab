package telemetry

import (
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

type PrometheusHandler struct {
	WorkerID      string
	workerIDBytes []byte
	Registry      *prometheus.Registry
	reqCounter    *prometheus.CounterVec
	reqDuration   *prometheus.HistogramVec
}

func NewPrometheusHandler(workerID string) *PrometheusHandler {
	reg := prometheus.NewRegistry()

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "telemetry_tax",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests processed.",
		},
		[]string{"route", "status"},
	)

	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "telemetry_tax",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency distribution in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"route"},
	)

	reg.MustRegister(counter, duration)

	return &PrometheusHandler{
		WorkerID:      workerID,
		workerIDBytes: []byte(workerID),
		Registry:      reg,
		reqCounter:    counter,
		reqDuration:   duration,
	}
}

// HandleFastHTTP represents Tier 1: Standard Prometheus Metrics Instrumentation.
// It tracks request volume and latency distribution across concurrent requests.
func (h *PrometheusHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

	// 1. Prometheus Telemetry Recording (The Metrics Tax)
	elapsed := time.Since(start).Seconds()
	h.reqCounter.WithLabelValues("/events", "200").Inc()
	h.reqDuration.WithLabelValues("/events").Observe(elapsed)

	// Response formatting
	buf := bytebufferpool.Get()
	buf.Reset()
	defer bytebufferpool.Put(buf)

	var numBuf [32]byte
	buf.WriteString(`{"tier":"tier1_prometheus","status":"SUCCESS","event_id":"`)
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
