# Rust vs. Go

This repo was forked from [https://github.com/misikdmytro/text-similarity](https://github.com/misikdmytro/text-similarity).

Comparison of Rust and Go in terms of performance. This fork has optimized and rewritten to Go solution for improved performance.

Thanks to [Dmytro Misik](https://github.com/misikdmytro) for the article on Medium which inspired me to optimize the Go code, and for open-sourcing the underlying source code.

## Optimizations

* Text sanitation done using https://github.com/eriklupander/replacer instead of stdlib `regexp` package. `strings.Replacer` from stdlib works almost as well.
* JSON parsing done using `buger/jsonparser` instead of `encoding/json`.
* "Set" code for unique words rewritten with `strset` and more effective algorithm.
* Specified initial capacity of `map` instances, i.e. `make(map[string]int, someLength)`
* Eliminated unnecessary for-loop with `map` allocation in `calculateIDF` and pre-computed scores.
* Slightly more efficient implementation of `calculateTF` with pre-computed fraction multiplication replacing division.

The largest gains was text sanitation (doubling the throughput) and (more surprisingly) setting initial capacities when creating `map` instances.

## Performance figures
All benchmarks from my machine using Go 1.23.4. Rust binary built with --release flag.

```
Original Go implementation                  : 2266 req/s, avg: 168 ms.
Original Rust implementation                : 5045 req/s, avg: 74 ms.
Optimized Go implementation                 : 6882 req/s, avg: 43.5 ms.
```
