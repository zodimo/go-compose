package build

import (
	"flag"
	"fmt"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	var (
		target     = fs.String("target", "desktop", "Target platform (android, js, desktop)")
		output     = fs.String("o", "", "Output file or directory")
		ldflags    = fs.String("ldflags", "", "Arguments to pass on each go tool link invocation")
		androidAPI = fs.Int("api", 35, "Android API version (e.g., 35)")
		ndkVersion = fs.String("ndk", "", "NDK version (e.g., 27.2.12479018)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	pkgPath := "."
	if fs.NArg() > 0 {
		pkgPath = fs.Arg(0)
	}

	switch *target {
	case "js":
		if *output == "" {
			*output = "dist"
		}
		return BuildJS(*output, pkgPath, *ldflags)
	case "android":
		if *output == "" {
			// output is optional in BuildAndroid logic, handled there
		}
		return BuildAndroid(*output, []string{pkgPath}, *androidAPI, *ndkVersion, *ldflags)
	case "desktop":
		return BuildDesktop(*output, pkgPath, *ldflags)
	default:
		return fmt.Errorf("unknown target: %s", *target)
	}
}
