package build

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func BuildDesktop(output string, pkgPath string) error {
	// Standard go build
	args := []string{"build"}
	if output != "" {
		args = append(args, "-o", output)
	}
	// args = append(args, "-tags", "example") // Add default tags if needed
	args = append(args, pkgPath)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // Inherit environment

	fmt.Printf("Building for desktop (%s/%s)...\n", runtime.GOOS, runtime.GOARCH)
	return cmd.Run()
}
