package build

// BuildInfo contains the configuration for a build.
type BuildInfo struct {
	AppID     string
	Version   Semver
	MinSDK    int
	TargetSDK int
	Name      string
	PkgPath   string // Go package path (e.g. github.com/user/app)
	IconPath  string
	Tags      string
	Key       string // Path to keystore
	Password  string // Keystore password
	Archs     []string
}

type Semver struct {
	Major, Minor, Patch int
	VersionCode         uint32
}

func (s Semver) String() string {
	// Simple string representation, can be expanded
	return "1.0.0.1" // Placeholder implementation
}
