# go-compose

## Vision
To provide a feature-complete, idiomatic Go implementation of Jetpack Compose's Material 3 API, enabling developers to build beautiful, cross-platform applications with a familiar declarative syntax.


# Run the demo's in the browser
```bash

go run ./cmd/go-compose  serve ./cmd/demo/kitchen/

```


# Status
- in early development

# Known issues
This is a non exhaustive list of known issues:

- the custom colors does not always propagate to the children
- package level typography is not implemented
- limit custom font support


# CLI Tool

The `go-compose` CLI provides build and serve commands for all supported platforms.

## Installation

```bash
go install github.com/zodimo/go-compose/cmd/go-compose@latest
```

## Commands

### `serve` - Development Server

Run your application locally with hot-reload support:

```bash
# Basic usage
go-compose serve ./cmd/demo/kitchen

# Custom port
go-compose serve -http :3000 ./cmd/demo/kitchen

# With TinyGo for smaller binaries
go-compose serve -tinygo ./cmd/demo/kitchen
```

**Flags:**
- `-http` - HTTP bind address (default: `:8080`)
- `-tinygo` - Use TinyGo compiler for smaller binaries

### `build` - Build for Production

Build your application for various targets.

#### Web (WASM)

```bash
# Build for web (outputs to dist/)
go-compose build -target js ./cmd/demo/kitchen

# Custom output directory
go-compose build -target js -o ./build ./cmd/demo/kitchen

# With TinyGo (significantly smaller WASM)
go-compose build -target js -tinygo -o ./dist ./cmd/demo/kitchen
```

**Output:** Generates `dist/` with `main.wasm`, `index.html`, and `wasm.js`

#### Desktop

```bash
# Build for current OS
go-compose build -target desktop ./cmd/demo/kitchen

# With custom output name
go-compose build -target desktop -o myapp ./cmd/demo/kitchen

# With TinyGo (smaller binary)
go-compose build -target desktop -tinygo -o myapp ./cmd/demo/kitchen
```

#### Android

```bash
# Build APK
go-compose build -target android ./cmd/demo/kitchen

# With custom output
go-compose build -target android -o myapp.apk ./cmd/demo/kitchen

# Specify API and NDK versions
go-compose build -target android -api 33 -ndk 27.2.12479018 ./cmd/demo/kitchen

# With TinyGo (experimental)
go-compose build -target android -tinygo ./cmd/demo/kitchen
```

**Requirements:**
- `ANDROID_HOME` environment variable must be set
- Android NDK and SDK Build Tools must be installed

**Build Flags:**
- `-target` - Target platform: `js`, `desktop`, `android` (default: `desktop`)
- `-o` - Output file or directory
- `-tinygo` - Use TinyGo compiler for smaller binaries
- `-ldflags` - Arguments to pass to the Go linker
- `-api` - Android API version (default: `35`)
- `-ndk` - NDK version (e.g., `27.2.12479018`)

## TinyGo Support

TinyGo can produce significantly smaller binaries, especially for WASM targets.

### Installation

```bash
# macOS
brew install tinygo

# Ubuntu/Debian
wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.2/tinygo_0.31.2_amd64.deb
dpkg -i tinygo_0.31.2_amd64.deb

# See https://tinygo.org/getting-started/install/ for other platforms
```

### Usage Examples

```bash
# WASM: ~10x smaller binaries
go-compose build -target js -tinygo ./cmd/demo/kitchen

# Desktop: Smaller native binaries
go-compose build -target desktop -tinygo -o app ./cmd/demo/kitchen

# Development server with TinyGo
go-compose serve -tinygo ./cmd/demo/kitchen
```

**Notes:**
- TinyGo has limited support for some Go features (reflection, some stdlib packages)
- Not all applications may compile with TinyGo
- Android support is experimental

# License
Distributed under the MIT License. See `LICENSE` for more information.
