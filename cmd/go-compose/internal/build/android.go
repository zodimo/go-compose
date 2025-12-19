package build

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/packages"
)

type androidTools struct {
	buildtools string
	androidjar string
}

type manifestData struct {
	AppID       string
	Version     Semver
	MinSDK      int
	TargetSDK   int
	Permissions []string
	Features    []string
	IconSnip    string
	AppName     string
}

const (
	themes = `<?xml version="1.0" encoding="utf-8"?>
<resources>
	<style name="Theme.GioApp" parent="android:style/Theme.NoTitleBar">
		<item name="android:windowBackground">@android:color/white</item>
	</style>
</resources>`
	themesV21 = `<?xml version="1.0" encoding="utf-8"?>
<resources>
	<style name="Theme.GioApp" parent="android:style/Theme.NoTitleBar">
		<item name="android:windowBackground">@android:color/white</item>

		<item name="android:windowDrawsSystemBarBackgrounds">true</item>
		<item name="android:navigationBarColor">#40000000</item>
		<item name="android:statusBarColor">#40000000</item>
	</style>
</resources>`
)

var (
	AndroidPermissions = map[string][]string{
		"network":      {"android.permission.INTERNET"},
		"networkstate": {"android.permission.ACCESS_NETWORK_STATE"},
		"storage":      {"android.permission.READ_EXTERNAL_STORAGE", "android.permission.WRITE_EXTERNAL_STORAGE"},
	}
	AndroidFeatures = map[string][]string{
		"default": {`glEsVersion="0x00020000"`, `name="android.hardware.type.pc"`},
	}
)

func BuildAndroid(output string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package path required")
	}
	pkgPath := args[0]
	// TODO: Parse ID, Version, Key, etc. from args or flags.
	// For now using defaults for proof of concept.

	bi := &BuildInfo{
		AppID:     "org.gioui.experiment", // default
		Version:   Semver{Major: 1, Minor: 0, Patch: 0, VersionCode: 1},
		MinSDK:    21, // Good default
		TargetSDK: 33,
		Name:      "GioApp", // default
		PkgPath:   pkgPath,
		Archs:     []string{"arm64", "amd64"}, // Common archs
		Tags:      "",
	}

	tmpDir, err := os.MkdirTemp("", "go-compose-android-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	return buildAndroid(tmpDir, bi, output)
}

func buildAndroid(tmpDir string, bi *BuildInfo, outputFile string) error {
	sdk := os.Getenv("ANDROID_HOME")
	if sdk == "" {
		return errors.New("please set ANDROID_HOME to the Android SDK path")
	}
	platform, err := latestPlatform(sdk)
	if err != nil {
		return err
	}
	buildtools, err := latestTools(sdk)
	if err != nil {
		return err
	}

	tools := &androidTools{
		buildtools: buildtools,
		androidjar: filepath.Join(platform, "android.jar"),
	}

	// Scan permissions
	perms := []string{"default"}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps,
		Env:  append(os.Environ(), "GOOS=android", "CGO_ENABLED=1"),
	}
	pkgs, err := packages.Load(cfg, bi.PkgPath)
	if err != nil {
		return err
	}

	// Simple permission check (can be expanded)
	visitedPkgs := make(map[string]bool)
	var visitPkg func(*packages.Package)
	visitPkg = func(p *packages.Package) {
		if p.PkgPath == "net" {
			perms = append(perms, "network")
		}
		for _, imp := range p.Imports {
			if !visitedPkgs[imp.ID] {
				visitedPkgs[imp.ID] = true
				visitPkg(imp)
			}
		}
	}
	visitPkg(pkgs[0])

	if err := compileAndroid(tmpDir, tools, bi); err != nil {
		return err
	}

	// Exe mode (APK)
	if outputFile == "" {
		outputFile = bi.Name + ".apk"
	}

	extraJars := []string{} // TODO: Support finding jars in packages

	if err := exeAndroid(tmpDir, tools, bi, extraJars, perms, false); err != nil {
		return err
	}

	return signAPK(tmpDir, outputFile, tools, bi)
}

