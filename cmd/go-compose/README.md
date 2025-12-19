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
