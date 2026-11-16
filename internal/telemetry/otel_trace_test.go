package telemetry

import (
	"testing"

	"github.com/arifrhm/telemetry-tax-lab/internal/testutil"
	"github.com/valyala/fasthttp"
)

func TestOTelTraceHandler(t *testing.T) {
	h := NewOTelTraceHandler("node-test", 1.0)
	defer h.Close()

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(testutil.SamplePayload)

	h.HandleFastHTTP(&ctx)
	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", ctx.Response.StatusCode(), ctx.Response.Body())
	}
}

func BenchmarkOTel_Tracing_100Percent(b *testing.B) {
	h := NewOTelTraceHandler("node-bench", 1.0)
	defer h.Close()

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

func BenchmarkOTel_Tracing_1Percent(b *testing.B) {
	h := NewOTelTraceHandler("node-bench", 0.01)
	defer h.Close()

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

func BenchmarkOTel_Tracing_ZeroPercent(b *testing.B) {
	h := NewOTelTraceHandler("node-bench", 0.0)
	defer h.Close()

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
