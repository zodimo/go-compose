package build

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(cmd *exec.Cmd) (string, error) {
	fmt.Printf("Running: %s\n", strings.Join(cmd.Args, " "))
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s failed: %s%s", strings.Join(cmd.Args, " "), out, exitErr.Stderr)
		}
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}

func runCmdRaw(cmd *exec.Cmd) ([]byte, error) {
	fmt.Printf("Running: %s\n", strings.Join(cmd.Args, " "))
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s failed: %s%s", strings.Join(cmd.Args, " "), out, exitErr.Stderr)
		}
		return nil, err
	}
	return out, nil
}

// copyFile is already defined in js.go, but we might want to consolidate or make sure they don't conflict if they are in same package.
// If they are in the same package 'build', we can't have duplicate functions.
// I should check if I already defined copyFile in js.go. Yes I did.
// So I should REMOVE copyFile from here if it's already in js.go, or rename the one in js.go to be this one.
// Logic in js.go was unexported `copyFile`.
// I will NOT redeclare it here to avoid conflicts, assuming js.go is in the same package.
// But wait, if I build just this file it's fine, but as a package, it will conflict.
// I will rely on the one in js.go or move it here later.
// Actually, `js.go` is in `package build`. So `copyFile` is visible to `android.go`.
// I'll skip copyFile here.
