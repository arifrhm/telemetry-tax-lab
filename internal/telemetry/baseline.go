package telemetry

import (
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

type BaselineHandler struct {
	WorkerID      string
	workerIDBytes []byte
}

func NewBaselineHandler(workerID string) *BaselineHandler {
	return &BaselineHandler{
		WorkerID:      workerID,
		workerIDBytes: []byte(workerID),
	}
}

// HandleFastHTTP represents Tier 0: Pure Baseline with Zero Telemetry.
// It performs schema extraction and validation without touching metrics, logging, or tracing.
func (h *BaselineHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBodyString("Method not allowed")
		return
	}

	start := time.Now().UnixNano()
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

	// Zero-allocation response formatting via bytebufferpool
	buf := bytebufferpool.Get()
	buf.Reset()
	defer bytebufferpool.Put(buf)

	var numBuf [32]byte
	buf.WriteString(`{"tier":"tier0_baseline","status":"SUCCESS","event_id":"`)
	buf.Write(eventIDBytes)
	buf.WriteString(`","tenant_id":"`)
	buf.Write(tenantIDBytes)
	buf.WriteString(`","ingest_ns":`)
	buf.Write(strconv.AppendInt(numBuf[:0], start, 10))
	buf.WriteString(`,"worker_node":"`)
	buf.Write(h.workerIDBytes)
	buf.WriteString("\"}\n")

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(buf.Bytes())
}
