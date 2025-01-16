Analyze the k6 test script and provision k6 with extensions based on dependencies.

### Sources

Dependencies can come from three sources: k6 test script, manifest file, `K6_DEPENDENCIES` environment variable. Instead of these three sources, a k6 archive can also be specified, which can contain all three sources.

### Output

By default, the k6 executable is created in the current directory with the name `k6` (`k6.exe` on Windows). This can be overridden by using the `--o/-output` flag.
