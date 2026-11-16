# 🧵 X / Twitter Thread: The Tracing Shock
## *Why OpenTelemetry Cut Our Go Ingestion by 57% and How 1% Sampling Fixed It*
## *Mengapa OpenTelemetry Memotong Throughput Go Ingestion Sebesar 57% dan Bagaimana 1% Sampling Memperbaikinya*

---

### Tweet 1: The Hook / Pembuka
**[EN]**
We turned on OpenTelemetry distributed tracing in our high-throughput Go ingestion pipeline.
Result: Throughput crashed from **155,000 req/s down to 65,000 req/s (-57.6%)**.
Memory allocation exploded to **5.3 KB per single HTTP request**.

Here is the empirical breakdown of "The Telemetry Tax" and how 1% sampling recovered +111% throughput 🧵👇

**[ID]**
Kami menyalakan OpenTelemetry distributed tracing pada pipeline ingestion Go berkecepatan tinggi.
Hasilnya: Throughput anjlok dari **155.000 req/s ke 65.000 req/s (-57,6%)**.
Alokasi memori meledak menjadi **5,3 KB per satu HTTP request**.

Berikut bedah empiris "Pajak Telemetri" dan bagaimana sampling 1% memulihkan throughput sebesar +111% 🧵👇

---

### Tweet 2: The Benchmark Matrix / Matriks Benchmark
**[EN]**
We benchmarked 5 tiers under identical hardware (10 vCPUs Xeon, 23GB RAM, 125 concurrent connections, bombardier):

0️⃣ Pure Baseline: **154,999 RPS** (0 B/op, 0 allocs)
1️⃣ Prometheus: **159,323 RPS** (0 B/op, 0 allocs, +413 ns)
2️⃣ Zerolog: **136,148 RPS** (0 B/op, 0 allocs)
3️⃣ OTel AlwaysOn (100%): **65,642 RPS** (5,287 B/op, 12 allocs) 📉
4️⃣ OTel Adaptive (1%): **138,533 RPS** (942 B/op, 8 allocs) 🚀

Receipts & benchmarks: github.com/arifrhm/telemetry-tax-lab

**[ID]**
Kami menguji 5 tingkatan pada hardware identik (10 vCPUs Xeon, 23GB RAM, 125 concurrent connections, bombardier):

0️⃣ Pure Baseline: **154.999 RPS** (0 B/op, 0 allocs)
1️⃣ Prometheus: **159.323 RPS** (0 B/op, 0 allocs, +413 ns)
2️⃣ Zerolog: **136.148 RPS** (0 B/op, 0 allocs)
3️⃣ OTel AlwaysOn (100%): **65.642 RPS** (5.287 B/op, 12 allocs) 📉
4️⃣ OTel Adaptive (1%): **138.533 RPS** (942 B/op, 8 allocs) 🚀

Bukti & benchmark: github.com/arifrhm/telemetry-tax-lab

---

### Tweet 3: Prometheus is Near-Zero Cost / Prometheus Nyaris Gratis
**[EN]**
Myth: "Metrics slow down high-concurrency microservices."
Reality: CNCF `client_golang` adds only **+413 nanoseconds** overhead.
Zero heap allocations (0 B/op).
Why? Prometheus uses lock-free atomic 64-bit integer counters in memory. Scraping happens out-of-band on `/metrics`. Metrics are essentially free.

**[ID]**
Mitos: "Metrik memperlambat microservice berkonsentrasi tinggi."
Realita: `client_golang` CNCF hanya menambah latensi **+413 nanodetik**.
Nol alokasi heap (0 B/op).
Mengapa? Prometheus memakai atomic 64-bit counter di memory. Scraping terjadi out-of-band via `/metrics`. Metrik pada dasarnya gratis.

---

### Tweet 4: Zerolog vs Slog / Perbandingan Logging
**[EN]**
Logging is where heap allocations begin:
- `rs/zerolog`: **136k RPS**, **0 B/op, 0 allocs**, 1,805 ns/op.
- Standard `log/slog`: **131k RPS**, **320 B/op, 7 allocs**, 4,689 ns/op.

