default: mint

mint: generator/*.go parser/*.go go.mod go.sum *.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@ ./cmd/

benchmarks/benchmarker.go benchmarks/benchmarker.custom.go: mint testdata/benchmarks/benchmark_types.mint
	./$< generate testdata/benchmarks --dest benchmarks/ --mkdir --package benchmarks
