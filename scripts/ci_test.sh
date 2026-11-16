#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "================================================================"
echo "🛡️  [Telemetry Tax Lab] Local CI & Automated Verification"
echo "================================================================"

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "\n${BLUE}🔍 [Stage 1/4] Code Formatting, Vet & Compile Check...${NC}"
cd "${ROOT_DIR}"

UNFORMATTED=$(gofmt -l . | grep -v "^vendor" || true)
if [[ -n "${UNFORMATTED}" ]]; then
    echo -e "${RED}❌ Unformatted Go files detected:${NC}"
    echo "${UNFORMATTED}"
    exit 1
fi
echo -e "${GREEN}✅ Formatting: All Go files properly formatted.${NC}"

go vet ./...
echo -e "${GREEN}✅ Go Vet: Static analysis clean.${NC}"

go build -o /dev/null ./cmd/server/...
echo -e "${GREEN}✅ Build Check: Server compiles cleanly.${NC}"

echo -e "\n${BLUE}🧪 [Stage 2/4] Concurrency Race Detector & Unit Tests...${NC}"
go test -v -race -timeout=30s ./internal/...
echo -e "${GREEN}✅ Unit Tests & Race Detection: 100% PASS with zero data races.${NC}"

echo -e "\n${BLUE}⚡ [Stage 3/4] Telemetry Microbenchmarks Check...${NC}"
go test -bench='^Benchmark(Baseline|Prometheus|Logging_Zerolog|OTel_Tracing_1Percent)' -benchmem -benchtime=1s ./internal/telemetry/...
echo -e "${GREEN}✅ Microbenchmarks: Successfully executed all target tiers.${NC}"

echo -e "\n${BLUE}🚀 [Stage 4/4] Smoke Testing Server Lifecycle & Routing...${NC}"
TEST_PORT="8998"
mkdir -p "${ROOT_DIR}/bin"
go build -o "${ROOT_DIR}/bin/ci_server" cmd/server/main.go

HANDLER_MODE="baseline" PORT="${TEST_PORT}" "${ROOT_DIR}/bin/ci_server" > /tmp/ci_telemetry_server.log 2>&1 &
SERVER_PID=$!

cleanup() {
    if kill -0 "${SERVER_PID}" 2>/dev/null; then
        kill -TERM "${SERVER_PID}" 2>/dev/null || true
        wait "${SERVER_PID}" 2>/dev/null || true
    fi
    rm -f "${ROOT_DIR}/bin/ci_server"
}
trap cleanup EXIT

HEALTH_URL="http://127.0.0.1:${TEST_PORT}/health"
READY=0
for i in {1..20}; do
    if curl -s "${HEALTH_URL}" | grep -q "HEALTHY"; then
        READY=1
        break
    fi
    sleep 0.1
done

if [[ "${READY}" -ne 1 ]]; then
    echo -e "${RED}❌ Server failed to respond on ${HEALTH_URL}! Log:${NC}"
    cat /tmp/ci_telemetry_server.log
    exit 1
fi
echo -e "${GREEN}✅ Health Check: Server started and responded HEALTHY on port ${TEST_PORT}.${NC}"

PAYLOAD_FILE="${ROOT_DIR}/scripts/payload_sample.json"
RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "http://127.0.0.1:${TEST_PORT}/events" \
    -H "Content-Type: application/json" \
    -d @"${PAYLOAD_FILE}")

HTTP_STATUS=$(echo "${RESP}" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "${RESP}" | grep -v "HTTP_STATUS")

if [[ "${HTTP_STATUS}" != "200" ]]; then
    echo -e "${RED}❌ Ingestion test failed! Status: ${HTTP_STATUS}, Body: ${BODY}${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Event Ingestion: Payload successfully ingested (HTTP 200 OK).${NC}"

cleanup
trap - EXIT

echo -e "\n================================================================"
echo -e "${GREEN}🎉 ALL LOCAL CI CHECKS & VERIFICATIONS PASSED!${NC}"
echo -e "================================================================\n"