Slog is 2.6x slower in execution microbenchmark due to `slog.Attr` boxing interface allocations. At 100k RPS, Slog allocates **32 MB/sec of garbage** for the GC to collect!

**[ID]**
Logging adalah titik awal alokasi heap:
- `rs/zerolog`: **136k RPS**, **0 B/op, 0 allocs**, 1.805 ns/op.
- Bawaan `log/slog`: **131k RPS**, **320 B/op, 7 allocs**, 4.689 ns/op.

Slog 2,6x lebih lambat di microbenchmark karena boxing interface pada `slog.Attr`. Pada 100k RPS, Slog memproduksi sampah memori **32 MB/detik** yang membebani Go Garbage Collector!

---

### Tweet 5: The Anatomy of "The Tracing Shock" / Anatomi Tracing Shock
**[EN]**
Why does OpenTelemetry 100% AlwaysOn cut throughput by 57%?
1. Span Context propagation: Hex encoding/decoding `traceparent` headers.
2. Synchronous Span Lifecycle: Allocating span attributes, events, timestamps on the hot path.
3. Heap bloat: **5.3 KB per request** triggered aggressive GC stop-the-world cycles.
Micro-latency jumped from **910 ns to 10,082 ns/op (11x slower)**!

**[ID]**
Mengapa OpenTelemetry 100% AlwaysOn memangkas throughput hingga 57%?
1. Propagasi Span Context: Serialisasi hex header `traceparent`.
2. Lifecycle Span Sinkron: Alokasi atribut span, event, dan timestamp di hot-path.
3. Ledakan Heap: **5,3 KB per request** memicu siklus GC stop-the-world yang intens.
Latensi mikro melonjak dari **910 ns ke 10.082 ns/op (11x lebih lambat)**!

---

### Tweet 6: The Fix: Adaptive Ratio-based Sampling / Solusi: Adaptive Sampling
**[EN]**
The cure isn't disabling tracing. It's **TraceIDRatioBased(0.01)** (1% sampling):
- Throughput surged back from **65k to 138k RPS (+111% recovery)**.
- Heap allocations dropped from **5,287 B/op to 942 B/op (-82.2% reduction)**.
- p99 latency stabilized at **6.22 ms**.
99% of requests pass through a fast-path no-op span context without allocating memory.

**[ID]**
Solusinya bukan mematikan tracing. Melainkan **TraceIDRatioBased(0.01)** (1% sampling):
- Throughput melonjak dari **65k kembali ke 138k RPS (+111% pemulihan)**.
- Alokasi heap terpangkas dari **5.287 B/op menjadi 942 B/op (-82,2% penurunan)**.
- Latensi p99 stabil di **6,22 ms**.
99% request melewati fast-path no-op span context tanpa mengalokasikan memori.

---

### Tweet 7: Conclusion & Open Source Lab / Kesimpulan & Lab Terbuka
**[EN]**
Key Takeaways for SREs & Platform Engineers:
1. Metrics = Safe anywhere.
2. Choose zero-alloc loggers (`zerolog`) for ingest paths.
3. NEVER run OTel tracing at 100% in high-RPS services unless you want to pay a 57% compute tax and huge Datadog bills. Use 1% ratio sampling.

Explore code, flamegraphs, and local CI scripts:
👉 https://github.com/arifrhm/telemetry-tax-lab
Star ⭐️ the repo if you build high-scale Go systems!

**[ID]**
Pelajaran Kunci untuk SRE & Platform Engineer:
1. Metrik = Aman di segala layer.
2. Gunakan zero-alloc logger (`zerolog`) untuk jalur ingest.
3. JANGAN PERNAH menyalakan OTel tracing 100% pada service ber-RPS tinggi kecuali Anda siap membayar penalti komputasi 57% dan tagihan Datadog/NewRelic membengkak. Pakai 1% ratio sampling.

Eksplorasi kode, flamegraphs, dan skrip CI lokal:
👉 https://github.com/arifrhm/telemetry-tax-lab
Beri bintang ⭐️ jika Anda mendesain sistem Go berskala besar!
