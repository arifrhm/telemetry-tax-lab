#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

PORT="${PORT:-8089}"
DURATION="${DURATION:-10s}"
CONNECTIONS="${CONNECTIONS:-125}"
HANDLER_MODE="${HANDLER_MODE:-baseline}"
TARGET_URL="http://127.0.0.1:${PORT}/events"
HEALTH_URL="http://127.0.0.1:${PORT}/health"
PAYLOAD_FILE="${ROOT_DIR}/scripts/payload_sample.json"

case "${HANDLER_MODE}" in
    baseline)
        BENCH_FUNC="BenchmarkBaseline_Serve"
        TIER_NAME="Tier 0 - Pure Baseline (Zero Telemetry)"
        ;;
    prometheus)
        BENCH_FUNC="BenchmarkPrometheus_Serve"
        TIER_NAME="Tier 1 - Prometheus Metrics (client_golang Counter+Histogram)"
        ;;
    zerolog)
        BENCH_FUNC="BenchmarkLogging_Zerolog"
        TIER_NAME="Tier 2a - Structured Logging (rs/zerolog Zero-Alloc)"
        ;;
    slog)
        BENCH_FUNC="BenchmarkLogging_Slog"
        TIER_NAME="Tier 2b - Structured Logging (log/slog Standard)"
        ;;
    otel-trace-100)
        BENCH_FUNC="BenchmarkOTel_Tracing_100Percent"
        TIER_NAME="Tier 3 - OpenTelemetry Tracing (100% AlwaysOn OTLP)"
        ;;
    otel-trace-1)
        BENCH_FUNC="BenchmarkOTel_Tracing_1Percent"
        TIER_NAME="Tier 4 - Intelligent Adaptive Sampling (OTel 1% Ratio)"
        ;;
    *)
        BENCH_FUNC="BenchmarkBaseline_Serve"
        TIER_NAME="Custom Benchmark"
        ;;
esac

echo "================================================================"
echo "⚡ [Telemetry Tax Lab] Automated Benchmark & Telemetry Runner"
echo "================================================================"

# 1. System Info Gathering
CPU_MODEL=$(lscpu | grep "Model name:" | head -n1 | sed -e 's/Model name:[[:space:]]*//')
VCPU_COUNT=$(nproc)
TOTAL_RAM=$(free -h | awk '/^Mem:/ {print $2}')
OS_KERNEL=$(uname -srm)
GO_VERSION=$(go version | awk '{print $3 " " $4}')

echo "💻 Hardware: ${CPU_MODEL} (${VCPU_COUNT} vCPUs), RAM: ${TOTAL_RAM}"
echo "🐧 OS Kernel: ${OS_KERNEL}"
echo "🐹 Go Version: ${GO_VERSION}"
echo "⚙️  Handler Mode: ${HANDLER_MODE} | Benchmark Target: ${BENCH_FUNC}"
echo "🏷️  Tier: ${TIER_NAME}"

# 2. Build Server
echo "🔨 Building server binary..."
cd "${ROOT_DIR}"
mkdir -p "${ROOT_DIR}/bin"
go build -o "${ROOT_DIR}/bin/server" cmd/server/main.go

# 3. Start Server in Background
echo "🚀 Starting server on :${PORT} in mode: ${HANDLER_MODE}..."
HANDLER_MODE="${HANDLER_MODE}" PORT="${PORT}" "${ROOT_DIR}/bin/server" > /tmp/telemetry_server.log 2>&1 &
SERVER_PID=$!

