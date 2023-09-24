# About

[![build](https://github.com/rgl/python-wazero-poc/actions/workflows/build.yml/badge.svg)](https://github.com/rgl/python-wazero-poc/actions/workflows/build.yml)

**THIS A WIP PoC. Once done, the final repo will be moved to another location.**

See https://github.com/tetratelabs/wazero/issues/1733.

This will try to create a single-file binary to execute an embedded Python script.

This is implemented as a Go application that executes an embedded Python script with [CPython WASI/WASM](https://github.com/brettcannon/cpython-wasi-build/) and [wazero](https://github.com/tetratelabs/wazero).

# Usage (Ubuntu 22.04)

Install `go`, `make`, and `unzip`.

Build and run the application:

```bash
make
```

# References

* https://github.com/tetratelabs/wazero
* https://github.com/brettcannon/cpython-wasi-build/
* https://snarky.ca/wasi-support-for-cpython-june-2023/
