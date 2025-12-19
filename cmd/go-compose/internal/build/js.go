package build

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

// BuildJS builds the package at pkgPath for the JS/WASM target and writes the output to outputDir.
func BuildJS(outputDir string, pkgPath string) error {

	if err := os.MkdirAll(outputDir, 0o700); err != nil {
		return err
	}

	// 1. Build WASM binary
	cmd := exec.Command(
		"go",
		"build",
		"-o", filepath.Join(outputDir, "main.wasm"),
		pkgPath,
	)
	cmd.Env = append(
		os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)
	// Add stderr to see build errors
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Println("Building WASM binary...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	// 2. Generate HTML
	// TODO: Get real name/icon from build info later
	name := filepath.Base(pkgPath)
	if name == "." || name == "/" {
		wd, _ := os.Getwd()
		name = filepath.Base(wd)
	}

	indexTemplate, err := template.New("").Parse(jsIndex)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err := indexTemplate.Execute(&b, struct {
		Name string
		Icon string
	}{
		Name: name,
		Icon: "", // TODO: Support icon
	}); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(outputDir, "index.html"), b.Bytes(), 0o600); err != nil {
		return err
	}

	// 3. Copy wasm_exec.js
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return fmt.Errorf("failed to get GOROOT: %w", err)
	}
	trimmedGoroot := strings.TrimSpace(string(goroot))

	// Try new location (Go 1.24+)
	wasmJS := filepath.Join(trimmedGoroot, "lib", "wasm", "wasm_exec.js")
	if _, err := os.Stat(wasmJS); err != nil {
		// Try old location
		wasmJS = filepath.Join(trimmedGoroot, "misc", "wasm", "wasm_exec.js")
		if _, err := os.Stat(wasmJS); err != nil {
			return fmt.Errorf("failed to find wasm_exec.js in GOROOT: %v", err)
		}
	}

	// TODO: Handle extra JS from packages (like gio-mw might need)
	// For now just basic wasm_exec.js

	// Copy wasm_exec.js to outputDir
	if err := copyFile(wasmJS, filepath.Join(outputDir, "wasm_exec.js")); err != nil {
		return err
	}

	// Create wasm.js wrapper (merging wasm_exec.js and our glue code)
	// For simplicity in this first pass, we'll just use individual files in HTML
	// but gio-cmd merges them into wasm.js. Let's follow gio-cmd's pattern for better compatibility.

	// Allow finding extra JS files from dependencies
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps,
		Env:  append(os.Environ(), "GOOS=js", "GOARCH=wasm"),
	}, pkgPath)
	if err != nil {
		return fmt.Errorf("failed to load packages: %w", err)
	}

	extraJS, err := findPackagesJS(pkgs[0], make(map[string]bool))
	if err != nil {
		return fmt.Errorf("failed to find extra JS: %w", err)
	}

	return mergeJSFiles(filepath.Join(outputDir, "wasm.js"), append([]string{wasmJS}, extraJS...)...)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func findPackagesJS(p *packages.Package, visited map[string]bool) (extraJS []string, err error) {
	if len(p.GoFiles) == 0 {
		return nil, nil
	}
	// Looking for *_js.js files in the package directory
	js, err := filepath.Glob(filepath.Join(filepath.Dir(p.GoFiles[0]), "*_js.js"))
	if err != nil {
		return nil, err
	}
	extraJS = append(extraJS, js...)

	for _, imp := range p.Imports {
		if !visited[imp.ID] {
			extra, err := findPackagesJS(imp, visited)
			if err != nil {
				return nil, err
			}
			extraJS = append(extraJS, extra...)
			visited[imp.ID] = true
		}
	}
	return extraJS, nil
}

func mergeJSFiles(dst string, files ...string) (err error) {
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := w.Close(); err != nil {
			err = cerr
		}
	}()

	// Prepend jsSetGo
	if _, err := io.Copy(w, strings.NewReader(jsSetGo)); err != nil {
		return err
	}

	for i := range files {
		r, err := os.Open(files[i])
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			r.Close()
			return err
		}
		r.Close()
		// Add newline between files to be safe
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// Append jsStartGo
	if _, err := io.Copy(w, strings.NewReader(jsStartGo)); err != nil {
		return err
	}
	return nil
}

const (
	jsIndex = `<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, user-scalable=no">
		<meta name="mobile-web-app-capable" content="yes">
		{{ if .Icon }}<link rel="icon" href="{{.Icon}}" type="image/x-icon" />{{ end }}
		{{ if .Name }}<title>{{.Name}}</title>{{ end }}
		<script src="wasm.js"></script>
		<style>
			body,pre { margin:0;padding:0; }
		</style>
	</head>
	<body>
	</body>
</html>`

	// jsSetGo sets the `window.go` variable.
	jsSetGo = `(() => {
    window.go = {argv: [], env: {}, importObject: {go: {}}};
	const argv = new URLSearchParams(location.search).get("argv");
	if (argv) {
		window.go["argv"] = argv.split(" ");
	}
})();`

	// jsStartGo initializes the main.wasm.
	jsStartGo = `(() => {
	defaultGo = new Go();
	Object.assign(defaultGo["argv"], defaultGo["argv"].concat(go["argv"]));
	Object.assign(defaultGo["env"], go["env"]);
	for (let key in go["importObject"]) {
		if (typeof defaultGo["importObject"][key] === "undefined") {
			defaultGo["importObject"][key] = {};
		}
		Object.assign(defaultGo["importObject"][key], go["importObject"][key]);
	}
	window.go = defaultGo;
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
})();`
)
