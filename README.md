This repo demonstrates a minimal HTTP server that contains a data race if the `-unsafe` flag is set.

To trigger the data race run `go test -race . -run TestUnsafeHandler`. This will (usually…) print out a report showing exactly what lines of code were running in each goroutine when the data race was detected.

See https://go.dev/doc/articles/race_detector for more details about the Go's data race detector.
