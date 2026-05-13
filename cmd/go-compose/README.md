# go-compose CLI

The `go-compose` command provides multiplatform build and serve tools for GoCompose applications.

## Installation

```bash
go install github.com/zodimo/go-compose/cmd/go-compose@latest
```

Or from source:

```bash
cd cmd/go-compose
go install .
```

## Usage

### Web (WASM)

Serve your application locally with hot-reload support (future):

```bash
go-compose serve -http :8080 ./cmd/demo/kitchen
```

Build for production (generates `dist/` with `main.wasm`, `index.html`, and `wasm_exec.js`):

```bash
go-compose build -target js -o dist ./cmd/demo/kitchen
```

### TinyGo Support (Smaller Binaries)

Use the `-tinygo` flag with build commands to create significantly smaller binaries:

```bash
# Native desktop with TinyGo (recommended)
go-compose build -target desktop -tinygo -o myapp ./cmd/demo/kitchen

# Android with TinyGo
go-compose build -target android -tinygo -o myapp.apk ./cmd/demo/kitchen
```

**Note:** TinyGo is automatically disabled for WASM/web targets due to [net/http compatibility issues](https://github.com/tinygo-org/tinygo/issues/4420). Standard Go is used instead with automatic size optimizations applied by default.

*Note:* TinyGo must be installed separately. Some features may not be fully supported depending on the target platform.

### WASM Binary Size

Go WASM binaries are inherently large (30-40MB for GUI applications) because they include the full Go runtime, garbage collector, and standard library. This is a known limitation of the standard Go compiler.

**Size Reduction Strategies:**

1. **Compression** (Recommended): Serve the WASM file with gzip or brotli compression. This typically achieves 60-70% size reduction:
   ```bash
   # Build generates both main.wasm and main.wasm.gz
   go-compose build -target js -o dist ./cmd/demo/kitchen
   
   # Serve with compression enabled
   # nginx example:
   # gzip_static on;
   # brotli_static on;
   ```

2. **Use TinyGo for Desktop**: If binary size is critical, consider building for desktop instead:
   ```bash
   go-compose build -target desktop -tinygo -o myapp ./cmd/demo/kitchen
   ```

**Expected Sizes:**
- Kitchen demo (WASM): ~38MB uncompressed, ~12MB gzipped
- Desktop (with TinyGo): ~5-15MB depending on dependencies

### Understanding Lazy Evaluation

Go-compose uses **lazy evaluation** (not lazy loading). This means composables are only evaluated/rendered when their conditions are met, saving memory by not creating UI elements that aren't currently visible.

**Key difference:**
- **Lazy Evaluation**: Components are only evaluated when needed (conditional rendering). The code is still in the binary, but UI elements aren't created until their condition is true.
- **Lazy Loading**: Code modules are loaded on demand (not supported by Go WASM).

**Example of lazy evaluation in go-compose:**

```go
// Only evaluates HeavyScreen() when showHeavy is true
c.WhenLazy(showHeavy, func() api.Composable {
    return HeavyScreen()
})

// vs c.When() which evaluates immediately even if not shown
c.When(showHeavy, HeavyScreen)  // HeavyScreen() is called every frame
```

See `cmd/demo/kitchen` for a practical example using navigation with `c.WhenLazy()` to only render the active screen.

### Android

Build an APK:

```bash
go-compose build -target android -o myapp.apk ./cmd/demo/kitchen
```

*Requirements:*
- `ANDROID_HOME` environment variable must be set.
- Android NDK and SDK Build Tools must be installed.

### Desktop

Build for the current OS:

```bash
go-compose build -target desktop -o myapp ./cmd/demo/kitchen
```
