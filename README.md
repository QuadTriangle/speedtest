# Speedtest Server

A self-contained speed test server. Download a single binary and run it — no dependencies needed.

## Quick Start

```bash
# Download the latest release for your platform
curl -fsSL https://github.com/QuadTriangle/speedtest/releases/latest/download/speedtest-linux-amd64 -o speedtest
chmod +x speedtest
./speedtest
```

Open http://localhost:8080 in your browser.

## Docker

```bash
# Linux (host networking, recommended)
docker run --rm -it --net=host ghcr.io/quadtriangle/speedtest:latest

# macOS / Windows (Docker Desktop)
docker run --rm -it -p 8080:8080 ghcr.io/quadtriangle/speedtest:latest

# Custom port
docker run --rm -it -p 3000:3000 ghcr.io/quadtriangle/speedtest:latest 3000
```

## Usage

```bash
./speedtest [port]
```

Default port is 8080.

Example:
```bash
./speedtest 8075
```

## Features

- Single binary with embedded web UI
- Measures download speed, upload speed, and latency
- Scales display from Kbps to Gbps automatically
- Zero configuration required
- Cross-platform (Linux, macOS, Windows)

## Building from source

```bash
go build -o speedtest .
```
