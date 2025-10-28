# eBPF Profiler Trigger

Small executable that checks every X seconds the contents of `config.json`.
If the content changes to true, execute a command (default `ping 8.8.8.8`).
If the content changes to false, kill the running process.

## Executing application

1. Set values for `config.json`
2. Run `./ebpf-profiler-trigger`
3. Activate profiler by setting `enabled` to `true`
4. Disable profiler by setting `enabled` to `false`