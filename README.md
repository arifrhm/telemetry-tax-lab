# ⚡ Telemetry Tax Lab (`telemetry-tax-lab`)
> *"Quantifying the Telemetry Tax: An Empirical Benchmark of OpenTelemetry, Prometheus, Structured Logging & Adaptive Sampling in High-Throughput Go Systems."*  
> *"Mengukur 'Pajak Telemetri': Benchmark empiris beban OpenTelemetry, Prometheus, Structured Logging & Adaptive Sampling pada sistem Go performa tinggi."*

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI Status](https://github.com/arifrhm/telemetry-tax-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/arifrhm/telemetry-tax-lab/actions)

---

## 🌍 Language / Bahasa
* [English](#-english)
* [Bahasa Indonesia](#-bahasa-indonesia)

---

## 🇬🇧 English

### 🎯 Manifesto & Lab Philosophy
In cloud-native distributed systems, engineering teams are repeatedly told to adopt **OpenTelemetry (OTel)**, Prometheus metrics, and structured logging for full observability. However, observability is not free. In high-throughput architectures, it introduces **The Telemetry Tax**:
1. **Throughput Collapse:** CPU cycles diverted from business logic to span context serialization.
2. **Tail Latency Spikes:** Mutex contention across concurrent metric buckets and trace buffers degrading p99 latencies.
3. **Memory & Cloud Bill Churn:** Unbounded heap allocations causing severe Garbage Collection pressure and astronomical Datadog/CloudWatch ingestion invoices.

**Telemetry Tax Lab** is an open-source, reproducible research environment that measures this tax across 5 architectural tiers on a single commodity Linux node with 100% synthetic mock telemetry.

### 💻 Hardware & Environment Benchmark Lab
All benchmarks in this series are executed automatically against real hardware:

| Resource | Specification |
| :--- | :--- |
| **CPU** | Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs) |
| **Memory (RAM)** | 23 GiB Total |
| **Operating System** | Linux 6.8.0-147-generic x86_64 (Ubuntu) |
| **Runtime** | `go1.26.0 linux/amd64` |
| **Load Tool** | `bombardier v1.2.6` (125 concurrent connections, 10s duration) |
| **Target Endpoint** | `POST /events` (~1KB JSON Ingestion Payload) |

### 📈 Empirical Evolution Matrix: The Telemetry Tax

| Tier | Architectural Description | Throughput (RPS) | Tail Latency p99 | Heap Alloc (B/op) | Allocs/Op | Micro Speed | Telemetry Tax (% Drop) |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Tier 0: Baseline** | Pure FastHTTP Ingestion (Zero Telemetry) | **154,999 req/s** | **4.46 ms** | **0 B/op** | **0 allocs** | **910 ns/op** | **0.0% (Ground Truth)** |
| **Tier 1: Prometheus**| `client_golang` Counter & Histogram Vec | **159,323 req/s** | **4.53 ms** | **0 B/op** | **0 allocs** | **1,323 ns/op** | **+2.8% (Negligible)** |
| **Tier 2a: Zerolog** | Zero-Allocation Structured JSON Logger | **136,148 req/s** | **6.00 ms** | **0 B/op** | **0 allocs** | **1,805 ns/op** | **-12.1% (Mild Tax)** |
| **Tier 2b: Slog** | Standard Library `log/slog` JSON Logger | **131,766 req/s** | **6.81 ms** | **320 B/op** | **7 allocs** | **4,689 ns/op** | **-15.0% (Heap Escape)** |
| **Tier 3: OTel 100%** | OpenTelemetry Tracing (100% AlwaysOn) | **65,642 req/s** | **6.86 ms** | **5,287 B/op** | **12 allocs** | **10,082 ns/op** | **-57.6% (Tracing Shock! 💥)** |
| **Tier 4: OTel 1%** | Intelligent Adaptive Ratio Sampling (1%)| **138,533 req/s** | **6.22 ms** | **942 B/op** | **8 allocs** | **3,965 ns/op** | **-10.6% (90% Recovery! 🚀)** |

---

## 📊 Visual ASCII Comparisons

### 1. Throughput Comparison (Higher is Better)
```
Tier 0 (Pure Baseline)    [██████████████████████████████████████]  154,999 RPS (Baseline)
Tier 1 (Prometheus)       [██████████████████████████████████████]  159,323 RPS (+2.8%)
Tier 2a (Zerolog)         [███████████████████████████████       ]  136,148 RPS (-12.1%)
Tier 2b (Slog Standard)   [█████████████████████████████         ]  131,766 RPS (-15.0%)
Tier 3 (OTel 100% Trace)  [██████████████                        ]   65,642 RPS (-57.6% 💥)
Tier 4 (OTel 1% Sampling) [███████████████████████████████       ]  138,533 RPS (+111% vs Tier 3 🚀)
                          +--------------------------------------+
                          0k     40k     80k    120k   160k
```

### 2. Heap Memory Allocation per Request (Lower is Better)
```
Tier 0 (Pure Baseline)    [                                      ]      0 B/op
Tier 1 (Prometheus)       [                                      ]      0 B/op
Tier 2a (Zerolog)         [                                      ]      0 B/op
Tier 2b (Slog Standard)   [██                                    ]    320 B/op
Tier 3 (OTel 100% Trace)  [██████████████████████████████████████]  5,287 B/op  (Heap Explosion!)
Tier 4 (OTel 1% Sampling) [██████                                ]    942 B/op  (-82.2% drop 📉)
                          +--------------------------------------+
                          0 B   1,000 B  2,500 B  5,000 B
```

### 🔬 Interactive Vector SVG Flamegraphs
Inspect where runtime CPU and memory are spent across tiers:
* **Heap Profiles:** [Baseline Heap Profile (0 B/op)](./docs/flamegraphs/baseline_heap.svg) vs [OTel 100% Heap Profile (5.3 KB)](./docs/flamegraphs/otel100_heap.svg) vs [OTel 1% Ratio Heap Profile](./docs/flamegraphs/otel1_heap.svg)
* **CPU Profiles:** [Baseline CPU Execution Profile](./docs/flamegraphs/baseline_cpu.svg) vs [OTel 100% Tracing CPU Profile](./docs/flamegraphs/otel100_cpu.svg) vs [OTel 1% CPU Profile](./docs/flamegraphs/otel1_cpu.svg)

---

## 🇮🇩 Bahasa Indonesia

### 🎯 Manifesto & Filosofi Lab
Dalam rekayasa sistem terdistribusi modern, tim engineering sering diwajibkan memasang **OpenTelemetry (OTel)**, metrik Prometheus, dan structured logging. Namun, *observability* tidak gratis. Pada arsitektur dengan beban volume tinggi, muncul fenomena **The Telemetry Tax (Pajak Telemetri)**:
1. **Penurunan Throughput:** Siklus CPU yang seharusnya memproses logika bisnis teralihkan untuk serialisasi atribut span dan context propagation.
2. **Lonjakan Latensi Ekor (Tail Latency):** Perebutan kunci mutex antar-koneksi konkurensi pada bucket histogram dan buffer trace yang merusak kestabilan latensi p99.
3. **Ledakan Memori & Biaya Cloud:** Alokasi heap yang membengkak memicu tekanan Garbage Collection (GC) berat dan tagihan invoice Datadog/CloudWatch yang melonjak drastis.

**Telemetry Tax Lab** mengukur fenomena ini secara terbuka dan berbasis angka empiris riil di atas satu mesin Linux komoditas menggunakan 100% data sintetis (bebas dari pelanggaran NDA kantor).

---

## 🚀 How to Run & Reproduce / Cara Menjalankan

### 1. Prerequisites / Prasyarat
* [Go 1.26+](https://golang.org)
* [Bombardier](https://github.com/codesenberg/bombardier) (`go install github.com/codesenberg/bombardier@latest`)
* `curl` and `python3`

### 2. Run Local CI Suite / Menjalankan Pengujian CI Lokal
```bash
make ci
# Runs gofmt check, go vet, concurrency race detector (-race), and lifecycle smoke test
```

### 3. Run Automated Benchmark / Menjalankan Otomasi Benchmark
```bash
# Run Baseline (Zero Telemetry)
HANDLER_MODE=baseline make bench

# Run Prometheus Metrics
HANDLER_MODE=prometheus make bench

# Run OpenTelemetry Tracing 100% AlwaysOn
HANDLER_MODE=otel-trace-100 make bench

# Run OpenTelemetry 1% Adaptive Ratio Sampling
HANDLER_MODE=otel-trace-1 make bench
```

---

## 📁 Repository Structure / Struktur Repositori
* `cmd/server/main.go`: Server HTTP dengan router mode telemetri (`baseline`, `prometheus`, `zerolog`, `slog`, `otel-trace-100`, `otel-trace-1`).
* `internal/telemetry/`: Implementasi handler berinstrumen telemetri beserta microbenchmark per-tier.
* `scripts/bench.sh`: Otomasi pengujian beban `bombardier`, pengukuran RSS memori proses, dan generator Markdown `CHANGELOG.md`.
* `scripts/ci_test.sh`: Script validator CI lokal (gofmt, vet, race detector, smoke test).
* `docs/flamegraphs/`: Aset visual interaktif SVG CPU & Heap Flamegraph dari runtime pprof.

---
*Built for open software engineering observability performance experimentation.*  
*Dibuat untuk eksperimen performa observability rekayasa perangkat lunak terbuka.*
