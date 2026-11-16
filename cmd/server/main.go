package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arifrhm/telemetry-tax-lab/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	workerID := os.Getenv("WORKER_ID")
	if workerID == "" {
		workerID = "node-tax-lab-01"
	}

	mode := os.Getenv("HANDLER_MODE")
	if mode == "" {
		mode = "baseline"
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Instantiate handlers
	baselineHandler := telemetry.NewBaselineHandler(workerID)
	promHandler := telemetry.NewPrometheusHandler(workerID)
	zerologHandler := telemetry.NewLoggingHandler(workerID, telemetry.LogModeZerolog)
	slogHandler := telemetry.NewLoggingHandler(workerID, telemetry.LogModeSlog)
	otel100Handler := telemetry.NewOTelTraceHandler(workerID, 1.0)
	defer otel100Handler.Close()
	otel1Handler := telemetry.NewOTelTraceHandler(workerID, 0.01)
	defer otel1Handler.Close()

	promAdaptor := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(promHandler.Registry, promhttp.HandlerOpts{}))

	router := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		switch path {
		case "/health":
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString(fmt.Sprintf(`{"status":"HEALTHY","mode":"%s","worker":"%s"}`, mode, workerID))
		case "/metrics":
			promAdaptor(ctx)
		case "/events":
			switch mode {
			case "prometheus":
				promHandler.HandleFastHTTP(ctx)
			case "zerolog":
				zerologHandler.HandleFastHTTP(ctx)
			case "slog":
				slogHandler.HandleFastHTTP(ctx)
			case "otel-trace-100":
				otel100Handler.HandleFastHTTP(ctx)
			case "otel-trace-1":
				otel1Handler.HandleFastHTTP(ctx)
			default:
				baselineHandler.HandleFastHTTP(ctx)
			}
		// Explicit path testing
		case "/events/baseline":
			baselineHandler.HandleFastHTTP(ctx)
		case "/events/prometheus":
			promHandler.HandleFastHTTP(ctx)
		case "/events/zerolog":
			zerologHandler.HandleFastHTTP(ctx)
		case "/events/slog":
			slogHandler.HandleFastHTTP(ctx)
		case "/events/otel-trace-100":
			otel100Handler.HandleFastHTTP(ctx)
		case "/events/otel-trace-1":
			otel1Handler.HandleFastHTTP(ctx)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Not Found")
		}
	}

	server := &fasthttp.Server{
		Handler:            router,
		Name:               "TelemetryTax-FastHTTP",
		Concurrency:        256 * 1024,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		MaxRequestBodySize: 10 * 1024 * 1024,
	}

	go func() {
		log.Printf("🚀 [Telemetry Tax Lab] Server listening on :%s (Mode: %s, Worker: %s)", port, mode, workerID)
		if err := server.ListenAndServe(":" + port); err != nil {
			log.Fatalf("FastHTTP ListenAndServe failed: %v", err)
		}
	}()

	<-stopChan
	log.Println("🛑 Shutting down server gracefully...")
	if err := server.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("✅ Server stopped gracefully")
}