func compileAndroid(tmpDir string, tools *androidTools, bi *BuildInfo) error {
	androidHome := os.Getenv("ANDROID_HOME")
	ndkRoot, err := findNDK(androidHome)
	if err != nil {
		return err
	}
	minSDK := bi.MinSDK
	if minSDK < 17 {
		minSDK = 17
	}

	tcRoot := filepath.Join(ndkRoot, "toolchains", "llvm", "prebuilt", archNDK())
	var builds errgroup.Group

	for _, a := range bi.Archs {
		arch := allArchs[a]
		clang, err := latestCompiler(tcRoot, a, minSDK)
		if err != nil {
			return err
		}

		archDir := filepath.Join(tmpDir, "jni", arch.jniArch)
		if err := os.MkdirAll(archDir, 0o755); err != nil {
			return err
		}
		libFile := filepath.Join(archDir, "libgio.so")
		cmd := exec.Command(
			"go",
			"build",
			"-ldflags=-w -s -extldflags \"-Wl,-z,max-page-size=65536\"",
			"-buildmode=c-shared",
			"-tags", bi.Tags,
			"-o", libFile,
			bi.PkgPath,
		)
		cmd.Env = append(
			os.Environ(),
			"GOOS=android",
			"GOARCH="+a,
			"GOARM=7",
			"CGO_ENABLED=1",
			"CC="+clang,
		)
		// Redirect output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		builds.Go(func() error {
			fmt.Printf("Building for %s...\n", a)
			return cmd.Run()
		})
	}

	// Compile Java shim (gioui.org/app)
	// For simplicity, we assume gioui.org/app is available in module cache or vendor
	// We need to find where gioui.org/app is located on disk.
	appDirVal, err := runCmd(exec.Command("go", "list", "-f", "{{.Dir}}", "gioui.org/app"))
	if err != nil {
		return fmt.Errorf("failed to find gioui.org/app: %w", err)
	}
	appDir := strings.TrimSpace(appDirVal)

	javaFiles, _ := filepath.Glob(filepath.Join(appDir, "*.java"))
	if len(javaFiles) > 0 {
		classes := filepath.Join(tmpDir, "classes")
		if err := os.MkdirAll(classes, 0o755); err != nil {
			return err
		}
		javac, err := findJavaC()
		if err != nil {
			return err
		}
		javacCmd := exec.Command(
			javac,
			"-target", "1.8",
			"-source", "1.8",
			"-sourcepath", appDir,
			"-bootclasspath", tools.androidjar,
			"-d", classes,
		)
		javacCmd.Args = append(javacCmd.Args, javaFiles...)
		builds.Go(func() error {
			_, err := runCmd(javacCmd)
			return err
		})
	}

	return builds.Wait()
}

