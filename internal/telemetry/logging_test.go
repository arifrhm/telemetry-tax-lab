package telemetry

import (
	"testing"

	"github.com/arifrhm/telemetry-tax-lab/internal/testutil"
	"github.com/valyala/fasthttp"
)

func TestLoggingHandler(t *testing.T) {
	modes := []LoggingMode{LogModeZerolog, LogModeSlog, LogModeStdLog}
	for _, m := range modes {
		h := NewLoggingHandler("node-test", m)
		var ctx fasthttp.RequestCtx
		ctx.Request.Header.SetMethod("POST")
		ctx.Request.Header.SetContentType("application/json")
		ctx.Request.SetBody(testutil.SamplePayload)

		h.HandleFastHTTP(&ctx)
		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("Mode %s failed with status %d", m, ctx.Response.StatusCode())
		}
	}
}

func BenchmarkLogging_Zerolog(b *testing.B) {
	h := NewLoggingHandler("node-bench", LogModeZerolog)
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(testutil.SamplePayload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.HandleFastHTTP(&ctx)
		ctx.Response.Reset()
	}
}

func BenchmarkLogging_Slog(b *testing.B) {
	h := NewLoggingHandler("node-bench", LogModeSlog)
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(testutil.SamplePayload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.HandleFastHTTP(&ctx)
		ctx.Response.Reset()
	}
}

func BenchmarkLogging_StdLog(b *testing.B) {
	h := NewLoggingHandler("node-bench", LogModeStdLog)
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(testutil.SamplePayload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.HandleFastHTTP(&ctx)
		ctx.Response.Reset()
	}
}
