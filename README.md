# ipmitool-api

A simple REST API wrapper around `ipmitool` originally built for use within Home Assistant but could be re-used elsewhere.

## Features

- **Power control** — on, off, reset, cycle, soft
- **Sensor readings** — temperature, fan speed, voltage
- **Chassis status** — power state, faults, last power event
- **API key authentication** — via `X-API-Key` header or `Authorization: Bearer <key>`
- **Config file** — define your IPMI servers in YAML

## Quick Start

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your IPMI hosts and a secure API key

go build -o ipmitool-api .
./ipmitool-api -config config.yaml
```

## Docker

```bash
docker build -t ipmitool-api .
docker run -d \
  -p 8080:8080 \
  -v /path/to/config.yaml:/etc/ipmitool-api/config.yaml:ro \
  ipmitool-api
```

## API Endpoints

All endpoints require the `X-API-Key` header (or `Authorization: Bearer <key>`).

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/servers` | List configured servers |
| `GET` | `/api/servers/{name}/power` | Get power status |
| `POST` | `/api/servers/{name}/power` | Set power state |
| `GET` | `/api/servers/{name}/sensors` | Get sensor readings |
| `GET` | `/api/servers/{name}/chassis` | Get chassis status |

### Examples

```bash
# List servers
curl -H "X-API-Key: your-key" http://localhost:8080/api/servers

# Get power status
curl -H "X-API-Key: your-key" http://localhost:8080/api/servers/server1/power

# Power on
curl -X POST -H "X-API-Key: your-key" \
  -d '{"action":"on"}' \
  http://localhost:8080/api/servers/server1/power

# Get sensors
curl -H "X-API-Key: your-key" http://localhost:8080/api/servers/server1/sensors

# Get chassis status
curl -H "X-API-Key: your-key" http://localhost:8080/api/servers/server1/chassis
```

### Power Actions

| Action | Description |
|--------|-------------|
| `on` | Power on |
| `off` | Hard power off |
| `soft` | Graceful shutdown (ACPI) |
| `reset` | Hard reset |
| `cycle` | Power cycle (off then on) |

## Home Assistant Integration

Use the [RESTful](https://www.home-assistant.io/integrations/rest/) integration:

```yaml
rest_command:
  server_power_on:
    url: "http://ipmitool-api:8080/api/servers/server1/power"
    method: POST
    headers:
      X-API-Key: "your-key"
    payload: '{"action":"on"}'
    content_type: "application/json"

sensor:
  - platform: rest
    name: "Server1 Power"
    resource: "http://ipmitool-api:8080/api/servers/server1/power"
    headers:
      X-API-Key: "your-key"
    value_template: "{{ value_json.status }}"
```

## Running Tests

```bash
go test ./...
```
