# Performance

Measured on 2026-11-25, `linux/amd64`, Go 1.27.4, AMD Ryzen 7 9700X (16 logical CPUs).
The `Pack`, converter, MGF and TXT benchmarks use the real `testdata/gen` and
`testdata/txt` fixtures. The `KnapsackPack` microbenchmark uses synthetic
items to isolate the DP. Numbers are machine-specific.

## Commands

```bash
# Packer, knapsack and default best-fit-decreasing
go test -run '^$' -bench 'Benchmark(KnapsackPack|Pack|PackBestFit)$' \
  -benchmem -count=6 ./pkg/binpack ./internal/nesting

# MGF rendering and TXT parsing
go test -run '^$' -bench 'Benchmark(Format|ParseFile)$' \
  -benchmem -count=6 ./internal/mgf ./internal/source/txt

# End-to-end converter (GEN fixtures, reports, Excel, remnant CSV)
go test -run '^$' -bench BenchmarkConvert10510Gen \
  -benchmem -count=6 ./internal/app

# CPU scaling of the packer (single-threaded by design)
go test -run '^$' -bench BenchmarkPackBestFit -cpu=1,2,4,8 \
  ./internal/nesting
```

## Results

| Benchmark | Before | After |
| --- | --- | --- |
| `Pack` (239 GEN profiles, knapsack) | 32.98 ms, 208.4 MB, 12,184 allocs | 8.14 ms, 2.48 MB, 2,589 allocs |
| `KnapsackPack/items=84` | 16.65 ms, 130.8 MB, 1,810 allocs | 4.79 ms, 237 KB, 13 allocs |
| `KnapsackPack/items=239` | 146.0 ms, 1.015 GB, 11,874 allocs | 35.7 ms, 482 KB, 19 allocs |
| `PackBestFit` (default, 239 profiles) | — | **0.303 ms, 1.46 MB, 2,587 allocs** |
| `Convert10510Gen` (end-to-end) | — | **10.8 ms, 10.6 MB, 92.5k allocs** |
| `Format` (71 MGF nests) | — | 1.24 ms, 1.44 MB, 30.2k allocs |
| `ParseFile` (5 TXT files) | — | 0.44 ms, 734 KB, 6.1k allocs |

The rerun used the current Go 1.27.4 toolchain and dependency versions. The
direct runtime dependencies are already current (`aveva-go v0.1.1`,
`dxf v1.0.2`, `excelize v2.11.1`); only tool-only indirect modules
(`x/mod`, `x/tools`) have newer releases. The differences from the previous
table are within run noise (about 1–5%); only the BFD path improved slightly.
There is no measured evidence that the dependency updates alone made the
pipeline materially faster.

An earlier `benchstat` run on the 239-profile `Pack` benchmark measured
**-75.51% time**, **-98.89%** allocated bytes and **-83.70%** allocations
(`p=0.002`, `n=6`).

The default `best-fit-decreasing` path is about **27× faster** than the
optimized knapsack path and uses less memory. The domain aggregates deep-copy
parts and cutouts on read, which costs roughly 25% on the short BFD path; this
is the price of not leaking aggregate state and is still sub-millisecond. It
also improved converter quality on the 9 GEN fixtures: **66 nests at 91.65%
utilization** versus 67 nests at 90.29% for knapsack. Knapsack remains available with
`--algorithm knapsack`.

## Why knapsack is not globally optimal

The knapsack DP is exact for one bin: it finds the subset with the largest
packed length that fits the bar. Bin packing itself is NP-hard, and applying
that per-bin optimum sequentially is not a global optimum. A maximal first
bin can strand items that a global strategy would combine differently.

A small concrete case, encoded by
`TestKnapsackPerBinOptimalIsNotGlobalOptimal`: capacity 30, items
`[21, 20, 18, 11, 8, 4]`.

- Per-bin knapsack picks `18 + 8 + 4 = 30` first, then must use separate bars
  for 21, 20 and 11: **4 bars**.
- Best-fit decreasing packs `21 + 8`, `20 + 4`, `18 + 11`: **3 bars**.

On the 9 GEN fixtures this shows up as 67 knapsack bars versus 66 best-fit
bars. Knapsack is still available with `--algorithm knapsack`; it remains the
right choice when the per-bin maximum packed length is the requirement.

## What changed

1. **Knapsack DP memory**: replaced the `n × capacity` integer matrix with a
   single rolling row plus a packed decision bitset. The objective and
   inclusion rule are unchanged; selected-item reconstruction is done from the
   bitset.
2. **Workspace reuse**: `KnapsackAlgorithm.Pack` now allocates its DP,
   bitset, length, selection and remaining-item buffers once per pack call and
   reuses them for every bin.
3. **Spacing cache**: endcut spacing is resolved once per
   `(profile type, endcut)` pair instead of re-parsing the profile type for
   every end of every profile.
4. **Default algorithm**: `best-fit-decreasing` is the converter default;
   `knapsack` stays selectable for exact per-bin objective comparisons.

The CPU profile before the change put 57% of packer time in
`KnapsackAlgorithm.solveKnapsack`; after workspace reuse the remaining
knapsack cost is the DP loop itself, not allocation or GC. The default
best-fit path avoids that loop entirely.

## Concurrency

The pipeline spawns no goroutines. The only concurrency primitive on this path
is the algorithm registry's `sync.RWMutex`, which is not touched by
`conv`/`parser`. `go test -race ./...` passes.

`BenchmarkPackBestFit` varies only slightly at `GOMAXPROCS=1,2,4,8`
(about 0.31–0.32 ms/op), so parallelising the packer would add coordination without throughput. The
current 9 GEN files and 5 TXT files also finish end-to-end in about 11 ms;
file-level worker pools are not justified yet. If inputs grow to hundreds of
files, the independent source files are the first candidate for bounded
parallel parsing, not the packer.

## Remaining costs

The converter profile is now dominated by dependency and report work:

- `aveva-go` scanner/contour parsing (its `bufio.Scanner.Text` dominates
  allocation objects);
- Excel XML generation in `xuri/excelize`;
- report string building.

These are outside the hot algorithm and currently fit comfortably in the
11 ms end-to-end budget. The Excel writer could be moved to `StreamWriter` if
workbooks grow, but it is not a present bottleneck.