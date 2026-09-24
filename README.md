# golog

A small levelled logger for command-line Go programs.

- printf-style API: `Info`, `Warn`, `Error`, `Debug`, `Fatal`
- colourised level prefixes
- every entry reports the caller's `file:line`
- safe for concurrent use
- no configuration, no dependencies, no redaction rules

## Install

```bash
go get github.com/finnwu0x1/golog
```

## Usage

```go
package main

import "github.com/finnwu0x1/golog"

func main() {
	logger := golog.NewLogger()
	logger.Info("starting: listen=%s", "0.0.0.0:8080")
	logger.Error("request failed: %v", err)
}
```

Output:

```
2026/09/24 22:15:03 main.go:9: [I] starting: listen=0.0.0.0:8080
```

The location is the line that called the level method, so a logged error points
straight at the code that observed it.

Send output elsewhere (a file, a buffer in tests) with `NewLoggerWithWriter`:

```go
logger := golog.NewLoggerWithWriter(os.Stderr)
```

## Scope

This package logs what it is given. It does not redact secrets, classify
dependency errors, or filter messages; deciding what is safe to log belongs to
the caller.

`Fatal` writes the entry and then calls `os.Exit(1)`.

## License

Apache License 2.0. See [LICENSE](LICENSE).