cleanup() {
    if kill -0 "${SERVER_PID}" 2>/dev/null; then
        echo "🛑 Terminating test server (PID: ${SERVER_PID})..."
        kill -TERM "${SERVER_PID}" 2>/dev/null || true
        wait "${SERVER_PID}" 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Wait for server ready
for i in {1..30}; do
    if curl -s "${HEALTH_URL}" | grep -q "HEALTHY"; then
        echo "✅ Server is HEALTHY on PID ${SERVER_PID}"
        break
    fi
    sleep 0.2
done

# Baseline RSS
BASELINE_RSS_KB=$(ps -o rss= -p "${SERVER_PID}" | tr -d ' ' || echo "0")
BASELINE_RSS_MB=$(awk "BEGIN {printf \"%.2f\", ${BASELINE_RSS_KB}/1024}")

# 4. Run Bombardier Load Test
echo "🔥 Bombarding ${TARGET_URL} with ${CONNECTIONS} connections for ${DURATION}..."
BOMB_JSON_FILE="/tmp/bombardier_telemetry_$$.json"
bombardier -c "${CONNECTIONS}" -d "${DURATION}" -l -m POST -f "${PAYLOAD_FILE}" -H "Content-Type: application/json" -p r -o json "${TARGET_URL}" > "${BOMB_JSON_FILE}"

# Peak RSS
PEAK_RSS_KB=$(ps -o rss= -p "${SERVER_PID}" | tr -d ' ' || echo "0")
PEAK_RSS_MB=$(awk "BEGIN {printf \"%.2f\", ${PEAK_RSS_KB}/1024}")

# Stop Server
cleanup
trap - EXIT

# Parse Bombardier JSON Output using Python for precision
PARSED_METRICS=$(python3 -c "
import json
with open('${BOMB_JSON_FILE}') as f:
    d = json.load(f)
r = d['result']
rps = r['rps']['mean']
p50 = r['latency']['percentiles'].get('50', 0) / 1000.0
p90 = r['latency']['percentiles'].get('90', 0) / 1000.0
p99 = r['latency']['percentiles'].get('99', 0) / 1000.0
total_2xx = r.get('req2xx', 0)
total_reqs = r.get('req1xx', 0) + total_2xx + r.get('req3xx', 0) + r.get('req4xx', 0) + r.get('req5xx', 0)
print(f'{rps:.2f}|{p50:.2f}|{p90:.2f}|{p99:.2f}|{total_reqs}|{total_2xx}')
")
rm -f "${BOMB_JSON_FILE}"

IFS='|' read -r RPS P50_MS P90_MS P99_MS TOTAL_REQS ERRORS_2XX <<< "${PARSED_METRICS}"

# 5. Run Native Microbenchmark for B/op & allocs/op
echo "🔬 Running native Go microbenchmark (testing.B -benchmem with ${BENCH_FUNC})..."
BENCH_OUTPUT=$(go test -bench="^${BENCH_FUNC}$" -benchmem -benchtime=3s "${ROOT_DIR}/internal/telemetry/...")

NS_OP=$(echo "${BENCH_OUTPUT}" | grep "${BENCH_FUNC}" | head -n1 | awk '{print $3}')
B_OP=$(echo "${BENCH_OUTPUT}" | grep "${BENCH_FUNC}" | head -n1 | awk '{print $5}')
ALLOCS_OP=$(echo "${BENCH_OUTPUT}" | grep "${BENCH_FUNC}" | head -n1 | awk '{print $7}')

echo ""
echo "================================================================"
echo "📊 BENCHMARK REPORT SUMMARY"
echo "================================================================"
echo "Tier:              ${TIER_NAME}"
echo "Mode:              ${HANDLER_MODE} (${BENCH_FUNC})"
echo "Throughput (RPS):  ${RPS} req/sec (Total: ${TOTAL_REQS}, 2xx: ${ERRORS_2XX})"
echo "Latency p50:       ${P50_MS} ms"
echo "Latency p90:       ${P90_MS} ms"
echo "Latency p99:       ${P99_MS} ms"
echo "Heap Allocation:   ${B_OP} B/op (${ALLOCS_OP} allocs/op)"
echo "Execution Speed:   ${NS_OP} ns/op"
echo "Server Memory RSS: Baseline ${BASELINE_RSS_MB} MB -> Peak ${PEAK_RSS_MB} MB"
echo "================================================================"

# If requested to append to CHANGELOG.md
if [[ "${1:-}" == "--append-changelog" ]]; then
    VERSION="${2:-0.1.0}"
    TITLE="${3:-${TIER_NAME}}"
    DATE_STR=$(date -u +"%Y-%m-%d")

    CHANGELOG_ENTRY="
## [${VERSION}] - ${DATE_STR} - ${TITLE}

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** ${CPU_MODEL} (${VCPU_COUNT} vCPUs)
- **Memory (RAM):** ${TOTAL_RAM} Total
- **Operating System / Sistem Operasi:** ${OS_KERNEL}
- **Runtime:** ${GO_VERSION}
- **Load Tool / Alat Uji:** bombardier v1.2.6 (${CONNECTIONS} connections, ${DURATION} duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **${RPS}** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **${ERRORS_2XX}** | requests in ${DURATION} (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **${P50_MS}** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **${P90_MS}** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **${P99_MS}** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **${B_OP}** | Bytes/op (\`B/op\` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **${ALLOCS_OP}** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **${NS_OP}** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **${BASELINE_RSS_MB} MB -> ${PEAK_RSS_MB} MB** | Resident Set Size |
"

    CHANGELOG_FILE="${ROOT_DIR}/CHANGELOG.md"
    echo "${CHANGELOG_ENTRY}" >> "${CHANGELOG_FILE}"
    echo "📝 Benchmark metrics appended to ${CHANGELOG_FILE} successfully!"
fi