func exeAndroid(tmpDir string, tools *androidTools, bi *BuildInfo, extraJars, perms []string, isBundle bool) error {
	classes := filepath.Join(tmpDir, "classes")
	var classFiles []string
	filepath.Walk(classes, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".class" {
			classFiles = append(classFiles, path)
		}
		return nil
	})
	classFiles = append(classFiles, extraJars...)

	dexDir := filepath.Join(tmpDir, "apk")
	if err := os.MkdirAll(dexDir, 0o755); err != nil {
		return err
	}

	minSDK := bi.MinSDK // Simplified logic

	if len(classFiles) > 0 {
		d8 := exec.Command(
			filepath.Join(tools.buildtools, "d8"),
			"--lib", tools.androidjar,
			"--output", dexDir,
			"--min-api", strconv.Itoa(minSDK),
		)
		d8.Args = append(d8.Args, classFiles...)
		if _, err := runCmd(d8); err != nil {
			return err
		}
	}

	// Resources
	resDir := filepath.Join(tmpDir, "res")
	valDir := filepath.Join(resDir, "values")
	v21Dir := filepath.Join(resDir, "values-v21")
	os.MkdirAll(valDir, 0o755)
	os.MkdirAll(v21Dir, 0o755)

	os.WriteFile(filepath.Join(valDir, "themes.xml"), []byte(themes), 0o660)
	os.WriteFile(filepath.Join(v21Dir, "themes.xml"), []byte(themesV21), 0o660)

	resZip := filepath.Join(tmpDir, "resources.zip")
	aapt2 := filepath.Join(tools.buildtools, "aapt2")
	if _, err := runCmd(exec.Command(aapt2, "compile", "-o", resZip, "--dir", resDir)); err != nil {
		return err
	}

	permissions, features := getPermissions(perms)
	manifestSrc := manifestData{
		AppID:       bi.AppID,
		MinSDK:      bi.MinSDK,
		TargetSDK:   bi.TargetSDK,
		Permissions: permissions,
		Features:    features,
		AppName:     bi.Name,
	}

	tmpl, err := template.New("manifest").Parse(
		`<manifest xmlns:android="http://schemas.android.com/apk/res/android" package="{{.AppID}}">
        <uses-sdk android:minSdkVersion="{{.MinSDK}}" android:targetSdkVersion="{{.TargetSDK}}"/>
{{range .Permissions}}	<uses-permission android:name="{{.}}"/>
{{end}}{{range .Features}}	<uses-feature android:{{.}} android:required="false"/>
{{end}}	<application android:label="{{.AppName}}">
		<activity android:name="org.gioui.GioActivity"
			android:label="{{.AppName}}"
			android:theme="@style/Theme.GioApp"
			android:configChanges="screenSize|screenLayout|smallestScreenSize|orientation|keyboardHidden"
			android:windowSoftInputMode="adjustResize"
			android:exported="true">
			<intent-filter>
				<action android:name="android.intent.action.MAIN" />
				<category android:name="android.intent.category.LAUNCHER" />
			</intent-filter>
		</activity>
	</application>
</manifest>`)
	if err != nil {
		return err
	}

	var manifestBuffer bytes.Buffer
	if err := tmpl.Execute(&manifestBuffer, manifestSrc); err != nil {
		return err
	}
	manifest := filepath.Join(tmpDir, "AndroidManifest.xml")
	os.WriteFile(manifest, manifestBuffer.Bytes(), 0o660)

	linkAPK := filepath.Join(tmpDir, "link.apk")
	if _, err := runCmd(exec.Command(aapt2, "link", "--manifest", manifest, "-I", tools.androidjar, "-o", linkAPK, resZip)); err != nil {
		return err
	}

	// Create final APK by merging link.apk + dex + libs
	unsignedAPK := filepath.Join(tmpDir, "app.zip")
	// ... logic to zip implementation ...
	// Re-using simplified zip merging logic

	return createUnsignedAPK(unsignedAPK, linkAPK, dexDir, tmpDir, bi)
}

func createUnsignedAPK(unsignedAPK, linkAPK, dexDir, tmpDir string, bi *BuildInfo) error {
	// ... (Implementation of zip merging similar to androidbuild.go)
	// Simplified copy:
	f, err := os.Create(unsignedAPK)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	// Copy link.apk
	r, err := zip.OpenReader(linkAPK)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fw, _ := w.Create(file.Name)
		fr, _ := file.Open()
		io.Copy(fw, fr)
		fr.Close()
	}

	// Add classes.dex
	dexFile := filepath.Join(dexDir, "classes.dex")
	if _, err := os.Stat(dexFile); err == nil {
		fw, _ := w.Create("classes.dex")
		fr, _ := os.Open(dexFile)
		io.Copy(fw, fr)
		fr.Close()
	}

	// Add libs
	for _, a := range bi.Archs {
		arch := allArchs[a]
		libName := filepath.Join("lib", arch.jniArch, "libgio.so")
		srcLib := filepath.Join(tmpDir, "jni", arch.jniArch, "libgio.so")
		if _, err := os.Stat(srcLib); err == nil {
			fw, _ := w.Create(filepath.ToSlash(libName))
			fr, _ := os.Open(srcLib)
			io.Copy(fw, fr)
			fr.Close()
		}
	}

	return nil
}

