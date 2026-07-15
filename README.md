# cypher

A fast, concurrent TCP port scanner built as a command-line tool in Go.

`cypher` scans a range of ports on a host using a bounded worker pool, so it can check thousands of ports in seconds instead of minutes.

```
$ cypher -host scanme.nmap.org -start 1 -end 1024
22: Open
80: Open

scanned 1024 ports in 891ms, 2 open, 1022 closed⏎
```

## Features

- Concurrent scanning with a configurable worker pool
- Sorted, deterministic output
- Optional verbose mode to show closed ports
- Scan summary with timing and open/closed counts

## Installation

### Prerequisites

You need [Go](https://go.dev/dl/) 1.22 or later installed. Check your version:

```bash
go version
```

### Install directly (recommended)

```bash
go install github.com/YOUR_USERNAME/cypher@latest
```

This downloads, builds, and installs the `cypher` binary into your Go bin directory.

Make sure that directory is on your `PATH` — see the OS-specific steps below if the `cypher` command isn't found after installing.

### Build from source

```bash
git clone https://github.com/YOUR_USERNAME/cypher.git
cd cypher
go build -o cypher .
```

This produces a `cypher` (or `cypher.exe` on Windows) binary in the current directory. Move it anywhere on your `PATH` to run it from anywhere.

## Adding Go's bin directory to your PATH

`go install` places binaries in `$(go env GOPATH)/bin` (commonly `~/go/bin`). If `cypher` isn't found after installing, add that directory to your `PATH`.

### macOS / Linux (bash or zsh)

Add this line to `~/.zshrc` (zsh) or `~/.bashrc` (bash):

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell config:

```bash
source ~/.zshrc   # or source ~/.bashrc
```

### Fish shell

```fish
fish_add_path (go env GOPATH)/bin
```

This persists automatically — no config file editing or shell restart needed.

### Windows (PowerShell)

```powershell
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$(go env GOPATH)\bin", "User")
```

Restart your terminal for the change to take effect.

### Verify installation

```bash
which cypher      # macOS/Linux/fish
where.exe cypher  # Windows
```

Or just run it:

```bash
cypher -host scanme.nmap.org -start 20 -end 80
```

## Usage

```
cypher [flags]
```

| Flag | Default | Description |
|---|---|---|
| `-host` | `scanme.nmap.org` | Target host to scan |
| `-start` | `1` | First port in the range |
| `-end` | `1024` | Last port in the range |
| `-worker` | `100` | Number of concurrent workers |
| `-verbose` | `false` | Show closed ports in addition to open ones |

### Examples

Scan the default range on a host:

```bash
cypher -host example.com
```

Scan a specific range with more workers (faster on larger ranges):

```bash
cypher -host example.com -start 1 -end 65535 -worker 300
```

Show every port checked, not just open ones:

```bash
cypher -host example.com -start 20 -end 100 -verbose
```

Scan a local service:

```bash
cypher -host 127.0.0.1 -start 1 -end 1024
```

## Uninstallation

`go install` places a single binary in your Go bin directory — uninstalling just means removing that file.

### macOS / Linux / fish

```bash
rm $(go env GOPATH)/bin/cypher
```

### Windows (PowerShell)

```powershell
Remove-Item "$(go env GOPATH)\bin\cypher.exe"
```

Confirm it's gone:

```bash
which cypher   # macOS/Linux/fish — should show no output
```

If you also cloned the source separately (e.g. via `git clone`), remove that folder too:

```bash
rm -rf cypher   # wherever you cloned the repo
```

## Notes

- Only scan hosts you own or are explicitly authorized to test. [`scanme.nmap.org`](https://nmap.org/book/testing.html) is provided by the Nmap project specifically for testing scanners like this one.
- Scan timing depends heavily on network latency to the target and the number of workers — a local scan (`127.0.0.1`) will be dramatically faster than a scan against a remote host.

## License

MIT
