# eBPF Profiler Trigger

Small executable that checks every X seconds the contents of `enable_profiler.txt`.
If the content changes to true, execute a command (for now `ping 8.8.8.8`).
If the content changes to false, kill the running process.

## executing application

1. run `echo false > enable_profiler.txt`
2. run `./ebpf-profiler-trigger`
3. Activate profiler by setting the text to 'true' in `enable_profiler.txt` (e.g. `echo true > enable_profiler.txt`)
4. Disable profiler by setting the text to 'false' in `enable_profiler.txt` (e.g. `echo false > enable_profiler.txt`)