func signAPK(tmpDir string, apkFile string, tools *androidTools, bi *BuildInfo) error {
	// Align
	alignedPC := filepath.Join(tmpDir, "aligned.apk")
	zipalign := filepath.Join(tools.buildtools, "zipalign")
	runCmd(exec.Command(zipalign, "-f", "4", filepath.Join(tmpDir, "app.zip"), alignedPC))

	// Generate debug key if needed
	if bi.Key == "" {
		bi.Key = filepath.Join(tmpDir, "debug.keystore")
		bi.Password = "android"
		keytool, _ := findKeytool() // assume found
		runCmd(exec.Command(keytool, "-genkey", "-keystore", bi.Key, "-storepass", bi.Password, "-alias", "android", "-keyalg", "RSA", "-keysize", "2048", "-validity", "10000", "-noprompt", "-dname", "CN=android"))
	}

	apksigner := filepath.Join(tools.buildtools, "apksigner")
	_, err := runCmd(exec.Command(apksigner, "sign", "--ks-pass", "pass:"+bi.Password, "--ks", bi.Key, "--out", apkFile, alignedPC))
	return err
}

// Helpers

func latestPlatform(sdk string) (string, error) {
	platforms, _ := filepath.Glob(filepath.Join(sdk, "platforms", "android-*"))
	if len(platforms) == 0 {
		return "", errors.New("no android platform found")
	}
	return platforms[len(platforms)-1], nil
}

func latestTools(sdk string) (string, error) {
	tools, _ := filepath.Glob(filepath.Join(sdk, "build-tools", "*"))
	if len(tools) == 0 {
		return "", errors.New("no build-tools found")
	}
	return tools[len(tools)-1], nil
}

func findNDK(sdk string) (string, error) {
	ndks, _ := filepath.Glob(filepath.Join(sdk, "ndk", "*"))
	if len(ndks) > 0 {
		return ndks[len(ndks)-1], nil
	}
	return "", errors.New("no NDK found")
}

func archNDK() string {
	if runtime.GOOS == "linux" {
		return "linux-x86_64"
	}
	if runtime.GOOS == "darwin" {
		return "darwin-x86_64"
	}
	if runtime.GOOS == "windows" {
		return "windows-x86_64"
	}
	return ""
}

func latestCompiler(tcRoot, arch string, minSDK int) (string, error) {
	// Simplified detection
	pre := ""
	switch arch {
	case "arm":
		pre = "armv7a-linux-androideabi"
	case "arm64":
		pre = "aarch64-linux-android"
	case "amd64":
		pre = "x86_64-linux-android"
	case "386":
		pre = "i686-linux-android"
	}
	return filepath.Join(tcRoot, "bin", fmt.Sprintf("%s%d-clang", pre, minSDK)), nil
}

func findJavaC() (string, error) {
	// Simplified
	return "javac", nil
}

func findKeytool() (string, error) {
	return "keytool", nil
}

// Permissions helper
func getPermissions(perms []string) ([]string, []string) {
	var permissions []string
	var features []string
	seen := make(map[string]bool)
	for _, p := range perms {
		if seen[p] {
			continue
		}
		seen[p] = true
		if ps, ok := AndroidPermissions[p]; ok {
			permissions = append(permissions, ps...)
		}
		if fs, ok := AndroidFeatures[p]; ok {
			features = append(features, fs...)
		}
	}
	return permissions, features
}

type arch struct {
	jniArch string
}

var allArchs = map[string]arch{
	"arm":   {jniArch: "armeabi-v7a"},
	"arm64": {jniArch: "arm64-v8a"},
	"386":   {jniArch: "x86"},
	"amd64": {jniArch: "x86_64"},
}
