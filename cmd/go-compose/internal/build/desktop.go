package build

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func BuildDesktop(output string, pkgPath string, ldflags string, useTinygo bool) error {
	var args []string
	var cmd *exec.Cmd

	if useTinygo {
		fmt.Printf("Building for desktop (%s/%s) with TinyGo (smaller binary)...\n", runtime.GOOS, runtime.GOARCH)
		args = []string{"build"}
		// TinyGo doesn't support ldflags the same way, but we pass them anyway
		if ldflags != "" {
			args = append(args, "-ldflags", ldflags)
		}
		if output != "" {
			args = append(args, "-o", output)
		}
		args = append(args, pkgPath)
		cmd = exec.Command("tinygo", args...)
	} else {
		fmt.Printf("Building for desktop (%s/%s)...\n", runtime.GOOS, runtime.GOARCH)
		args = []string{"build"}
		if ldflags != "" {
			args = append(args, "-ldflags", ldflags)
		}
		if output != "" {
			args = append(args, "-o", output)
		}
		args = append(args, pkgPath)
		cmd = exec.Command("go", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // Inherit environment

	return cmd.Run()
}
