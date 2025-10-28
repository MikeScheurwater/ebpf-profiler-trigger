# eBPF Profiler Trigger

A lightweight Go-based utility that automatically starts and stops a custom command based on configuration stored in a
JSON file.

## Overview

This service monitors a `config.json` file and reacts to changes in the `enabled` flag:

- When `"enabled": true`, it starts the configured command.
- When `"enabled": false`, it stops the running process.
- The file is checked periodically based on the value of `"poll_interval"` (in seconds).

This allows you to remotely toggle profiling or tracing without manually managing background
processes.

---

## Example Configuration (`config.json`)

```json
{
  "enabled": false,
  "poll_interval": 2,
  "command": "ping",
  "args": [
    "8.8.8.8"
  ]
}
```