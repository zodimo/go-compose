package build

import (
	"flag"
	"fmt"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	var (
		target = fs.String("target", "desktop", "Target platform (android, js, desktop)")
		output = fs.String("o", "", "Output file or directory")
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
		return BuildJS(*output, pkgPath)
	case "android":
		if *output == "" {
			// output is optional in BuildAndroid logic, handled there
		}
		return BuildAndroid(*output, []string{pkgPath})
	case "desktop":
		return BuildDesktop(*output, pkgPath)
	default:
		return fmt.Errorf("unknown target: %s", *target)
	}
}
