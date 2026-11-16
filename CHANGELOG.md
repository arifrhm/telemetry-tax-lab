# Changelog - Telemetry Tax Lab
## *Catatan Rilis & Telemetri Benchmark Beban Empiris*

All benchmarks are executed automatically with real hardware resource telemetry attached to each architectural release.  
*Semua benchmark dijalankan secara otomatis dengan data resource hardware riil terlampir pada setiap rilis arsitektur.*

---

## [0.1.0] - 2026-11-16 - Tier 0: Pure Baseline (Zero Telemetry)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **154999.59** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **1544923** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.59** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **1.47** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **4.46** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **0** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **0** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **910.7** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.64 MB -> 16.19 MB** | Resident Set Size |


## [0.2.0] - 2026-11-16 - Tier 1: Prometheus Metrics Tax (client_golang Counter & Histogram)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **159323.16** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **1580434** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.56** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **1.46** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **4.53** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **0** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **0** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **1323** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.43 MB -> 16.09 MB** | Resident Set Size |


## [0.3.0] - 2026-11-16 - Tier 2a: Structured Logging (rs/zerolog Zero-Alloc)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **136148.27** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **1360390** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.62** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **1.70** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **6.00** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **0** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **0** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **1805** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.46 MB -> 16.08 MB** | Resident Set Size |


## [0.3.1] - 2026-11-16 - Tier 2b: Structured Logging (log/slog Standard)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **131766.03** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **1317097** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.62** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **1.78** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **6.81** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **320** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **7** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **4689** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.62 MB -> 19.98 MB** | Resident Set Size |


## [0.4.0] - 2026-11-16 - Tier 3: OpenTelemetry Tracing Shock (100% AlwaysOn OTLP)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **65642.96** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **657020** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.91** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **2.65** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **6.86** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **5287** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **12** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **10082** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.66 MB -> 1479.69 MB** | Resident Set Size |


## [0.5.0] - 2026-11-16 - Tier 4: Intelligent Adaptive Sampling (OTel 1% Ratio Sampling)

### 💻 System & Environment Telemetry / Telemetri Lingkungan Sistem
- **Host CPU:** Intel(R) Xeon(R) Platinum 8173M CPU @ 2.00GHz (10 vCPUs)
- **Memory (RAM):** 23Gi Total
- **Operating System / Sistem Operasi:** Linux 6.8.0-147-generic x86_64
- **Runtime:** go1.26.0 linux/amd64
- **Load Tool / Alat Uji:** bombardier v1.2.6 (125 connections, 10s duration, ~1KB JSON payload)

### 📊 Benchmark Results / Hasil Pengujian Benchmark
| Metric Category / Kategori | Metric Name / Nama Metrik | Measured Value / Nilai Riil | Unit & Notes / Satuan & Catatan |
| :--- | :--- | :--- | :--- |
| **Throughput** | Mean Request Rate | **138533.94** | requests/sec (RPS) |
| **Total Volume** | Successful 2xx Requests | **1383232** | requests in 10s (0 errors / 100% 2xx) |
| **Tail Latency** | Median (p50) | **0.59** | milliseconds (ms) |
| **Tail Latency** | 90th percentile (p90) | **1.67** | milliseconds (ms) |
| **Tail Latency** | 99th percentile (p99) | **6.22** | milliseconds (ms) |
| **Memory Alloc** | Per-operation Heap Alloc | **942** | Bytes/op (`B/op` memory allocation) |
| **Memory Alloc** | Per-operation Allocs Count | **8** | allocs/op (heap allocations per request) |
| **Microbenchmark**| Single Operation Duration | **3965** | ns/op |
| **Process RSS** | Baseline -> Peak RAM | **11.44 MB -> 46.59 MB** | Resident Set Size |

