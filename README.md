# matelabdev-watcher-agent

A single static Go binary that runs on low-power devices (Raspberry Pi, Jetson, etc.) and actively monitors internal network hardware — cameras (RTSP/TCP), POS terminals, kiosks, barriers, and HTTP endpoints — reporting results to [Matelabdev Watcher](https://github.com/matelabdev/watcher).

## Quick Install

Copy the one-line install command from your project page in Watcher Master:

```bash
curl -sSL https://watcher.example.com/install/prj_xxxxxxxxxxxx | sudo sh
```

The script automatically:
1. Detects device architecture (`amd64` / `arm64` / `armv7`)
2. Downloads the correct binary
3. Places it at `/usr/local/bin/matelabdev-watcher-agent`
4. Writes the config to `/etc/matelabdev-watcher/config.yaml`
5. Installs and starts the systemd service

---

## Manual Install

### 1. Download Binary

Download the appropriate binary for your device from the [Releases](https://github.com/matelabdev/watcher-agent/releases) page:

| Target | File |
|--------|------|
| Linux x86-64 (server, PC) | `matelabdev-watcher-agent-linux-amd64` |
| Linux ARM 64-bit (RPi 4/5, Jetson) | `matelabdev-watcher-agent-linux-arm64` |
| Linux ARM 32-bit (RPi 2/3, older devices) | `matelabdev-watcher-agent-linux-armv7` |

```bash
# Example: Raspberry Pi 4
wget https://github.com/matelabdev/watcher-agent/releases/latest/download/matelabdev-watcher-agent-linux-arm64
chmod +x matelabdev-watcher-agent-linux-arm64
sudo mv matelabdev-watcher-agent-linux-arm64 /usr/local/bin/matelabdev-watcher-agent
```

### 2. Config File

```bash
sudo mkdir -p /etc/matelabdev-watcher
sudo nano /etc/matelabdev-watcher/config.yaml
```

```yaml
master_url: "https://watcher.example.com"
project_token: "prj_xxxxxxxxxxxx"
report_interval: 30        # seconds — used when a monitor has no per-monitor interval

monitors:
  - key: camera-entrance
    name: Entrance Camera
    type: rtsp
    host: 192.168.1.100    # host only → defaults to port 554
    timeout: 5

  - key: pos-terminal
    name: POS Terminal
    type: tcp
    host: 192.168.1.200
    port: 9100
    timeout: 5
    interval: 60           # per-monitor interval (seconds)

  - key: payment-api
    name: Payment API
    type: http
    url: "https://payment.example.com/ping"
    timeout: 10
    interval: 120
```

### 3. Run

```bash
matelabdev-watcher-agent -config /etc/matelabdev-watcher/config.yaml
```

---

## systemd Service

```bash
sudo nano /etc/systemd/system/matelabdev-watcher-agent.service
```

```ini
[Unit]
Description=Matelabdev Watcher Agent
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/usr/local/bin/matelabdev-watcher-agent -config /etc/matelabdev-watcher/config.yaml
Restart=always
RestartSec=10
User=root

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable matelabdev-watcher-agent
sudo systemctl start matelabdev-watcher-agent
sudo systemctl status matelabdev-watcher-agent
```

---

## Monitor Types

### `tcp` — TCP Port Check

For hardware accessible over TCP: POS terminals, kiosks, barriers.

```yaml
- key: pos-terminal
  type: tcp
  host: 192.168.1.200
  port: 9100
  timeout: 5
```

### `rtsp` — Camera (RTSP / TCP:554)

IP cameras. Establishes a TCP connection to port 554 and checks the banner.

```yaml
- key: camera-entrance
  type: rtsp
  host: 192.168.1.100     # host only → port 554 is used
  timeout: 5
```

Use `host: 192.168.1.100:8554` to specify a non-standard port.

### `http` — HTTP Endpoint Check

Payment APIs, service health endpoints.

```yaml
- key: payment-api
  type: http
  url: "https://payment.example.com/ping"
  timeout: 10
```

| HTTP Response | Status |
|---------------|--------|
| 1xx–4xx | `online` |
| 5xx | `degraded` |
| Connection error / timeout | `offline` |

---

## Config Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `master_url` | string | — | Watcher Master URL (required) |
| `project_token` | string | — | Project token starting with `prj_` (required) |
| `report_interval` | int | `30` | Default check interval for all monitors (seconds) |
| `monitors[].key` | string | — | Unique identifier (required) |
| `monitors[].name` | string | — | Human-readable label |
| `monitors[].type` | string | — | `tcp`, `http`, or `rtsp` |
| `monitors[].host` | string | — | Target IP or hostname (tcp/rtsp) |
| `monitors[].port` | int | — | Target port (tcp) |
| `monitors[].url` | string | — | Target URL (http) |
| `monitors[].timeout` | int | `5` | Connection timeout (seconds) |
| `monitors[].interval` | int | `report_interval` | Per-monitor check interval |

---

## Build from Source

```bash
git clone https://github.com/matelabdev/watcher-agent
cd watcher-agent

# Local binary
make build

# Cross-compile for all platforms
make build-all
# → dist/matelabdev-watcher-agent-linux-amd64
# → dist/matelabdev-watcher-agent-linux-arm64
# → dist/matelabdev-watcher-agent-linux-armv7

# Tests
make test
```

**Requirements:** Go 1.22+

---

## Architecture

```
config.yaml
    │
    ▼
main.go ──► reporter.SyncMonitors() ──► POST /api/monitors/sync
    │
    ▼
scheduler ──► [goroutine + ticker per monitor]
    │
    ├── checker.New(tcp)  ──► net.DialTimeout
    ├── checker.New(rtsp) ──► TCP:554 connection
    └── checker.New(http) ──► net/http GET
              │
              ▼
        reporter.SendHeartbeat() ──► POST /api/heartbeat
```

All notification decisions are made by the Master — the agent only reports raw status.

---

## License

MIT