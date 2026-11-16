package telemetry

import (
	"testing"

	"github.com/arifrhm/telemetry-tax-lab/internal/testutil"
	"github.com/valyala/fasthttp"
)

func TestBaselineHandler_HandleFastHTTP(t *testing.T) {
	h := NewBaselineHandler("node-tax-01")

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody(testutil.SamplePayload)

	h.HandleFastHTTP(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", ctx.Response.StatusCode(), ctx.Response.Body())
	}
}

func BenchmarkBaseline_Serve(b *testing.B) {
	h := NewBaselineHandler("node-tax-bench")

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
