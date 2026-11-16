package telemetry

import (
	"io"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/rs/zerolog"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

type LoggingMode string

const (
	LogModeZerolog LoggingMode = "zerolog"
	LogModeSlog    LoggingMode = "slog"
	LogModeStdLog  LoggingMode = "stdlog"
)

type LoggingHandler struct {
	WorkerID      string
	workerIDBytes []byte
	Mode          LoggingMode
	zLogger       zerolog.Logger
	sLogger       *slog.Logger
	stdLogger     *log.Logger
}

func NewLoggingHandler(workerID string, mode LoggingMode) *LoggingHandler {
	// Send log bytes to io.Discard to purely measure CPU & Memory allocation tax
	zLog := zerolog.New(io.Discard).With().Timestamp().Str("worker", workerID).Logger()
	sLog := slog.New(slog.NewJSONHandler(io.Discard, nil))
	stdLog := log.New(io.Discard, "", 0)

	return &LoggingHandler{
		WorkerID:      workerID,
		workerIDBytes: []byte(workerID),
		Mode:          mode,
		zLogger:       zLog,
		sLogger:       sLog,
		stdLogger:     stdLog,
	}
}

// HandleFastHTTP represents Tier 2: The Structured Logging Tax.
// It formats and outputs a structured log event on every single request.
func (h *LoggingHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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

	// 1. Execute Logging based on Mode
	switch h.Mode {
	case LogModeZerolog:
		// Zero-allocation structured logging
		h.zLogger.Info().
			Bytes("event_id", eventIDBytes).
			Bytes("tenant_id", tenantIDBytes).
			Int64("timestamp", timestamp).
			Int("status", 200).
			Dur("elapsed", time.Since(start)).
			Msg("event_processed")
	case LogModeSlog:
		// Standard slog (reflection and interface heap allocation)
		h.sLogger.Info("event_processed",
			slog.String("event_id", string(eventIDBytes)),
			slog.String("tenant_id", string(tenantIDBytes)),
			slog.Int64("timestamp", timestamp),
			slog.Int("status", 200),
			slog.Duration("elapsed", time.Since(start)),
		)
	case LogModeStdLog:
		// Traditional string formatting
		h.stdLogger.Printf("event_processed event_id=%s tenant_id=%s timestamp=%d status=200 elapsed=%s",
			eventIDBytes, tenantIDBytes, timestamp, time.Since(start))
	}

	// Response formatting
	buf := bytebufferpool.Get()
	buf.Reset()
	defer bytebufferpool.Put(buf)

	var numBuf [32]byte
	buf.WriteString(`{"tier":"tier2_logging","mode":"`)
	buf.WriteString(string(h.Mode))
	buf.WriteString(`","status":"SUCCESS","event_id":"`)
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
